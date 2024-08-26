package app

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"sync"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/s0vunia/auth_microservice/internal/interceptor"
	"github.com/s0vunia/auth_microservice/internal/logger"
	"github.com/s0vunia/auth_microservice/internal/metric"
	"github.com/s0vunia/auth_microservice/internal/tracing"
	descAccess "github.com/s0vunia/auth_microservice/pkg/access_v1"
	descAuth "github.com/s0vunia/auth_microservice/pkg/auth_v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	// nolint
	_ "github.com/s0vunia/auth_microservice/statik"
	"github.com/s0vunia/platform_common/pkg/closer"

	"github.com/s0vunia/auth_microservice/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	desc "github.com/s0vunia/auth_microservice/pkg/user_v1"
)

var (
	configPath  string
	serviceName = "auth-service"
)

func init() {
	configPath = os.Getenv("CONFIG_PATH")
}

// App represents the app.
type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
}

// NewApp creates a new app.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run runs the app.
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	logger.Info(a.serviceProvider.LoggerConfig().FileName())

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logger.Fatal(
				"failed to run GRPC server",
				zap.Error(err),
			)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logger.Fatal(
				"failed to run HTTP server",
				zap.Error(err),
			)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			logger.Fatal(
				"failed to run Swagger server",
				zap.Error(err),
			)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.serviceProvider.UserSaverConsumer(ctx).RunConsumer(ctx)
		if err != nil {
			logger.Fatal(
				"failed to run user saver consumer",
				zap.Error(err),
			)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.runPrometheus()
		if err != nil {
			logger.Fatal(
				"failed to run prometheus",
				zap.Error(err),
			)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initMetric,
		a.initTracing,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initPrometheusServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.MetricsInterceptor,
				interceptor.ServerTracingInterceptor,
				interceptor.LogInterceptor,
				interceptor.ValidateInterceptor,
			),
		),
	)

	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HTTPConfig().Address(),
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: a.serviceProvider.HTTPConfig().ReadHeaderTimeout(),
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.serviceProvider.SwaggerConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: a.serviceProvider.HTTPConfig().ReadHeaderTimeout(),
	}

	return nil
}

func (a *App) initPrometheusServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:              a.serviceProvider.PrometheusConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: a.serviceProvider.HTTPConfig().ReadHeaderTimeout(),
	}
	return nil
}

func (a *App) initMetric(ctx context.Context) error {
	err := metric.Init(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.getCore(a.getAtomicLevel()))
	return nil
}

func (a *App) initTracing(_ context.Context) error {
	tracing.Init(logger.Logger(), serviceName)
	return nil
}

func (a *App) runGRPCServer() error {
	logger.Info("GRPC server is running",
		zap.String("address", a.serviceProvider.GRPCConfig().Address()),
	)

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	logger.Info("HTTP server is running",
		zap.String("address", a.serviceProvider.HTTPConfig().Address()),
	)

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runPrometheus() error {
	logger.Info("Prometheus server is running",
		zap.String("address", a.serviceProvider.PrometheusConfig().Address()),
	)

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	logger.Info("Swagger server is running",
		zap.String("address", a.serviceProvider.SwaggerConfig().Address()),
	)

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		logger.Info(
			"Serve swagger file",
			zap.String("path", path),
		)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info(
			"Get swagger file",
			zap.String("path", path),
		)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(file http.File) {
			err = file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}(file)

		logger.Info(
			"Read swagger file",
			zap.String("path", path),
		)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info(
			"Serve swagger file",
			zap.String("path", path),
		)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info(
			"Serve swagger file",
			zap.String("path", path),
		)
	}
}

func (a *App) getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   a.serviceProvider.LoggerConfig().FileName(),
		MaxSize:    a.serviceProvider.LoggerConfig().MaxSize(), // megabytes
		MaxBackups: a.serviceProvider.LoggerConfig().MaxBackups(),
		MaxAge:     a.serviceProvider.LoggerConfig().MaxAge(), // days
	})
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func (a *App) getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(a.serviceProvider.LoggerConfig().Level()); err != nil {
		logger.Fatal("failed to set log level", zap.Error(err))
	}

	return zap.NewAtomicLevelAt(level)
}

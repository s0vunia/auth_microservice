package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/config"
	"github.com/s0vunia/auth_microservices_course_boilerplate/pkg/common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/s0vunia/auth_microservices_course_boilerplate/pkg/auth_v1"
)

var configPath string

func init() {
	configPath = os.Getenv("CONFIG_PATH")
}

type server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

// Create user
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("User email: %s", req.GetInfo().GetEmail())
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "len of password must be positive")
	}
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Error(codes.InvalidArgument, "password must be equal to password_confirm")
	}
	passHash, err := utils.GetHashPassword(req.GetPassword())
	if err != nil {
		log.Printf("err: %e", err)
		return nil, status.Error(codes.Internal, "failed to generate password hash")
	}
	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "pass_hash", "role").
		Values(req.GetInfo().GetName(), req.GetInfo().GetEmail(), passHash, req.GetInfo().GetRole()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("err: %e", err)
		return nil, status.Error(codes.Internal, "failed to build query")
	}
	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("err: %e", err)
		return nil, status.Error(codes.Internal, "failed to insert user")
	}

	log.Printf("inserted user with id: %d", userID)
	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

// Get user by id
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	builderSelectOne := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, status.Error(codes.Internal, "failed to build query")
	}

	var id, role int64
	var name, email string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("failed to select users: %v", err)
		return nil, status.Error(codes.Internal, "failed to select users")
	}

	log.Printf("id: %d, name: %s, email: %s, role: %d, created_at: %v, updated_at: %v\n", id, name, email, role, createdAt, updatedAt)

	var updatedAtTime *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtTime = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetResponse{
		User: &desc.User{
			Id: req.GetId(),
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  desc.Role(role),
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: updatedAtTime,
		},
	}, nil
}

// Update user credentials
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("User id: %d", req.GetId())

	if req.GetInfo().GetName() == nil && req.GetInfo().GetEmail() == nil && req.GetInfo().GetRole() == desc.Role_unknown {
		return nil, status.Error(codes.InvalidArgument, "update required at least one column")
	}
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar)
	if req.GetInfo().GetName() != nil {
		builderUpdate = builderUpdate.Set("name", req.GetInfo().GetName().GetValue())
	}
	if req.GetInfo().GetEmail() != nil {
		builderUpdate = builderUpdate.Set("email", req.GetInfo().GetEmail().GetValue())
	}
	if req.GetInfo().GetRole() != desc.Role_unknown {
		builderUpdate = builderUpdate.Set("role", int(req.GetInfo().GetRole()))
	}
	builderUpdate = builderUpdate.Set("updated_at", time.Now())

	builderUpdate = builderUpdate.
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, status.Error(codes.Internal, "failed to build query")
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	log.Printf("updated %d rows", res.RowsAffected())
	return nil, nil
}

// Delete user
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("User id: %d", req.GetId())

	builderDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, status.Error(codes.Internal, "failed to build query")
	}

	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	log.Printf("deleted %d rows", res.RowsAffected())
	return nil, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

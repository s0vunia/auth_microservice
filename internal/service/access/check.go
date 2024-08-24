package access

import (
	"context"
	"strings"

	"github.com/s0vunia/auth_microservice/internal/model"
	"github.com/s0vunia/auth_microservice/internal/utils"
	"google.golang.org/grpc/metadata"
)

func (s *serv) Check(ctx context.Context, endpointAddress string) error {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ErrMetadataNotProvided
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return ErrAuthorizationHeader
	}

	if !strings.HasPrefix(authHeader[0], s.authPrefix) {
		return ErrInvalidHeaderFormat
	}

	accessToken := strings.TrimPrefix(authHeader[0], s.authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(s.accessTokenSecretKey))
	if err != nil {
		return ErrInvalidAccessToken
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return ErrGetAccessibleRoles
	}

	roles, ok := accessibleMap[endpointAddress]
	if !ok {
		return nil
	}

	if _, ok = roles[claims.Role]; ok {
		return nil
	}

	return ErrAccessDenied
}

// Возвращает мапу с адресом эндпоинта и ролью, которая имеет доступ к нему
func (s *serv) accessibleRoles(_ context.Context) (map[string]map[model.Role]struct{}, error) {
	return accessibleRoles, nil
}

package access

import "github.com/s0vunia/auth_microservice/internal/model"

var (
	accessibleRoles = map[string]map[model.Role]struct{}{
		"/chat_v1.ChatV1/Create": {
			model.RoleAdmin: {},
		},
		"/chat_v1.ChatV1/Delete": {
			model.RoleAdmin: {},
		},
		"/chat_v1.ChatV1/SendMessage": {
			model.RoleUser:  {},
			model.RoleAdmin: {},
		},
	}
)

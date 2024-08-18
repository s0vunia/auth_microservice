package access

import "github.com/s0vunia/auth_microservice/internal/model"

var (
	accessibleRoles = map[string]map[model.Role]struct{}{
		"/chat_v1.ChatService/Create": {
			model.RoleAdmin: {},
		},
		"/chat_v1.ChatService/Delete": {
			model.RoleAdmin: {},
		},
		"/chat_v1.ChatService/SendMessage": {
			model.RoleUser:  {},
			model.RoleAdmin: {},
		},
	}
)

package access

import "errors"

var (
	ErrMetadataNotProvided = errors.New("metadata is not provided")
	ErrAuthorizationHeader = errors.New("authorization header is not provided")
	ErrInvalidHeaderFormat = errors.New("invalid authorization header format")
	ErrInvalidAccessToken  = errors.New("access token is invalid")
	ErrGetAccessibleRoles  = errors.New("failed to get accessible roles")
	ErrAccessDenied        = errors.New("access denied")
)

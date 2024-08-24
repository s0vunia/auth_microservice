package access

import "errors"

var (
	// ErrMetadataNotProvided - metadata not provided
	ErrMetadataNotProvided = errors.New("metadata is not provided")

	// ErrAuthorizationHeader - authorization header is not provided
	ErrAuthorizationHeader = errors.New("authorization header is not provided")

	// ErrInvalidHeaderFormat - invalid authorization header format
	ErrInvalidHeaderFormat = errors.New("invalid authorization header format")

	// ErrInvalidAccessToken  - invalid access token
	ErrInvalidAccessToken = errors.New("access token is invalid")

	// ErrGetAccessibleRoles  - failed to get accessible roles
	ErrGetAccessibleRoles = errors.New("failed to get accessible roles")

	// ErrAccessDenied        - access denied
	ErrAccessDenied = errors.New("access denied")
)

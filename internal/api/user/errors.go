package user

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	PasswordConfirmErr = status.Error(codes.InvalidArgument, "password and password confirm do not match")
)

package middleware

import (
	"context"
	auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/ohdat/login/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
//Auth grpc Auth 中间件
func Auth(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	uid, err := getUid(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	tags.Extract(ctx).Set("auth.uid", uid)
	//WARNING: in production define your own type to avoid context collisions
	newCtx := context.WithValue(ctx, "uid", uid)
	return newCtx, nil
}

func getUid(token string) (uid int, err error) {
	tokenInfo, err := jwt.ParseToken(token)
	if err != nil {
		return
	}
	uid = tokenInfo.Uid
	return
}

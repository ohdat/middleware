package middleware

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/ohdat/login/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"path"
)

type AuthFunc func(ctx context.Context) (context.Context, error)

// ServiceAuthFuncOverride allows a given gRPC service implementation to override the global `AuthFunc`.
//
// If a service implements the AuthFuncOverride method, it takes precedence over the `AuthFunc` method,
// and will be called instead of AuthFunc for all method invocations within that service.
type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)
}

// AuthUnaryServer returns a new unary server interceptors that performs per-request auth.
func AuthUnaryServer(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var err error
		var method = path.Base(info.FullMethod)
		if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		} else if method == "Ping" {
			//Ping跳过Auth
			newCtx = context.Background()
		} else {
			newCtx, err = authFunc(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// AuthStreamServer returns a new unary server interceptors that performs per-request auth.
func AuthStreamServer(authFunc AuthFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		var method = path.Base(info.FullMethod)
		if overrideSrv, ok := srv.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
		} else if method == "Ping" || method == "ServerReflectionInfo" {
			//跳过验证auth ServerReflectionInfo 测试环境
			newCtx = context.Background()
		} else {
			newCtx, err = authFunc(stream.Context())
		}
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

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

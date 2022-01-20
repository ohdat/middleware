package grpc_middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ohdat/game/logger"
	"google.golang.org/grpc"
	"runtime/debug"
	"time"
)

func ServerLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if e := recover(); e != nil {
				stack := debug.Stack()
				logger.Error(ctx, fmt.Sprintf("grpc-client has err:%v, stack:%v", e, string(stack)))
			}
		}()

		startTime := time.Now().UnixNano()
		ret, err := handler(ctx, req)
		duration := (time.Now().UnixNano() - startTime) / 1e6
		requestByte, _ := json.Marshal(req)
		responseStr := ""
		if err == nil {
			responseByte, _ := json.Marshal(ret)
			responseStr = string(responseByte)
		}

		logger.Info(ctx, fmt.Sprintf("方法名:%v,耗时:%vms,请求数据:%v,返回数据:%v", info.FullMethod, duration, string(requestByte), responseStr))
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("方法名:%v,耗时:%vms,请求数据:%v,返回错误:%v", info.FullMethod, duration, string(requestByte), err))
		}

		return ret, err
	}
}

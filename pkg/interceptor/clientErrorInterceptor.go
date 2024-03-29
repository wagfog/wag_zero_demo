package interceptor

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wagfog/wag_zero_demo/pkg/xcode"
	"google.golang.org/grpc"

	// "google.golang.org/grpc/internal/status"
	"google.golang.org/grpc/status"
)

func ClientErrorInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			xc := xcode.GrpcStatusToXCode(grpcStatus)
			err = errors.WithMessage(xc, grpcStatus.Message())
		}
		return err
	}

}

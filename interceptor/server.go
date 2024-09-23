package interceptor

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"jaeger-interceptors/config"
	"jaeger-interceptors/models"
)

func TracingServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		ctx = otel.GetTextMapPropagator().Extract(ctx, models.GrpcCustomPropagationHeaderCarrier{Md: md})

		ctx, span := otel.Tracer(config.TracerNameGrpcServerInterceptor).Start(ctx, info.FullMethod)
		defer span.End()

		resp, err = handler(ctx, req)

		span.SetAttributes(attribute.String("rpc.method", info.FullMethod))

		if err != nil {
			grpcStatus, _ := status.FromError(err)
			span.SetAttributes(
				attribute.Int("rpc.grpc.status_code", int(grpcStatus.Code())),
				attribute.String("rpc.error_message", grpcStatus.Message()),
			)
			span.SetStatus(otelCodes.Error, grpcStatus.Message())
		} else {
			span.SetAttributes(attribute.String("rpc.grpc.status_code", codes.OK.String()))
			span.SetStatus(otelCodes.Ok, "Success")
		}

		return resp, err
	}
}

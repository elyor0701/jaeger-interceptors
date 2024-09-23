package interceptor

import (
	"context"
	"go.opentelemetry.io/otel"
	otelCodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"jaeger-interceptors/config"
	"jaeger-interceptors/models"
	"log"
)

func TracingClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx, span := otel.Tracer(config.TracerNameGrpcClientInterceptor).Start(ctx, method)
		defer span.End()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		otel.GetTextMapPropagator().Inject(ctx, models.GrpcCustomPropagationHeaderCarrier{Md: md})

		ctx = metadata.NewOutgoingContext(ctx, md)

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
			log.Printf("gRPC error: %v", err)
		} else {
			span.SetStatus(otelCodes.Ok, "gRPC call successful")
		}

		return err
	}
}

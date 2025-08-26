package client

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func loggingUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)

	duration := time.Since(start)

	logLevel := slog.LevelInfo
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			switch s.Code() {
			case codes.Internal, codes.DataLoss, codes.Unknown:
				logLevel = slog.LevelError
			case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists:
				logLevel = slog.LevelWarn
			default:
				logLevel = slog.LevelError
			}
		} else {
			logLevel = slog.LevelError
		}
	}

	attrs := []slog.Attr{
		slog.String("type", "grpc_call"),
		slog.String("method", method),
		slog.Duration("duration", duration),
		slog.String("target", cc.Target()),
	}

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
		if s, ok := status.FromError(err); ok {
			attrs = append(attrs,
				slog.String("grpc_code", s.Code().String()),
				slog.String("grpc_message", s.Message()),
			)
		}
		slog.LogAttrs(ctx, logLevel, "gRPC Call Failed", attrs...)
	} else {
		slog.LogAttrs(ctx, slog.LevelDebug, "gRPC Call Success", attrs...)
	}

	return err
}

func retryUnaryInterceptor(maxRetries int, retryDelay time.Duration) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var lastErr error

		for attempt := 0; attempt <= maxRetries; attempt++ {
			if attempt > 0 {
				slog.WarnContext(ctx, "Retrying gRPC call",
					slog.String("method", method),
					slog.Int("attempt", attempt),
					slog.Int("max_retries", maxRetries),
					slog.String("last_error", lastErr.Error()),
				)

				select {
				case <-time.After(retryDelay):
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			err := invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				if attempt > 0 {
					slog.InfoContext(ctx, "gRPC call succeeded after retry",
						slog.String("method", method),
						slog.Int("attempts", attempt+1),
					)
				}
				return nil
			}

			lastErr = err

			if !shouldRetry(err) {
				break
			}

			if attempt == maxRetries {
				break
			}
		}

		slog.ErrorContext(ctx, "gRPC call failed after all retries",
			slog.String("method", method),
			slog.Int("total_attempts", maxRetries+1),
			slog.String("final_error", lastErr.Error()),
		)

		return lastErr
	}
}

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	s, ok := status.FromError(err)
	if !ok {
		return false
	}

	switch s.Code() {
	case codes.Unavailable,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.Internal:
		return true
	default:
		return false
	}
}

func timeoutUnaryInterceptor(defaultTimeout time.Duration) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
			defer cancel()
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

package client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/Raisondetr3/checklist-api-service/pkg/pb"

	"github.com/Raisondetr3/checklist-api-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

type TaskClient interface {
	CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error)
	GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error)
	UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error)
	DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error)
	ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error)
	Close() error
}

type taskClient struct {
	client pb.TaskServiceClient
	conn   *grpc.ClientConn
	config config.DBServiceConfig
}

func NewTaskClient(dbConfig config.DBServiceConfig) (TaskClient, error) {
	kacp := keepalive.ClientParameters{
		Time:                dbConfig.KeepAliveTime,
		Timeout:             dbConfig.KeepAliveTimeout,
		PermitWithoutStream: true,
	}

	conn, err := grpc.NewClient(
		dbConfig.GRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithUnaryInterceptor(loggingUnaryInterceptor),
		grpc.WithUnaryInterceptor(retryUnaryInterceptor(dbConfig.MaxRetries, dbConfig.RetryDelay)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", dbConfig.GRPCAddress, err)
	}

	client := pb.NewTaskServiceClient(conn)

	slog.Info("Connected to gRPC server",
		slog.String("address", dbConfig.GRPCAddress),
		slog.Duration("timeout", dbConfig.Timeout),
		slog.Int("max_retries", dbConfig.MaxRetries),
	)

	return &taskClient{
		client: client,
		conn:   conn,
		config: dbConfig,
	}, nil
}

func (c *taskClient) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	ctx = c.addMetadata(ctx, "CreateTask")

	resp, err := c.client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create task failed: %w", err)
	}

	return resp, nil
}

func (c *taskClient) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	ctx = c.addMetadata(ctx, "GetTask")

	resp, err := c.client.GetTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get task failed: %w", err)
	}

	return resp, nil
}

func (c *taskClient) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	ctx = c.addMetadata(ctx, "UpdateTask")

	resp, err := c.client.UpdateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("update task failed: %w", err)
	}

	return resp, nil
}

func (c *taskClient) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	ctx = c.addMetadata(ctx, "DeleteTask")

	resp, err := c.client.DeleteTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("delete task failed: %w", err)
	}

	return resp, nil
}

func (c *taskClient) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	ctx = c.addMetadata(ctx, "ListTasks")

	resp, err := c.client.ListTasks(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list tasks failed: %w", err)
	}

	return resp, nil
}

func (c *taskClient) Close() error {
	slog.Info("Closing gRPC client connection")
	return c.conn.Close()
}

func (c *taskClient) addMetadata(ctx context.Context, method string) context.Context {
	md := metadata.New(map[string]string{
		"client":    "api-service",
		"method":    method,
		"timestamp": time.Now().Format(time.RFC3339),
	})

	if requestID := ctx.Value("request_id"); requestID != nil {
		md.Set("request-id", requestID.(string))
	}

	return metadata.NewOutgoingContext(ctx, md)
}

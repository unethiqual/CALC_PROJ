package client

import (
    "context"
    "log"
    "time"

    pb "github.com/unethiqual/CALC_PROJ/proto/taskservice"
    "google.golang.org/grpc"
)

var grpcConn *grpc.ClientConn
var taskServiceClient pb.TaskServiceClient

func InitGRPCClient() {
    var err error
    grpcConn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    taskServiceClient = pb.NewTaskServiceClient(grpcConn)
}

func GetTask() (*pb.Task, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    resp, err := taskServiceClient.GetTask(ctx, &pb.TaskRequest{})
    if err != nil {
        return nil, err
    }
    return resp.Task, nil
}

func SubmitResult(taskID int64, result float64) error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    _, err := taskServiceClient.SubmitResult(ctx, &pb.TaskResult{
		Id:     taskID,
        Result: result,
    })
    return err
}
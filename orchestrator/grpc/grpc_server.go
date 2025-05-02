package grpc

import (
    "context"
    "log"
    "net"

    "github.com/unethiqual/CALC_PROJ/orchestrator/models"
    "github.com/unethiqual/CALC_PROJ/orchestrator/scheduler"
    pb "github.com/unethiqual/CALC_PROJ/proto/taskservice"
    "google.golang.org/grpc"
)

type TaskServiceServer struct {
    pb.UnimplementedTaskServiceServer
}

func (s *TaskServiceServer) GetTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
    task, err := scheduler.GetNextTask()
    if err != nil {
        return nil, err
    }

    return &pb.TaskResponse{
        Task: &pb.Task{
            Id:            task.ID,
            Arg1:          task.Arg1,
            Arg2:          task.Arg2,
            Operation:     task.Operation,
            OperationTime: int32(task.OperationTime),
        },
    }, nil
}

func (s *TaskServiceServer) SubmitResult(ctx context.Context, req *pb.TaskResult) (*pb.ResultResponse, error) {
    err := scheduler.SubmitTaskResult(req.Id, req.Result)
    if err != nil {
        return &pb.ResultResponse{Status: "error"}, err
    }
    return &pb.ResultResponse{Status: "ok"}, nil
}

func StartGRPCServer() {
    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterTaskServiceServer(grpcServer, &TaskServiceServer{})

    log.Println("gRPC server is running on port 50051")
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
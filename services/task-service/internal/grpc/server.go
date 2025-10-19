package grpc

import (
	"context"
	"time"

	pb "github.com/todo/proto/task"
	"github.com/todo/services/task-service/internal/models"
	"github.com/todo/services/task-service/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
	repo *repository.PostgresRepository
}

func NewTaskServer(repo *repository.PostgresRepository) *TaskServer {
	return &TaskServer{repo: repo}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	priority := convertPriorityFromProto(req.Priority)

	var dueDate *time.Time
	if req.DueDate != nil {
		t := req.DueDate.AsTime()
		dueDate = &t
	}

	task, err := s.repo.CreateTask(req.Title, req.Description, req.UserId, priority, dueDate)
	if err != nil {
		return &pb.CreateTaskResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.CreateTaskResponse{
		Task: convertTaskToProto(task),
	}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	task, err := s.repo.GetTaskByID(req.Id)
	if err != nil {
		return &pb.GetTaskResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.GetTaskResponse{
		Task: convertTaskToProto(task),
	}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	status := convertStatusFromProto(req.Status)
	priority := convertPriorityFromProto(req.Priority)

	var dueDate *time.Time
	if req.DueDate != nil {
		t := req.DueDate.AsTime()
		dueDate = &t
	}

	task, err := s.repo.UpdateTask(req.Id, req.Title, req.Description, status, priority, dueDate)
	if err != nil {
		return &pb.UpdateTaskResponse{
			Error: err.Error(),
		}, nil
	}

	return &pb.UpdateTaskResponse{
		Task: convertTaskToProto(task),
	}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	err := s.repo.DeleteTask(req.Id)
	if err != nil {
		return &pb.DeleteTaskResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.DeleteTaskResponse{
		Success: true,
	}, nil
}

func (s *TaskServer) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	tasks, total, err := s.repo.ListTasks(int(req.Page), int(req.PageSize))
	if err != nil {
		return &pb.ListTasksResponse{
			Error: err.Error(),
		}, nil
	}

	pbTasks := make([]*pb.Task, len(tasks))
	for i, task := range tasks {
		pbTasks[i] = convertTaskToProto(task)
	}

	return &pb.ListTasksResponse{
		Tasks: pbTasks,
		Total: int32(total),
	}, nil
}

func (s *TaskServer) ListUserTasks(ctx context.Context, req *pb.ListUserTasksRequest) (*pb.ListUserTasksResponse, error) {
	status := convertStatusFromProto(req.Status)
	tasks, total, err := s.repo.ListUserTasks(req.UserId, int(req.Page), int(req.PageSize), status)
	if err != nil {
		return &pb.ListUserTasksResponse{
			Error: err.Error(),
		}, nil
	}

	pbTasks := make([]*pb.Task, len(tasks))
	for i, task := range tasks {
		pbTasks[i] = convertTaskToProto(task)
	}

	return &pb.ListUserTasksResponse{
		Tasks: pbTasks,
		Total: int32(total),
	}, nil
}

func convertTaskToProto(task *models.Task) *pb.Task {
	pbTask := &pb.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      convertStatusToProto(task.Status),
		Priority:    convertPriorityToProto(task.Priority),
		UserId:      task.UserID,
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}

	if task.DueDate != nil {
		pbTask.DueDate = timestamppb.New(*task.DueDate)
	}

	return pbTask
}

func convertStatusToProto(status models.TaskStatus) pb.TaskStatus {
	switch status {
	case models.StatusPending:
		return pb.TaskStatus_PENDING
	case models.StatusInProgress:
		return pb.TaskStatus_IN_PROGRESS
	case models.StatusCompleted:
		return pb.TaskStatus_COMPLETED
	case models.StatusCancelled:
		return pb.TaskStatus_CANCELLED
	default:
		return pb.TaskStatus_PENDING
	}
}

func convertStatusFromProto(status pb.TaskStatus) models.TaskStatus {
	switch status {
	case pb.TaskStatus_PENDING:
		return models.StatusPending
	case pb.TaskStatus_IN_PROGRESS:
		return models.StatusInProgress
	case pb.TaskStatus_COMPLETED:
		return models.StatusCompleted
	case pb.TaskStatus_CANCELLED:
		return models.StatusCancelled
	default:
		return models.StatusPending
	}
}

func convertPriorityToProto(priority models.TaskPriority) pb.TaskPriority {
	switch priority {
	case models.PriorityLow:
		return pb.TaskPriority_LOW
	case models.PriorityMedium:
		return pb.TaskPriority_MEDIUM
	case models.PriorityHigh:
		return pb.TaskPriority_HIGH
	case models.PriorityUrgent:
		return pb.TaskPriority_URGENT
	default:
		return pb.TaskPriority_MEDIUM
	}
}

func convertPriorityFromProto(priority pb.TaskPriority) models.TaskPriority {
	switch priority {
	case pb.TaskPriority_LOW:
		return models.PriorityLow
	case pb.TaskPriority_MEDIUM:
		return models.PriorityMedium
	case pb.TaskPriority_HIGH:
		return models.PriorityHigh
	case pb.TaskPriority_URGENT:
		return models.PriorityUrgent
	default:
		return models.PriorityMedium
	}
}

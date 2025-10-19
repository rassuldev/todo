package grpc

import (
	"context"
	"fmt"

	pb "github.com/todo/proto/notification"
	"github.com/todo/services/notification-service/internal/email"
	"github.com/todo/services/notification-service/internal/models"
	"github.com/todo/services/notification-service/internal/push"
	"github.com/todo/services/notification-service/internal/repository"
)

type NotificationServer struct {
	pb.UnimplementedNotificationServiceServer
	repo        *repository.PostgresRepository
	emailSender *email.EmailSender
	pushSender  *push.PushSender
}

func NewNotificationServer(repo *repository.PostgresRepository, emailSender *email.EmailSender, pushSender *push.PushSender) *NotificationServer {
	return &NotificationServer{
		repo:        repo,
		emailSender: emailSender,
		pushSender:  pushSender,
	}
}

func (s *NotificationServer) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	err := s.emailSender.SendEmail(req.To, req.Subject, req.Body)
	if err != nil {
		// Save failed notification
		s.repo.SaveNotification(models.TypeEmail, req.To, req.Subject, req.Body, false)
		return &pb.SendEmailResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Save successful notification
	s.repo.SaveNotification(models.TypeEmail, req.To, req.Subject, req.Body, true)

	return &pb.SendEmailResponse{
		Success: true,
	}, nil
}

func (s *NotificationServer) SendPushNotification(ctx context.Context, req *pb.SendPushNotificationRequest) (*pb.SendPushNotificationResponse, error) {
	err := s.pushSender.SendPushNotification(req.DeviceToken, req.Title, req.Body)
	if err != nil {
		// Save failed notification
		s.repo.SaveNotification(models.TypePush, req.DeviceToken, req.Title, req.Body, false)
		return &pb.SendPushNotificationResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Save successful notification
	s.repo.SaveNotification(models.TypePush, req.DeviceToken, req.Title, req.Body, true)

	return &pb.SendPushNotificationResponse{
		Success: true,
	}, nil
}

func (s *NotificationServer) SendTaskReminder(ctx context.Context, req *pb.SendTaskReminderRequest) (*pb.SendTaskReminderResponse, error) {
	// In a real implementation, you would fetch user details from user service
	// For now, we'll send both email and push notification

	subject := "Task Reminder"
	body := fmt.Sprintf("Don't forget about your task: %s\nDue date: %s",
		req.TaskTitle,
		req.DueDate.AsTime().Format("2006-01-02 15:04:05"))

	// For demonstration, we'll just send email
	// In production, you would get user's email and device token from user service
	userEmail := fmt.Sprintf("user-%s@example.com", req.UserId)

	err := s.emailSender.SendEmail(userEmail, subject, body)
	if err != nil {
		s.repo.SaveNotification(models.TypeEmail, userEmail, subject, body, false)
		return &pb.SendTaskReminderResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	s.repo.SaveNotification(models.TypeEmail, userEmail, subject, body, true)

	return &pb.SendTaskReminderResponse{
		Success: true,
	}, nil
}

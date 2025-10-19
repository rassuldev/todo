package push

import (
	"fmt"
	"log"
)

type PushSender struct {
	apiKey string
}

func NewPushSender(apiKey string) *PushSender {
	return &PushSender{
		apiKey: apiKey,
	}
}

func (s *PushSender) SendPushNotification(deviceToken, title, body string) error {
	// For demonstration purposes, we'll just log the notification
	// In production, you would use FCM (Firebase Cloud Messaging) or similar service
	log.Printf("Sending push notification to device: %s, title: %s, body: %s", deviceToken, title, body)

	if s.apiKey == "" {
		log.Println("Push notification API key not configured (simulation mode)")
		return nil
	}

	// Here you would implement actual push notification logic using FCM or similar
	// For now, we'll just simulate it
	log.Printf("Push notification sent successfully to device %s", deviceToken)
	return nil
}

func (s *PushSender) SendTaskReminder(deviceToken, taskTitle string) error {
	title := "Task Reminder"
	body := fmt.Sprintf("Don't forget about your task: %s", taskTitle)
	return s.SendPushNotification(deviceToken, title, body)
}

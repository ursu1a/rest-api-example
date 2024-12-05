package emails

import (
	"backend/utils"
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SenderEmail  string
}

// Constructor for EmailService
func InitEmailService() *EmailService {
	if err := utils.CheckEnvs([]string{"SMTP_PORT", "SMTP_HOST", "SMTP_EMAIL", "SMTP_PASSWORD", "SMTP_EMAIL_SENDER"}); err != nil {
		log.Fatalf("Error checking environment variables: %v", err)
		return nil
	}

	port, _ := strconv.Atoi(utils.GetEnv("SMTP_PORT", ""))
	return &EmailService{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     port,
		SMTPUser:     os.Getenv("SMTP_EMAIL"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SenderEmail:  os.Getenv("SMTP_EMAIL_SENDER"),
	}
}

// Common send email function
func (s *EmailService) sendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.SenderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.SMTPHost, s.SMTPPort, s.SMTPUser, s.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}
	return nil
}

// Send registration confirmation email
func (s *EmailService) SendRegistrationConfirmation(to string, token string) error {
	frontAddress := utils.GetEnv("FRONTEND_ADDRESS", "http://localhost:3001")
	link := fmt.Sprintf("%s/verify-email?token=%v", frontAddress, token)

	body := fmt.Sprintf(`
		<h1>Please confirm registration</h1>
		<p>Thanks for registration! Go to the link below to confirm your registration:</p>
		<a href="%s">Confirm Email</a>
	`, link)

	return s.sendEmail(to, "Registration confirmation", body)
}

// Reset password email
func (s *EmailService) SendPasswordReset(to string, token string) error {
	frontAddress := utils.GetEnv("FRONTEND_ADDRESS", "http://localhost:3001")
	link := fmt.Sprintf("%s/update-password?token=%v", frontAddress, token)

	body := fmt.Sprintf(`
		<h1>Password reset</h1>
		<p>We received a request for reset password. Please go to the link below for set new password:</p>
		<a href="%s">Reset password</a>
	`, link)
	return s.sendEmail(to, "Reset password", body)
}

// Transactional email
func (s *EmailService) SendTransactionalEmail(to string, subject string, content string) error {
	body := fmt.Sprintf(`
		<h1>%s</h1>
		<p>%s</p>
	`, subject, content)
	return s.sendEmail(to, subject, body)
}

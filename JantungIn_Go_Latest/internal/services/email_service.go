package services

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"jantungin-api-server/pkg/utils"
)

type EmailService interface {
	SendOTP(ctx context.Context, toEmail, toName, code string, ttl time.Duration) error
}

type smtpEmailService struct {
	cfg *utils.Config
}

func NewEmailService(cfg *utils.Config) EmailService {
	return &smtpEmailService{cfg: cfg}
}

func (s *smtpEmailService) SendOTP(ctx context.Context, toEmail, toName, code string, ttl time.Duration) error {
	_ = ctx

	if strings.TrimSpace(toEmail) == "" {
		return fmt.Errorf("email penerima kosong")
	}
	if s.cfg.SMTP.Host == "" || s.cfg.SMTP.Port == 0 || s.cfg.SMTP.Username == "" || s.cfg.SMTP.Password == "" {
		return fmt.Errorf("smtp belum dikonfigurasi")
	}

	fromEmail := s.cfg.SMTP.FromEmail
	if fromEmail == "" {
		fromEmail = s.cfg.SMTP.Username
	}

	fromName := s.cfg.SMTP.FromName
	if fromName == "" {
		fromName = "JantungIn"
	}

	if strings.TrimSpace(toName) == "" {
		toName = toEmail
	}

	subject := "Kode OTP JantungIn"
	body := fmt.Sprintf(
		"Halo %s,\n\nKode OTP kamu: %s\nBerlaku selama %d detik.\n\nJika bukan kamu, abaikan email ini.\n",
		toName,
		code,
		int(ttl.Seconds()),
	)

	headerFrom := fmt.Sprintf("%s <%s>", fromName, fromEmail)
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\r\n", headerFrom))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", toEmail))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTP.Host, s.cfg.SMTP.Port)
	auth := smtp.PlainAuth("", s.cfg.SMTP.Username, s.cfg.SMTP.Password, s.cfg.SMTP.Host)

	return smtp.SendMail(addr, auth, fromEmail, []string{toEmail}, []byte(msg.String()))
}

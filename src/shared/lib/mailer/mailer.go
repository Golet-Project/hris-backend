package mailer

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

var ErrorToEmailMissing = errors.New("missing to email")

type mailer struct {
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string

	senderName string
	toEmail   []string
	subject   string
}

func NewMailer() *mailer {
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	return &mailer{
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     smtpPort,
		smtpUser: os.Getenv("SMTP_AUTH_USER"),
		smtpPassword: os.Getenv("SMTP_AUTH_PASSWORD"),
	}
}

func (m *mailer) SenderName(name string) *mailer {
	m.senderName = name

	return m
}

func (m *mailer) To(email []string) *mailer {
	m.toEmail = email

	return m
}

func (m *mailer) Subject(subject string) *mailer {
	m.subject = subject

	return m
}

func (m *mailer) Send(body string) error {
	if len(m.toEmail) == 0 {
		return ErrorToEmailMissing
	}

	auth := smtp.PlainAuth("", m.smtpUser, m.smtpPassword, m.smtpHost)

	smtpAddr := fmt.Sprintf("%s:%d", m.smtpHost, m.smtpPort)

	builder := strings.Builder{}

	builder.WriteString("From: ")

	var senderName string
	if m.senderName == "" {
		senderName = os.Getenv("SMTP_SENDER_NAME")
	} else {
		senderName = m.senderName
	}

	builder.WriteString(senderName + "\n")
	builder.WriteString("To: ")
	builder.WriteString(strings.Join(m.toEmail, ", ") + "\n")

	if m.subject != "" {
		builder.WriteString("Subject: ")
		builder.WriteString(m.subject + "\n")
	}

	builder.WriteString(body)

	msg := builder.String()

	err := smtp.SendMail(smtpAddr, auth, m.smtpUser, m.toEmail, []byte(msg))

	return err
}

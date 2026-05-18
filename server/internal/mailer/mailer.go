package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"github.com/chiagxziem/snipper/internal/config"
	"github.com/resend/resend-go/v3"
)

//go:embed templates/*.html
var templates embed.FS

type Mailer struct {
	mailer *resend.Client
	domain string
}

func New(cfg *config.Config) *Mailer {
	return &Mailer{
		mailer: resend.NewClient(cfg.ResendAPIKey),
		domain: cfg.ResendDomain,
	}
}

func (m *Mailer) SendEmail(
	ctx context.Context,
	to []string,
	subject, html string,
) error {
	params := &resend.SendEmailRequest{
		To:      to,
		From:    getFrom(m.domain),
		Html:    html,
		Subject: subject,
	}

	_, err := m.mailer.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func getFrom(domain string) string {
	return fmt.Sprintf("Snipper <snipper@%s>", domain)
}

func (m *Mailer) SendVerificationEmail(
	ctx context.Context,
	to []string,
	name, token string,
) error {
	tmpl, err := template.ParseFS(templates, "templates/verification.html")
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, map[string]string{
		"Name":      name,
		"VerifyURL": "http://localhost:3000/verify?token=" + token,
	})
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return m.SendEmail(ctx, to, "Verify your email", body.String())
}

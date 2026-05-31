package mailer

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/chiagxziem/gater/internal/config"
	"github.com/resend/resend-go/v3"
)

type resendClient struct {
	mailer      *resend.Client
	domain      string
	frontendURL string
}

func NewResendClient(cfg *config.Config) *resendClient {
	return &resendClient{
		mailer:      resend.NewClient(cfg.ResendAPIKey),
		domain:      cfg.ResendDomain,
		frontendURL: cfg.FrontendURL,
	}
}

func (r *resendClient) SendEmail(
	ctx context.Context,
	to []string,
	subject, html string,
) error {
	params := &resend.SendEmailRequest{
		To:      to,
		From:    getFrom(r.domain),
		Html:    html,
		Subject: subject,
	}

	_, err := r.mailer.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func (r *resendClient) SendVerificationEmail(
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
		"VerifyURL": r.frontendURL + "/verify?token=" + token,
	})
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return r.SendEmail(ctx, to, "Verify your email", body.String())
}

func (r *resendClient) SendPasswordResetEmail(
	ctx context.Context,
	to []string,
	name, token string,
) error {
	tmpl, err := template.ParseFS(templates, "templates/password-reset.html")
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, map[string]string{
		"Name":     name,
		"ResetURL": r.frontendURL + "/password-reset?token=" + token,
	})
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return r.SendEmail(ctx, to, "Reset your password", body.String())
}

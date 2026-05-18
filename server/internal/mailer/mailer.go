package mailer

import (
	"context"
	"embed"
	"fmt"
)

//go:embed templates/*.html
var templates embed.FS

type Mailer interface {
	SendEmail(ctx context.Context, to []string, subject, html string) error
	SendVerificationEmail(ctx context.Context, to []string, name, token string) error
}

func getFrom(domain string) string {
	return fmt.Sprintf("Snipper <snipper@%s>", domain)
}

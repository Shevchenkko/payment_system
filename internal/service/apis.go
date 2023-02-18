package service

import "context"

// APIs contains all available APIs.
type APIs struct {
	Emails EmailsAPI
}

// EmailsAPI - represents emails api.
type EmailsAPI interface {
	SendEmail(ctx context.Context, inp SendEmailInput) error
}

type SendEmailInput struct {
	To          string
	Subject     string
	ContentType string
	Body        string
}

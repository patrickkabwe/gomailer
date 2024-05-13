package gomailer

import (
	"net/mail"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var u = mail.Address{Name: "Sender Name", Address: os.Getenv("EMAIL_USER")}

var to = []string{os.Getenv("EMAIL_USER_RECEIVER")}
var host = "smtp.gmail.com"
var port = 587
var password = os.Getenv("EMAIL_PASSWORD")

func TestGoMailer_New(t *testing.T) {
	testCases := []struct {
		name string
		opts GoMailerOption
	}{
		{
			name: "GoMailer",
			opts: GoMailerOption{
				Host:     host,
				Port:     port,
				Username: u.Address,
				Password: password,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mailer := New(tc.opts)
			assert.NotNil(t, mailer)
			assert.IsType(t, &goMailer{}, mailer)
			assert.Equal(t, tc.opts.Host, mailer.(*goMailer).host)
			assert.Equal(t, tc.opts.Port, mailer.(*goMailer).port)
		})
	}
}

func TestGoMailer_SendMail(t *testing.T) {
	testCases := []struct {
		name       string
		message    EmailMessage
		shouldFail bool
	}{
		{
			name: "Should send mail with text body",
			message: EmailMessage{
				From:    u.String(),
				To:      to,
				Subject: "Text Email",
				Body:    []byte("Hello Test"),
				CC:      []string{u.Address},
			},
			shouldFail: false,
		},
		{
			name: "Should send mail with attachment and text body",
			message: EmailMessage{
				From:    u.String(),
				To:      to,
				Subject: "Email with attachment",
				Body:    []byte("Hello Test"),
				CC:      []string{u.Address},
				Attachments: []Attachment{
					{
						Name: "attachment.txt",
						Path: "testdata/attachments/attachment.txt",
					},
				},
			},
			shouldFail: false,
		},
		{
			name: "Should send mail with template",
			message: EmailMessage{
				From:    u.String(),
				To:      to,
				Body:    []byte("Hello Test"),
				Subject: "Email with template and attachment",
				Attachments: []Attachment{
					{
						Name: "attachment.txt",
						Path: "testdata/attachments/attachment.txt",
					},
				},
				Template: Template{
					Path: "testdata/template/email.html",
					Data: struct {
						Name string
					}{
						Name: "Patrick",
					},
				},
			},
			shouldFail: false,
		},
		{
			name: "Should send mail with template and attachment",
			message: EmailMessage{
				From:    u.String(),
				To:      to,
				Body:    []byte("Hello Test"),
				Subject: "Email with template",
				CC:      []string{u.Address},
				Template: Template{
					Path: "testdata/template/email.html",
					Data: struct {
						Name string
					}{
						Name: "Patrick",
					},
				},
			},
			shouldFail: false,
		},
		{
			name: "Should not send mail when from is empty",
			message: EmailMessage{
				To:      to,
				Subject: "Test Mail",
				Body:    []byte("Hello Test"),
			},
			shouldFail: true,
		},
		{
			name: "Should not send mail when to is empty",
			message: EmailMessage{
				Subject: "Test Mail",
				Body:    []byte("Hello Test"),
			},
			shouldFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mailer := New(GoMailerOption{
				Host:     host,
				Port:     port,
				Username: u.Address,
				Password: password,
			})

			err := mailer.SendMail(tc.message)
			if tc.shouldFail && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

package gomailer

import (
	"net/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

var u = mail.Address{Name: "Sender Name", Address: "patrickckabwe@gmail.com"}

var to = []string{"patrickkabwe45@yahoo.com"}
var host = "smtp.gmail.com"
var port = 587
var password = "fqjfeqttraieafkt"

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
		opts       EmailOption
		shouldFail bool
	}{
		{
			name: "Should send mail",
			message: EmailMessage{
				From:    u.Address,
				To:      to,
				Subject: "Test Mail",
				Body:    []byte("Hello Test"),
			},
			opts:       EmailOption{},
			shouldFail: false,
		},
		{
			name: "Should not send mail when from is empty",
			message: EmailMessage{
				To:      to,
				Subject: "Test Mail",
				Body:    []byte("Hello Test"),
			},
			opts:       EmailOption{},
			shouldFail: true,
		},
		{
			name: "Should not send mail when to is empty",
			message: EmailMessage{
				Subject: "Test Mail",
				Body:    []byte("Hello Test"),
			},
			opts:       EmailOption{},
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

			err := mailer.SendMail(tc.message, tc.opts)
			if tc.shouldFail && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGoMailer_SendMailWithTemplate(t *testing.T) {
	testCases := []struct {
		name       string
		message    EmailMessage
		opts       EmailOption
		shouldFail bool
	}{
		{
			name: "Should not send mail with template",
			message: EmailMessage{
				From:    u.Address,
				To:      to,
				Subject: "Test Mail",
				Body:    []byte("Hello Test"),
			},
			opts:       EmailOption{},
			shouldFail: true,
		},
		{
			name: "Should send mail with template",
			message: EmailMessage{
				From:    u.Address,
				To:      to,
				Subject: "Test Mail",
				Body:    []byte("Hello Test"),
			},
			opts: EmailOption{
				TemPath: "testdata/template/test.html",
				TemData: struct{ Name string }{Name: "Patrick"},
			},
			shouldFail: false,
		},
		{
			name: "Should not send mail with template when template path is invalid",
			message: EmailMessage{
				From:    u.Address,
				To:      to,
				Subject: "Test Mail",
				Body:    []byte("Hello Test"),
			},
			opts: EmailOption{
				TemPath: "testdata/template/test.html",
				TemData: struct{ Name string }{Name: "Patrick"},
			},
			shouldFail: true,
		},
		{
			name: "Should not send mail with template when template path is invalid",
			message: EmailMessage{
				// From:    "",
				// To:      to,
				Subject: "Test Mail",
				Body:    []byte("Hello Test"),
			},
			opts: EmailOption{
				TemPath: "template/test.html",
				TemData: struct{ Name string }{Name: "Patrick"},
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

			err := mailer.SendMailWithTemplate(tc.message.From, tc.message.To, tc.message.Subject, tc.opts)
			if tc.shouldFail && err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

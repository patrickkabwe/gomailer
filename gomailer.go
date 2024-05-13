package gomailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
)

type GoMailerOption struct {
	// host to use when sending the email
	Host string
	// port to use when sending the email
	Port int
	// username to use when authenticating to the smtp server
	Username string
	// password to use when authenticating to the smtp server
	Password string
}

type EmailMessage struct {
	// From is the email address of the sender of the email
	From string
	// To is a slice of email addresses to send the email to
	To []string
	// Subject is the subject of the email
	Subject string
	// Body can be a string (text) or (html)
	Body []byte
}

// EmailOption is used to pass additional options to the email
type EmailOption struct {
	TemPath     string
	TemData     interface{}
	ContentType string // text/html or text/plain
	Secure      bool
	Attachments []Attachment
}

// GoMailer is the interface that wraps the basic SendMail method for sending emails
type GoMailer interface {
	// SendMail sends an email using the provided message and options
	SendMail(message EmailMessage, options EmailOption) error
	// SendMailWithTemplate sends an email with html template provided
	SendMailWithTemplate(from string, to []string, subject string, options EmailOption) error
}

// goMailer is the default implementation of the GoMailer interface
type goMailer struct {
	host     string
	port     int
	username string
	password string
}

// New returns a new GoMailer instance with the provided options
func New(options GoMailerOption) GoMailer {
	return &goMailer{
		host:     options.Host,
		port:     options.Port,
		username: options.Username,
		password: options.Password,
	}
}

func (gm *goMailer) SendMail(message EmailMessage, options EmailOption) error {
	if message.From == "" {
		return fmt.Errorf("from address is required")
	}

	auth := gm.auth()
	fullHost := fmt.Sprintf("%s:%d", gm.host, gm.port)
	msg := gm.preparePayload(message.From, message.To, message.Subject, message.Body, "text/html")

	err := smtp.SendMail(fullHost, auth, message.From, message.To, msg)
	if err != nil {
		return err
	}

	return nil
}

func (gm *goMailer) SendMailWithTemplate(from string, to []string, subject string, options EmailOption) error {
	if options.TemPath == "" {
		return fmt.Errorf("template path is required")
	}

	mgs, err := gm.parseHtmlTemplate(options.TemPath, options.TemData)
	if err != nil {
		return err
	}

	d := EmailMessage{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    mgs,
	}

	err = gm.SendMail(d, options)

	if err != nil {
		return err
	}

	return nil
}

func (gm *goMailer) auth() smtp.Auth {
	auth := smtp.PlainAuth("", gm.username, gm.password, gm.host)
	return auth
}

func (gm *goMailer) preparePayload(from string, to []string, subject string, body []byte, contentType string) []byte {
	headers := make(map[string]string)
	headers["MIME-version"] = "1.0"
	headers["From"] = from
	headers["To"] = strings.Join(to, ", ")
	headers["Subject"] = subject
	headers["Content-Type"] = fmt.Sprintf("%s; charset=utf-8", contentType)
	headers["Content-Transfer-Encoding"] = "quoted-printable"
	headers["Content-Disposition"] = "inline"
	headers["X-Mailer"] = "Go-Mailer"

	msg := gm.headersToBytes(headers, body)
	return msg
}

func (gm *goMailer) parseHtmlTemplate(path string, data interface{}) ([]byte, error) {
	message := new(bytes.Buffer)
	t, err := template.ParseFiles(path)

	if err != nil {
		return nil, err
	}
	if err := t.Execute(message, data); err != nil {

		return nil, err
	}

	return message.Bytes(), nil
}

func (gm *goMailer) headersToBytes(headers map[string]string, message []byte) []byte {
	var h []string
	for k, v := range headers {
		h = append(h, fmt.Sprintf("%s: %s \r\n", k, v))
	}
	h = append(h, string(message))

	msg := []byte(strings.Join(h, ""))
	gm.msg.WriteString("\r\n")

	return msg
}

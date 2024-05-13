package gomailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"net"
	"net/mail"
	"net/smtp"
	"regexp"
	"strings"
	"time"
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
	// secure indicates whether to use a secure connection when sending the email (e.g. SSL/TLS)
	Secure bool
}

type Template struct {
	// Path is the path to the template file
	Path string
	// Data is the data to use when rendering the template
	Data interface{}
}

type Attachment struct {
	// Name is the name of the attachment
	Name string
	// Path is the path to the attachment file
	Path string
}
type EmailMessage struct {
	// Name is the name of the sender of the email
	Name string
	// From is the email address of the sender of the email
	From string
	// To is a slice of email addresses to send the email to
	To []string
	// Subject is the subject of the email
	Subject string
	// Body can be a string (text) or (html)
	Body []byte
	// Attachments is a slice of file paths to attach to the email
	Attachments []Attachment
	// CC is a slice of email addresses to send a copy of the email to
	CC []string
	// BCC is a slice of email addresses to send a blind copy of the email to
	BCC []string
	// ReplyTo is the email address to use when replying to the email
	ReplyTo string
	// Template is the path to the template file to use when sending the email
	Template Template
}

// GoMailer is the interface that wraps the basic SendMail method for sending emails
type GoMailer interface {
	// SendMail sends an email using the provided message and options
	SendMail(message EmailMessage) error
}

// goMailer is the default implementation of the GoMailer interface
type goMailer struct {
	host     string
	port     int
	username string
	password string
	secure   bool
	msg      *bytes.Buffer
	writer   *multipart.Writer
}

// New returns a new GoMailer instance with the provided options
func New(options GoMailerOption) GoMailer {
	msg := new(bytes.Buffer)
	writer := multipart.NewWriter(msg)
	return &goMailer{
		host:     options.Host,
		port:     options.Port,
		username: options.Username,
		password: options.Password,
		secure:   options.Secure,
		msg:      msg,
		writer:   writer,
	}
}

func (gm *goMailer) SendMail(message EmailMessage) error {
	if message.From == "" {
		message.From = gm.username
	}
	if message.Template.Path != "" {
		if message.Template.Data == nil {
			return fmt.Errorf("TemplateData is required when using a template")
		}
		mgs, err := gm.parseHtmlTemplate(message.Template.Path, message.Template.Data)
		if err != nil {
			return err
		}
		message.Body = mgs
	}

	if message.Name != "" {
		add := mail.Address{Name: message.Name, Address: message.From}
		message.From = add.String()
	}

	return gm.sendMail(message)
}

func (gm *goMailer) sendMail(message EmailMessage) error {
	auth := gm.auth()
	fullHost := net.JoinHostPort(gm.host, fmt.Sprintf("%d", gm.port))

	msg := gm.createMessage(message)
	senderInfos := strings.Split(message.From, " ")
	re := regexp.MustCompile("[<>]")
	from := re.ReplaceAllString(senderInfos[len(senderInfos)-1], "")
	err := smtp.SendMail(fullHost, auth, from, message.To, msg)

	if err != nil {
		log.Println("Error sending email: ", err)
		return err
	}

	return nil
}

func (gm *goMailer) auth() smtp.Auth {
	auth := smtp.PlainAuth("", gm.username, gm.password, gm.host)
	return auth
}

func (gm *goMailer) createMessage(message EmailMessage) []byte {

	headers := make(map[string]string)
	headers["MIME-version"] = "1.0"
	headers["From"] = message.From
	headers["To"] = strings.Join(message.To, ", ")
	headers["Subject"] = message.Subject
	headers["Content-Type"] = "multipart/mixed; boundary=" + gm.writer.Boundary()
	headers["Content-Transfer-Encoding"] = "quoted-printable"
	headers["X-Mailer"] = createXMailer(gm.host)
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	headers["Message-ID"] = fmt.Sprintf("<%d.%d.%d@%s>", time.Now().Unix(), time.Now().Nanosecond(), time.Now().UnixNano(), gm.host)
	headers["List-Id"] = message.From
	headers["Reply-To"] = message.From

	if len(message.CC) > 0 {
		headers["Cc"] = strings.Join(message.CC, ", ")
	}
	if len(message.BCC) > 0 {
		headers["Bcc"] = strings.Join(message.BCC, ", ")
	}
	if message.ReplyTo != "" {
		headers["Reply-To"] = message.ReplyTo
	}

	gm.headersToBytes(headers, message)
	return gm.msg.Bytes()
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

func (gm *goMailer) headersToBytes(headers map[string]string, message EmailMessage) {
	var contentType string
	for key, value := range headers {
		gm.msg.WriteString(key + ": " + value + "\r\n")
	}

	gm.msg.WriteString("\r\n")

	gm.msg.WriteString("\r\n--" + gm.writer.Boundary() + "\r\n")
	hasHTML := strings.Contains(string(message.Body), "html")
	if hasHTML || (message.Template.Path != "" && hasHTML) {
		contentType = "text/html"
	} else {
		contentType = "text/plain"
	}
	gm.msg.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n", contentType))
	gm.msg.WriteString("\r\n")
	gm.msg.WriteString(string(message.Body))

	gm.msg.WriteString("\r\n")

	for _, attachment := range message.Attachments {
		attachFile(gm.writer, gm.msg, attachment)
	}

	gm.msg.WriteString("\r\n--" + gm.writer.Boundary() + "--")
}

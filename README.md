## Gomailer

Gomailer is a Go programming language package that provides functionality for sending emails. It allows Go developers to send emails from within their Go applications, making it easy to integrate email functionality into their software.

### Installation

To use Gomailer, you need to have Go installed. Once you have Go installed, you can use the following command to install Gomailer:

```golang
go get github.com/Kazion500/gomailer
```

### Usage

Here's an example of how to use Gomailer to send an email:

```golang
package main

import (
    "github.com/Kazion500/gomailer"
)

func main() {
    options := gomailer.GoMailerOption{
        Host:     "smtp.gmail.com",
        Port:     587,
        Username: "your_username",
        Password: "your_password",
    }

    gm := gomailer.New(options)

    message := gomailer.EmailMessage{
        From:    "your_email_address",
        To:      []string{"recipient_email_address"},
        Subject: "Test email from Gomailer",
        Body:    []byte("This is a test email sent from Gomailer."),
    }

    emailOptions := gomailer.EmailOption{
        ContentType: "text/plain",
        Secure:      true,
    }

    err := gm.SendMail(message, emailOptions)
    if err != nil {
        // Handle error
    }
}
```

### API

#### `GoMailerOption`

`Host` - host to use when sending the email
`Port` - port to use when sending the email
`Username` - username to use when authenticating to the SMTP server
`Password` - password to use when authenticating to the SMTP server

#### `EmailMessage`

`From` - the email address of the sender of the email
`To` - a slice of email addresses to send the email to
`Subject` - the subject of the email
`Body` - can be a string (text) or (html)

#### `EmailOption`

`TemPath` - path to the email template
`TemData` - data to be passed to the email template
`ContentType` - type of email content (text/html or text/plain)
`Secure` - whether to use secure connection when sending the email

#### `New(options GoMailerOption) GoMailer`

Returns a new GoMailer instance with the provided options.

#### `SendMail(message EmailMessage, options EmailOption) error`

Sends an email using the provided message and options.

#### `SendMailWithTemplate(from string, to []string, subject string, options EmailOption) error`

Sends an email with the provided email template.

### License

Gomailer is released under the MIT License. See [LICENSE](https://github.com/Kazion500/gomailer/blob/main/LICENSE) for details.

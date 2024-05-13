# Gomailer 

Gomailer is a simple Go library for sending emails from Go applications. It supports sending text or HTML emails, email templates, email attachments and many more.

### Installation 🛠️

To use Gomailer, you need to have Go installed. Once you have Go installed, you can use the following command to install Gomailer:

```golang
go get github.com/patrickkabwe/gomailer
```

### Features 🚀

- Send emails from Go applications
- Text or HTML email content supported
- Email templates
- Email Attachments - (WIP)
- Secure connections(SSL/TLS) - (WIP)


### Usage 📝

Here's an example of how to use Gomailer to send an email:

```golang
package main

import (
    "github.com/patrickkabwe/gomailer"
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

    err := gm.SendMail(message)
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

#### `New(options GoMailerOption) GoMailer`

Returns a new GoMailer instance with the provided options.

#### `SendMail(message EmailMessage) error`
Sends an email using the provided message and options.


### License

Gomailer is released under the MIT License. See [LICENSE](https://github.com/patrickkabwe/gomailer/blob/main/LICENSE) for details.

package gomailer

import (
	"bytes"
	"log"
	"mime"
	"mime/multipart"
	"net/textproto"
	"os"
)

func createXMailer(host string) string {
	switch host {
	case "gmail.com":
		return "Google Gmail"
	case "yahoo.com":
		return "Yahoo Mail"
	case "outlook.com":
		return "Microsoft Outlook"
	default:
		return "Golang"
	}
}

func attachFile(writer *multipart.Writer, msg *bytes.Buffer, attachment Attachment) {
	file, err := os.ReadFile(attachment.Path)
	if err != nil {
		log.Fatalf("Failed to open attachment file %s: %v", attachment.Path, err)
		return
	}

	part := make(textproto.MIMEHeader)
	part.Set("Content-Type", mime.TypeByExtension(attachment.Path))
	part.Set("Content-Disposition", `attachment; filename="`+attachment.Name+`"`)
	writer.CreatePart(part)
	msg.Write(file)
}

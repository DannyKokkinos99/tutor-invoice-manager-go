package gmail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
)

// GmailService handles Gmail API interactions
type EmailSender struct {
	SMTPServer string
	Port       string
	Username   string
	Password   string
}

func NewEmailSender(username, password string) *EmailSender {
	return &EmailSender{
		SMTPServer: "smtp.gmail.com",
		Port:       "587",
		Username:   username,
		Password:   password,
	}
}

// SendEmail sends an email using SMTP
func (e *EmailSender) SendEmail(to, subject, bodyText, attachmentFileName string) error {
	// SMTP Authentication
	auth := smtp.PlainAuth("", e.Username, e.Password, e.SMTPServer)

	// Create a buffer to hold the email
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the HTML part to the email body
	htmlPart, err := writer.CreatePart(
		textproto.MIMEHeader{
			"Content-Type":              []string{"text/html; charset=UTF-8"},
			"Content-Transfer-Encoding": []string{"quoted-printable"},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create HTML part: %v", err)
	}

	// Encode the HTML body text in quoted-printable format
	quotedprintable.NewWriter(htmlPart).Write([]byte(bodyText))

	// Open the attachment file (make sure to change the path as needed)
	attachment, err := os.Open(attachmentFileName)
	if err != nil {
		return fmt.Errorf("failed to open attachment file: %v", err)
	}
	defer attachment.Close()

	// Read the attachment content
	attachmentData, err := io.ReadAll(attachment)
	if err != nil {
		return fmt.Errorf("failed to read attachment: %v", err)
	}

	// Encode the attachment in base64
	encodedAttachment := base64.StdEncoding.EncodeToString(attachmentData)

	// Create the attachment part
	attachmentName := filepath.Base(attachmentFileName)
	attachmentPart, err := writer.CreatePart(
		textproto.MIMEHeader{
			"Content-Type":              []string{"application/pdf"},
			"Content-Disposition":       []string{fmt.Sprintf(`attachment; filename="%s"`, attachmentName)},
			"Content-Transfer-Encoding": []string{"base64"},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create attachment part: %v", err)
	}

	// Write the base64-encoded attachment data
	attachmentPart.Write([]byte(encodedAttachment))

	// Finalize the email
	writer.Close()

	// Compose the email with necessary headers
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n", e.Username, to, subject, writer.Boundary())
	message += buf.String()

	// Send the email using the SMTP protocol
	err = smtp.SendMail(e.SMTPServer+":"+e.Port, auth, e.Username, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully!")
	return nil
}

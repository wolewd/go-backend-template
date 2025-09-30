package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"
	"sync"
)

var (
	emailConfig     *EmailConfig
	emailConfigOnce sync.Once
	emailInitErr    error
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

func getEmailConfig() (*EmailConfig, error) {
	emailConfigOnce.Do(func() {
		host := GetEnv("SMTP_HOST", "")
		port := GetEnv("SMTP_PORT", "587")
		username := GetEnv("SMTP_USERNAME", "")
		password := GetEnv("SMTP_PASSWORD", "")
		fromEmail := GetEnv("SMTP_FROM_EMAIL", "")
		fromName := GetEnv("SMTP_FROM_NAME", "No Reply")

		if host == "" || username == "" || password == "" || fromEmail == "" {
			emailInitErr = fmt.Errorf("SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD, and SMTP_FROM_EMAIL must be set")
			return
		}

		emailConfig = &EmailConfig{
			SMTPHost:     host,
			SMTPPort:     port,
			SMTPUsername: username,
			SMTPPassword: password,
			FromEmail:    fromEmail,
			FromName:     fromName,
		}

		log.Println("Email config initialized successfully")
	})

	if emailInitErr != nil {
		return nil, emailInitErr
	}

	return emailConfig, nil
}

func SendEmail(to, subject, templatePath string, data interface{}) error {
	config, err := getEmailConfig()
	if err != nil {
		return fmt.Errorf("failed to get email config: %w", err)
	}

	recipients := parseEmails(to)
	if len(recipients) == 0 {
		return fmt.Errorf("no valid recipients provided")
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	from := fmt.Sprintf("%s <%s>", config.FromName, config.FromEmail)
	message := "From: " + from + "\r\n"
	message += "To: " + strings.Join(recipients, ", ") + "\r\n"
	message += "Subject: " + subject + "\r\n"
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body.String()

	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)

	if err := smtp.SendMail(addr, auth, config.FromEmail, recipients, []byte(message)); err != nil {
		return fmt.Errorf("failed to send email (to=%s, subject=%s): %w", to, subject, err)
	}

	return nil
}

func SendEmailWithCC(to, cc, subject, templatePath string, data interface{}) error {
	config, err := getEmailConfig()
	if err != nil {
		return fmt.Errorf("failed to get email config: %w", err)
	}

	toRecipients := parseEmails(to)
	ccRecipients := parseEmails(cc)

	if len(toRecipients) == 0 {
		return fmt.Errorf("no valid recipients provided")
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	from := fmt.Sprintf("%s <%s>", config.FromName, config.FromEmail)
	message := "From: " + from + "\r\n"
	message += "To: " + strings.Join(toRecipients, ", ") + "\r\n"

	if len(ccRecipients) > 0 {
		message += "Cc: " + strings.Join(ccRecipients, ", ") + "\r\n"
	}

	message += "Subject: " + subject + "\r\n"
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body.String()

	allRecipients := append(toRecipients, ccRecipients...)

	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)

	if err := smtp.SendMail(addr, auth, config.FromEmail, allRecipients, []byte(message)); err != nil {
		return fmt.Errorf("failed to send email (to=%s, cc=%s, subject=%s): %w", to, cc, subject, err)
	}

	return nil
}

func parseEmails(emails string) []string {
	if emails == "" {
		return []string{}
	}

	parts := strings.Split(emails, ",")
	result := make([]string, 0, len(parts))

	for _, email := range parts {
		email = strings.TrimSpace(email)
		if email != "" {
			result = append(result, email)
		}
	}

	return result
}

/* EXAMPLE
// Define your email data struct
type VerifyEmailData struct {
    UserName    string
    VerifyURL   string
    ExpiredTime string
}

// Usage using template/email/verify.html
func SendVerificationEmail(userEmail, userName, token string) error {
    // Create the struct
    data := VerifyEmailData{
        UserName:    userName,
        VerifyURL:   fmt.Sprintf("https://example.com/verify?token=%s", token),
        ExpiredTime: "24 hours",
    }

    // Send email with struct
    return utils.SendEmail(
        userEmail,
        "Verify Your Account",
        "templates/emails/verify.html",
        data,
    )

	// Send email with CC
	return utils.SendEmailWithCC(
	    userEmail,
        "admin@example.com,support@example.com",
        "Verify Your Account",
        "templates/emails/verify.html",
        data,
	)
}
*/

/* ANOTHER Example:
// Adding new email template in templates/email called annoucement.html
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Announcement</title>
</head>
<body>
<h2>Hello {{.TeamName}},</h2>
<p>We are excited to announce: <b>{{.Message}}</b></p>
<p>Date: {{.Date}}</p>
<p>regards, Admin Team</p>
</body>
</html>

// Usage
type AnnouncementData struct {
	TeamName string
	Message  string
	Date     string
}

err := utils.SendEmailWithCC(
	"user1@example.com,user2@example.com,user3@example.com",
	"manager@example.com,hr@example.com",
	"Important Announcement",
	"templates/emails/announcement.html",
	AnnouncementData{
		TeamName: "Development Team",
		Message:  "Version 2.0 Release",
		Date:     "30 September 2025",
	},
)
if err != nil {
	log.Printf("failed to send announcement email: %v", err)
} else {
	log.Println("announcement email sent successfully")
}
*/

package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"password-lock/templates"
)

func (s Service) SendPasswordEmail(emailAddress string, password string) error {

	subject := "Password Lock password"

	body, err := ParseTemplate("password_email.html", struct {
		Password string
	}{Password: password})
	if err != nil {
		return err
	}

	err = s.sendEmail(emailAddress, subject, body)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")
	return nil
}

func (s Service) SendVerificationLinkEmail(emailAddress string, token string) error {
	subject := "Account verification"

	body, err := ParseTemplate("verification_link_email.html", struct {
		BaseUrl string
		Token   string
	}{Token: token, BaseUrl: s.Cfg.BaseUrl})
	if err != nil {
		return err
	}

	err = s.sendEmail(emailAddress, subject, body)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")
	return nil
}

func (s Service) SendPasswordResetLinkEmail(emailAddress string, token string) error {
	subject := "Password reset"

	body, err := ParseTemplate("password_reset_link_email.html", struct {
		BaseUrl string
		Token   string
	}{BaseUrl: s.Cfg.BaseUrl, Token: token})
	if err != nil {
		return err
	}

	err = s.sendEmail(emailAddress, subject, body)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")
	return nil
}

func (s Service) SendNewPasswordEmail(emailAddress string, newPassword string) error {
	subject := "Password Lock new password"

	body, err := ParseTemplate("new_password_email.html", struct {
		Password string
	}{Password: newPassword})
	if err != nil {
		return err
	}

	err = s.sendEmail(emailAddress, subject, body)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")
	return nil
}

func (s Service) sendEmail(to string, subject string, body string) error {

	receivers := []string{
		to,
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject = "Subject: " + subject + "\n"
	msg := []byte(subject + mime + "\n" + body)

	auth := smtp.PlainAuth("", s.Cfg.SmtpFrom, s.Cfg.FirebaseAppPassword, s.Cfg.SmtpHost)

	// Sending email.
	err := smtp.SendMail(s.Cfg.SmtpHost+":"+s.Cfg.SmtpPort, auth, s.Cfg.SmtpFrom, receivers, msg)
	if err != nil {
		return err
	}

	return nil
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {

	t1 := template.Must(template.New(templateFileName).ParseFS(templates.Files, templateFileName))

	var tpl bytes.Buffer
	if err := t1.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

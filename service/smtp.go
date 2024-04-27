package service

import (
	"fmt"
	"net/smtp"
)

func (s Service) SendAccountVerifiedEmail(emailAddress string, secretKey string) error {

	messageText := "\nThis is your secret key : " + secretKey + "\nWe strongly advise you to write it down in a physical form and to delete this email.\nWe also remind you that without this secret key, you will not be able to access rest of your passwords.\nThank you for using Password Locker."

	// Message.
	message := []byte("Subject: Password Locker Secret Key\r\n" + messageText)

	err := s.sendEmail(emailAddress, message)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")
	return nil
}

func (s Service) SendVerificationLinkEmail(emailAddress string, token string) error {
	messageText := "\nThis is your verification token : " + token
	// Message.
	message := []byte("Subject: Password Locker Verification token\r\n" + messageText)

	err := s.sendEmail(emailAddress, message)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")
	return nil
}

func (s Service) SendNewPasswordEmail(emailAddress string, token string) error {
	messageText := "\nThis is your new password token : " + token
	// Message.
	message := []byte("Subject: Password Locker new password token\r\n" + messageText)

	err := s.sendEmail(emailAddress, message)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")
	return nil
}

func (s Service) SendPasswordResetLinkEmail(emailAddress string, token string) error {
	return nil
}

func (s Service) sendEmail(to string, message []byte) error {

	receivers := []string{
		to,
	}

	auth := smtp.PlainAuth("", s.cfg.SmtpFrom, s.cfg.FirebaseAppPassword, s.cfg.SmtpHost)

	// Sending email.
	err := smtp.SendMail(s.cfg.SmtpHost+":"+s.cfg.SmtpPort, auth, s.cfg.SmtpFrom, receivers, message)
	if err != nil {
		return err
	}

	return nil
}

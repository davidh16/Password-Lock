package service

import (
	"fmt"
	"net/smtp"
)

func SendEmail(emailAddress string, secretKey string) error {

	// Sender data.
	from := "hey.clothing.shop@gmail.com"
	appsPassword := "ihgceqqtzsoajrsp"

	// Receiver email address.
	to := []string{
		emailAddress,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	messageText := "\nThis is your secret key : " + secretKey + "\nWe strongly advise you to write it down in a physical form and to delete this email.\nWe also remind you that without this secret key, you will not be able to access rest of your passwords.\nThank you for using Password Locker."

	// Message.
	message := []byte("Subject: Password Locker Secret Key\r\n" + messageText)
	// Authentication.
	auth := smtp.PlainAuth("", from, appsPassword, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")
	return nil
}

package utils

import (
	"log"
	"time"

	"github.com/pquerna/otp/totp"
	"gopkg.in/gomail.v2"
)

// GenerateOTP generates a one-time password
func GenerateOTP(accountName string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Base-beego",
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}

	otp, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", "", err
	}

	return otp, key.Secret(), nil
}

// SendOTP sends the OTP to the user's email
func SendOTP(email, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "MS_is8x0T@trial-351ndgw8pjxgzqx8.mlsender.net")
	m.SetHeader("To", email)
	m.SetAddressHeader("Cc", "MS_is8x0T@trial-351ndgw8pjxgzqx8.mlsender.net", "dat")
	m.SetHeader("Subject", "Confirm your email")
	m.SetBody("text/plain", "OTP: "+otp)

	d := gomail.NewDialer("smtp.mailersend.net", 587, "MS_is8x0T@trial-351ndgw8pjxgzqx8.mlsender.net", "fhAKrkUGToqW3LZ4")

	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email: ", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}

// VerifyOTP verifies the OTP entered by the user
func VerifyOTP(otp, secret string) bool {
	return totp.Validate(otp, secret)
}

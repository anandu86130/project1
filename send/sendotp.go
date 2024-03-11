package send

import (
	"fmt"
	"net/smtp"
)

func SendOTPByEmail(Email, Otp string) {
	auth := smtp.PlainAuth(
		"rcdyr",
		"sonusuni2255@gmail.com",
		"ggnhnsxxsnvvonmm",
		"smtp.gmail.com",
	)

	msg := []byte(Otp)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"sonusuni2255@gmail.com",
		[]string{Email},
		msg,
	)
	if err != nil {
		fmt.Println(err)
	}
}

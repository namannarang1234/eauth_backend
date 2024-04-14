package handlers

import (
	"fmt"
	"log"
	"os"

	"github.com/skip2/go-qrcode"
	"gopkg.in/gomail.v2"
)

const from = ""
const password = ""
const smtpHost = "smtp.gmail.com"
const smtpPort = "587"

func SendMail(to string, otp string, l *log.Logger) {
	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Login OTP")
	m.SetBody(
		"text/html",
		"<p>Login with the OTP or scan the QR code</p>"+
			"<p>Your OTP is <b>"+otp+"</b>, do not share it with anyone</p>",
	)
	filename := "./qrs/" + otp + ".png"
	qrcode.WriteFile(fmt.Sprintf("http://192.168.0.103:6969/verifyqr/%s/%s", to, otp), qrcode.Medium, 256, filename)
	m.Attach(filename)

	d := gomail.NewPlainDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		l.Println(err)
	}

	os.Remove(filename)
}

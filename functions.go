package main

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func WriteToFile(content string, path string) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SendMail(host string, port int, user string, password string, from string, to string, subject string, body string, attach string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	//m.SetAddressHeader("Cc", cc, "autobot")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if attach != "" {
		m.Attach(attach)
	}

	d := gomail.NewDialer(host, port, user, password)

	if err := d.DialAndSend(m); err != nil {
		println(err.Error())
		//panic(err)
	}
}

package main

import (
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

func WriteToFile(content string, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	l, err := f.WriteString(content)
	if err != nil {
		f.Close()
		log.Fatalln(err)
	}
	log.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		log.Fatalln(err)
	}
}

func SendMail(host string, port int, user string, password string, from string, to []string, subject string, body string, attach string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	//m.SetAddressHeader("Cc", cc, "autobot")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if attach != "" {
		m.Attach(attach)
	}

	d := gomail.NewDialer(host, port, user, password)

	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

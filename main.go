package main

import (
	"bytes"
	"log"
	"net/smtp"
	"strconv"
	"text/template"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Licenses []License

var cfg Config
var licenses Licenses

var auth smtp.Auth

func main() {

	readConfigFile(&cfg)

	f, err := excelize.OpenFile(cfg.Paths.ExcelFile)
	if err != nil {
		println(err.Error())
		return
	}

	// Get all the rows in the Sheet1.
	for i := 11; i < 101; i++ {
		client, err := f.GetCellValue("Лист1", "D"+strconv.Itoa(i))
		if err != nil {
			println(err.Error())
			return
		}
		if client != "" {
			dueDate, err := f.GetCellValue("Лист1", "F"+strconv.Itoa(i))
			if err != nil {
				println(err.Error())
				return
			}
			licenses = append(licenses, License{client, dueDate})
		}
	}

	var tpl bytes.Buffer
	t := template.Must(template.New("").Parse(`<body><table border="1"><td><strong>Bundle ID</strong></td><td><strong>Expiration date</strong></td>{{range .}}<tr><td>{{.Client}}</td><td>{{.DueDate}}</td></tr>{{end}}</table></body>`))
	if err := t.Execute(&tpl, licenses); err != nil {
		log.Fatal(err)
	}

	WriteToFile(tpl.String(), "./index.html")

	SendMail(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.From, cfg.SMTP.To, cfg.SMTP.Subject, tpl.String(), "")
}

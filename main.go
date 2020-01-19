package main

import (
	"bytes"
	"log"
	"net/smtp"
	"strconv"
	"text/template"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Licenses []License

var cfg Config
var licenses Licenses
var expiringLicenses Licenses

var auth smtp.Auth

func main() {

	readConfigFile(&cfg)

	f, err := excelize.OpenFile(cfg.Excel.File)
	if err != nil {
		println(err.Error())
		return
	}

	// Get all the rows in the Sheet1.
	for i := cfg.Excel.CheckingRowStart; i <= cfg.Excel.CheckingRowEnd; i++ {
		client, err := f.GetCellValue(cfg.Excel.Sheet, cfg.Excel.NameColumn+strconv.Itoa(i))
		if err != nil {
			println(err.Error())
			return
		}
		if client != "" {
			dueDateStr, err := f.GetCellValue(cfg.Excel.Sheet, cfg.Excel.CheckingColumn+strconv.Itoa(i))
			if err != nil {
				println(err.Error())
				return
			}
			licenses = append(licenses, License{client, dueDateStr})
		}
	}

	//sending all licenses to admins
	var tpl bytes.Buffer
	t := template.Must(template.New("").Parse(`<body><h1>All data</h1><br><table border="1"><td><strong>` + cfg.Common.TableHeaderNameColumn + `</strong></td><td><strong>` + cfg.Common.TableHeaderCheckingColumn + `</strong></td>{{range .}}<tr><td>{{.Client}}</td><td>{{.DueDate}}</td></tr>{{end}}</table></body>`))
	if err := t.Execute(&tpl, licenses); err != nil {
		log.Fatal(err)
	}

	//WriteToFile(tpl.String(), "./index.html")
	SendMail(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.From, cfg.Common.AdminsEmails, cfg.SMTP.Subject, tpl.String(), "")

	//sending expiring licenses
	today := time.Now()
	for _, v := range licenses {
		err = nil
		dueDate, err := time.Parse("01-02-06", v.DueDate)
		if err == nil {
			if dueDate.Before(today.Add(cfg.Common.NotifyForDays * 24 * time.Hour)) {
				println(v.Client, dueDate.String())
				expiringLicenses = append(expiringLicenses, License{v.Client, v.DueDate})
			}
		}
	}

	tpl.Reset()
	t = template.Must(template.New("").Parse(`<body><h1>` + cfg.Common.TableCaption + `</h1><br><table border="1"><td><strong>` + cfg.Common.TableHeaderNameColumn + `</strong></td><td><strong>` + cfg.Common.TableHeaderCheckingColumn + `</strong></td>{{range .}}<tr><td>{{.Client}}</td><td>{{.DueDate}}</td></tr>{{end}}</table></body>`))
	if err := t.Execute(&tpl, expiringLicenses); err != nil {
		log.Fatal(err)
	}

	SendMail(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.From, cfg.SMTP.To, cfg.SMTP.Subject, tpl.String(), "")
}

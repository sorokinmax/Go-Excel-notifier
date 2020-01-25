package main

import (
	"bytes"
	"io"
	"log"
	"net/smtp"
	"os"
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
	log.SetFlags(log.LstdFlags)
	lf, err := os.OpenFile("output.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer lf.Close()
	multi := io.MultiWriter(os.Stdout, lf)
	log.SetOutput(multi)

	readConfigFile(&cfg)

	f, err := excelize.OpenFile(cfg.Excel.File)
	if err != nil {
		log.Fatalln(err)
	}

	today := time.Now()

	// Get all the rows in the Sheet.
	for i := cfg.Excel.CheckingRowStart; i <= cfg.Excel.CheckingRowEnd; i++ {
		client, err := f.GetCellValue(cfg.Excel.Sheet, cfg.Excel.NameColumn+strconv.Itoa(i))
		if err != nil {
			log.Fatalln(err)
		}
		if client != "" {
			dueDateStr, err := f.GetCellValue(cfg.Excel.Sheet, cfg.Excel.CheckingColumn+strconv.Itoa(i))
			if err != nil {
				log.Fatalln(err)
			}
			licenses = append(licenses, License{client, dueDateStr})
		}
	}

	//sending all licenses to admins
	var tpl bytes.Buffer
	t := template.Must(template.New("").Parse(`<body><h1>All data</h1><br><table border="1"><td><strong>` + cfg.Common.TableHeaderNameColumn + `</strong></td><td><strong>` + cfg.Common.TableHeaderCheckingColumn + `</strong></td>{{range .}}<tr><td>{{.BundleID}}</td><td>{{.DueDate}}</td></tr>{{end}}</table></body>`))
	if err := t.Execute(&tpl, licenses); err != nil {
		log.Fatalln(err)
	}

	//WriteToFile(tpl.String(), "./index.html")
	SendMail(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.From, cfg.Common.AdminsEmails, cfg.SMTP.Subject, tpl.String(), "")

	//sending all expiring licenses
	for _, v := range licenses {
		err = nil
		dueDate, err := time.Parse("01-02-06", v.DueDate)
		if err != nil {
			log.Println(err)
		}
		if dueDate.Before(today.Add(cfg.Common.NotifyForDays * 24 * time.Hour)) {
			//log.Println(v.BundleID, dueDate.Format("2006-01-02"))
			expiringLicenses = append(expiringLicenses, License{v.BundleID, dueDate.Format("2006-01-02")})
		}
	}

	tpl.Reset()
	t = template.Must(template.New("").Parse(`<body><h1>` + cfg.Common.TableCaption + `</h1><br><table border="1"><td><strong>` + cfg.Common.TableHeaderNameColumn + `</strong></td><td><strong>` + cfg.Common.TableHeaderCheckingColumn + `</strong></td>{{range .}}<tr><td>{{.BundleID}}</td><td>{{.DueDate}}</td></tr>{{end}}</table></body>`))
	if err := t.Execute(&tpl, expiringLicenses); err != nil {
		log.Fatal(err)
	}

	SendMail(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.From, cfg.SMTP.To, cfg.SMTP.Subject, tpl.String(), "")

	//sending personal expiring license
	expiringLicenses = nil
	for _, v := range licenses {
		err = nil
		dueDate, err := time.Parse("01-02-06", v.DueDate)
		if err != nil {
			log.Println(err)
		}
		if dueDate.Before(today.Add(cfg.Common.NotifyForDays*24*time.Hour)) && contains(cfg.Excel.PersonalNotificationBundles, v.BundleID) {
			//log.Println(v.BundleID, dueDate.Format("2006-01-02"))
			expiringLicenses = append(expiringLicenses, License{v.BundleID, dueDate.Format("2006-01-02")})
		}
	}

	tpl.Reset()
	t = template.Must(template.New("").Parse(`<body><h1>` + cfg.Common.TableCaption + `</h1><br><table border="1"><td><strong>` + cfg.Common.TableHeaderNameColumn + `</strong></td><td><strong>` + cfg.Common.TableHeaderCheckingColumn + `</strong></td>{{range .}}<tr><td>{{.BundleID}}</td><td>{{.DueDate}}</td></tr>{{end}}</table></body>`))
	if err := t.Execute(&tpl, expiringLicenses); err != nil {
		log.Fatal(err)
	}

	SendMail(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.From, cfg.Excel.PersonalNotificationEmails, cfg.SMTP.Subject, tpl.String(), "")
}

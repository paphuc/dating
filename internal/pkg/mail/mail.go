package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"dating/internal/app/config"
)

type Mail struct {
	Subject string
	Body    string
}

func Sendmail(content Mail, mails []string, conf *config.Configs) error {

	var (
		HostMail = conf.Mail.Smtp.HostMail
		PortMail = conf.Mail.Smtp.PortMail
		Address  = HostMail + ":" + PortMail

		email    = conf.Mail.Email
		password = conf.Mail.Password
	)

	auth := smtp.PlainAuth("", email, password, HostMail)

	t, err := template.ParseFiles(conf.Mail.SrcTemplate)
	fmt.Println(t, err)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", content.Subject, mimeHeaders)))

	t.Execute(&body, struct {
		Code string
	}{
		Code: content.Body,
	})

	error := smtp.SendMail(Address, auth, email, mails, body.Bytes())
	if error != nil {
		return error
	}

	return nil
}

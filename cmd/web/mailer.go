package main

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct { // Mail server
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string // Default from address
	FromName    string // Default from name
	Wait        *sync.WaitGroup
	MailerChan  chan Message
	ErrorChan   chan error
	DoneChan    chan bool
}

type Message struct { // Email message
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string // File paths
	Data        any      // the body of the message
	DataMap     map[string]any
	Template    string
}

func (app *Config) listenForMail() {
	for {
		select {
		case msg := <-app.Mailer.MailerChan:
			go app.Mailer.Send(msg, app.Mailer.ErrorChan)
		case err := <-app.Mailer.ErrorChan:
			app.ErrorLog.Println("Error sending email: ", err)
		case <-app.Mailer.DoneChan:
			app.InfoLog.Println("Stopping email service...")
			return // exit goroutine
		}
	}
}

// Helpful wrapper function to send email
func (app *Config) sendEmail(msg Message) {
	app.Wait.Add(1)
	app.Mailer.MailerChan <- msg
}

func (m *Mail) Send(msg Message, errorChan chan error) {
	defer m.Wait.Done()

	// set defaults
	if msg.Template == "" {
		msg.Template = "mail" // capture default template name
	}
	if msg.From == "" {
		msg.From = m.FromAddress // capture default from address
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName // capture default from name
	}

	// populate data
	data := map[string]any{
		"message": msg.Data,
	}
	msg.DataMap = data

	// build html mail
	formattedMessage, err := m.buildHTML(msg)
	if err != nil {
		errorChan <- err
	}

	// build plaintext mail
	plainMessage, err := m.buildPlainText(msg)
	if err != nil {
		errorChan <- err
	}

	// setup SMTP client
	mailServer := mail.NewSMTPClient()
	mailServer.Host = m.Host
	mailServer.Port = m.Port
	mailServer.Username = m.Username
	mailServer.Password = m.Password
	mailServer.Encryption = m.getEncryption(m.Encryption)
	mailServer.KeepAlive = false // not expecting to send mail constantly
	mailServer.ConnectTimeout = 10 * time.Second
	mailServer.SendTimeout = 10 * time.Second

	smptClient, err := mailServer.Connect()
	if err != nil {
		errorChan <- err
	}

	// create email
	email := mail.NewMSG()
	email.SetFrom(msg.From)
	email.AddTo(msg.To)
	email.SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)
	if len(msg.Attachments) > 0 {
		for _, file := range msg.Attachments {
			email.AddAttachment(file)
		}
	}

	// send email
	err = email.Send(smptClient)
	if err != nil {
		errorChan <- err
	}
}

func (m *Mail) buildHTML(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s.html.gohtml", pathToTemplates, msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil { // execute the template, populating the body tag
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     true,  // remove classes
		CssToAttributes:   false, // convert css to attributes
		KeepBangImportant: true,  // keep !important
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) buildPlainText(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("%s/%s.plain.gohtml", pathToTemplates, msg.Template)

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil { // execute the template, populating the body tag
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) getEncryption(e string) mail.Encryption {
	switch e {
	case "SSL":
		return mail.EncryptionSSLTLS
	case "TLS":
		return mail.EncryptionSTARTTLS
	case "none":
		return mail.EncryptionNone // development only
	default:
		return mail.EncryptionSTARTTLS
	}
}

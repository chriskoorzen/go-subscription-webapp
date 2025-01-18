package main

import (
	"sync"
	"time"

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

// TODO function to listen for messages to send on MailerChan

func (m *Mail) Send(msg Message, errorChan chan error) {
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
	return "", nil
}

func (m *Mail) buildPlainText(msg Message) (string, error) {
	return "", nil
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

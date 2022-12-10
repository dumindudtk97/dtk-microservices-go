package main

import (
	"bytes"
	"html/template"

	"log"
	"time"

	//Package premailer is for inline styling (use css and auto convert)
	"github.com/vanng822/go-premailer/premailer"
	// Go Simple Mail is a simple and efficient package to send emails.
	mail "github.com/xhit/go-simple-mail/v2"
)

// setup the instance of mail with appropriate configurations
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

// contents of email
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

// function to send a email message
func (m *Mail) SendSMTPMessage(msg Message) error {

	// make sure to fill from address and name
	if msg.From == "" {
		msg.From = m.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// data map to pass to template
	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	// build html version of the message (with inclineCSS)
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	// build plaintext version of the message
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	//building the mail server
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// start smtp client server
	smtpClient, err := server.Connect()
	if err != nil {
		log.Println("Can't do server.Connect() with mail.NewSMTPClient(), smtpClient not created")
		return err
	}

	// create email from recieved parameter msg
	email := mail.NewMSG() //NewMSG creates a new email.
	// set parameters
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)
	// adding body to email
	email.SetBody(mail.TextPlain, plainMessage) //SetBody sets the body of the email message.
	email.AddAlternative(mail.TextHTML, formattedMessage)

	// add attachments to email
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	//send email //email and client has all info needed
	err = email.Send(smtpClient)
	if err != nil {
		log.Println("Can't do email.send with smtpClient")
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	// html template
	templateToRender := "./templates/mail.html.gohtml"
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		log.Println("can't open html template file", err)
		return "", err
	}

	var tpl bytes.Buffer // to read template into
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	// inline CSS is used to apply a unique style to a single HTML element
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {

	// . is current directory then ./templates
	templateToRender := "./templates/mail.plain.gohtml"
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		log.Println("can't open plain template file", err)
		return "", err
	}

	var tpl bytes.Buffer
	// execute template and add datamap to tpl
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String() // tpl to string (plaintext)

	return plainMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	//create a premailer instance
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	//transform and inlining css
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

package main

//	go get github.com/vanng822/go-premailer/premailer
//	go get github.com/xhit/go-simple-mail/v2

type Mail struct {
	Domain string
	Host string
	Port int
	Username string
	Password string
	Encrption string
	FromAddress string
	FromName string
}

type Message struct {
	From string
	FromName string
	To string
	Subject string
	Attachment []string
	Data any
	DataMap map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error{
	
	//make sure from are set 
	if msg.From == "" {
		msg.From = m.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		log.Panic(err)
	}

	plainMsg, err := m.buidPlainTextMessage(msg)
	if err != nil {
		log.Panic(err)
	}

	server = mail.new(smtpClient)

	server.host = m.host 
	server.port = m.port
	server.username = m.username
	server.password = m.password
	server.encryption = m.getEncryption(m.Encrption)
	server.keepAlive = false
	server.timeout = time.Second * 10
	server.sendtimeout = time.Second * 10

	client, err = smtpClient()
	if err != nil {
		log.Panic(err)
	}
	
	
	

}

func (m *Mail) buildHTMLMessage(msg) (string, error) {
	templateToRender := ""./template/mail.gohtml

	t := template.FromHTMLTemplate


	if err != nil {
		log.Panic(err)
	}


	if err != nil {
		log.Panic(err)
	}

	formattedMessage := 

	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		log.Panic(err)
	}

	return formattedMessage, nil

}


func (m *Mail) inlineCSS(f string) (string, error) {
	var options = premailer.Options{
		removeClass: false,
	}

	prim = premailer.new(premailer)

	if err != nil {
		log.Panic(err)
	}

	html := prim.response

	return html, nil

}

func (m *Mail) buidPlainTextMessage(f string) (string, error){
	
	plainMsg 

	return plainMsg, nil
}

func (m *Mail) getEncryption(f s) (string, error){
	switch s {
	case SSL:
		mail.encryptssl = true
	case TLS:
		mail.encrypttls = true
	case None:
		mail.encrypt = false
	}
}

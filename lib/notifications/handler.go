package notifications

import (
	"crypto/tls"
	"log"
	"strings"

	"github.com/go-gomail/gomail"
)

type Handler interface {
	Post(text string) (int, error)
}

func NewLogHandler(
	logger *log.Logger,
) Handler {

	return &LogHandler{
		logger: logger,
	}
}

type LogHandler struct {
	logger *log.Logger
}

func (c *LogHandler) Post(text string) (int, error) {
	c.logger.Println(text)

	return 0, nil
}

type MailHandler struct {
	sender   string
	receiver string
	host     string
	port     int
	username string
	password string
}

func NewMailHandler(
	logger *log.Logger,
	sender string,
	receiver string,
	host string,
	port int,
	username string,
	password string,
) Handler {

	return &MailHandler{
		sender:   sender,
		receiver: receiver,
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (c *MailHandler) Post(text string) (int, error) {
	endIdx := strings.Index(text, "] ")
	subject := "QNAP Notification"
	if endIdx != -1 {
		subject = subject + " " + text[:endIdx+1]
		text = text[endIdx+2:]
	}

	m := gomail.NewMessage()
	m.SetHeader("From", c.sender)
	m.SetHeader("To", c.receiver)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(c.host, c.port, c.username, c.password)
	if c.port == 25 {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if err := d.DialAndSend(m); err != nil {
		log.Println(err.Error())
		return -1, err
	}

	return 0, nil
}

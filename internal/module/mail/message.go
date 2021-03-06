package mail

import (
	"github.com/teablog/tea/internal/config"
	"github.com/teablog/tea/internal/logger"
	"gopkg.in/gomail.v2"
	"time"
)

var ch chan *gomail.Message

func Init() {
	ch = make(chan *gomail.Message, 5)
	go start()
}

func Send(msg *Message) {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress("douyacun@163.com", "Douyacun"))
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Title)
	m.SetBody(msg.ContentType, msg.Body)
	ch <- m
	return
}

func start() {
	defer close(ch)
	d := gomail.NewDialer(config.Email.Host(), config.Email.Port(), config.Email.Username(), config.Email.Password())

	var s gomail.SendCloser
	var err error
	open := false
	for {
		select {
		case m, ok := <-ch:
			if !ok {
				return
			}
			// 禁止发送邮件
			if !config.Email.Enable() {
				return
			}
			if !open {
				if s, err = d.Dial(); err != nil {
					logger.Errorf("gomail dial %s err", config.Email.Host())
				}
				open = true
			}
			if err := gomail.Send(s, m); err != nil {
				logger.Wrapf(err, "gmail send err ")
			}
		// Close the connection to the SMTP server if no email was sent in
		// the last 30 seconds.
		case <-time.After(30 * time.Second):
			if open {
				_ = s.Close()
				open = false
			}
		}
	}
}

package mail

import (
	"github.com/teablog/tea/internal/config"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	config.Init("configs/debug.ini")
	Init()
	m := NewMessage().SetTo("douyacun@163.com").SetTitle("报警邮件").SetBody("出bug，快来修复")
	Send(m)
	time.Sleep(5 * time.Second)
}
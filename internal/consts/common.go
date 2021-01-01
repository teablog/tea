package consts

const (
	PageSize    = 10
	DefaultPage = 1
	Host        = "http://www.douyacun.com"
	HostDev     = "http://localhost:3000"
)

type Status int

const (
	StatusNil Status = iota
	StatusOn
	StatusOff
	StatusDel
)

const (
	TimeYMD = "2006-01-02"
	TimeYMDHIS = "2006-01-02 15:04:05"
)
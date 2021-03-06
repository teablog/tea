package consts

const (
	PageSize    = 10
	DefaultPage = 1
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
	TimeYM     = "2006-01-02"
	TimeYMD    = "2006-01-02"
	TimeYMDHIS = "2006-01-02 15:04:05"
)

const CookieName = "douyacun"
const CookieUUID = "douyacun-uuid"
const CookieUUIDV2 = "douyacun-uuid-v2"

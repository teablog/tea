package qqwry

import (
	"bufio"
	"fmt"
	iconv "github.com/qiniu/iconv"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type qqwryData struct {
	startIp int64
	endIp   int64
	country string
	local   string

	countryFlag int
	fp          *os.File

	firstStartIp int64
	lastStartIp  int64
	endIpOff     int64
}

// Getqqdata 打开qqwry.dat 并返回*os.File
func Getqqdata(qqwryfile string) (file *os.File, err error) {
	file, err = os.Open(qqwryfile) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	return
}

// Getlocation 根据ip（xxx.xxx.xxx.xxx）返回地区
func Getlocation(file *os.File, ip string) (country string, location string) {
	//初始化
	qqdata := qqwryData{}
	qqdata.startIp = 0
	qqdata.endIp = 0
	qqdata.country = ""
	qqdata.local = ""
	qqdata.countryFlag = 0
	qqdata.firstStartIp = 0
	qqdata.lastStartIp = 0
	qqdata.endIpOff = 0
	qqdata.fp = file
	qqdata.qqwry(ip)
	country = iconvStr(qqdata.country)
	location = iconvStr(qqdata.local)
	return
}

func (qqdata *qqwryData) qqwry(dotip string) (result string) {
	ipint := ip_to_int(dotip)

	qqdata.fp.Seek(0, os.SEEK_SET)
	buf := make([]byte, 8)
	_, err := qqdata.fp.Read(buf)
	if err != nil {
		log.Fatal(err)
		fmt.Println("qqwry:", err)
	}
	b := buf[0:4]
	a := buf[4:8]
	qqdata.firstStartIp = int64(b[0]) + (int64(b[1]) * 256) + (int64(b[2]) * 256 * 256) + (int64(b[3]) * 256 * 256 * 256)
	qqdata.lastStartIp = int64(a[0]) + (int64(a[1]) * 256) + (int64(a[2]) * 256 * 256) + (int64(a[3]) * 256 * 256 * 256)

	recordCount := (qqdata.lastStartIp - qqdata.firstStartIp) / 7

	if recordCount <= 1 {
		qqdata.country = "FileDataError"
		result = "2"
		return
	}

	// Match ...
	var rangB int64 = 0
	var rangE int64 = recordCount
	var recNo int64
	for rangB < rangE-1 {
		recNo = (rangB + rangE) / 2
		qqdata.getStartIp(recNo)

		if ipint == qqdata.startIp {
			rangB = recNo
			break
		}

		if ipint > qqdata.startIp {
			rangB = recNo
		} else {
			rangE = recNo
		}
	}

	qqdata.getStartIp(rangB)
	qqdata.getEndIp()

	if (qqdata.startIp <= ipint) && (qqdata.endIp >= ipint) {
		qqdata.getCountry()
	} else {
		qqdata.country = "未知"
		qqdata.local = ""
	}

	return
}

func (qqdata *qqwryData) getStartIp(recNo int64) (startIp int64) {
	offset := qqdata.firstStartIp + recNo*7
	qqdata.fp.Seek(offset, os.SEEK_SET)
	buf := make([]byte, 7)
	_, err := qqdata.fp.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	b := buf[4:7]
	a := buf[0:4]

	qqdata.endIpOff = int64(b[0]) + (int64(b[1]) * 256) + (int64(b[2]) * 256 * 256)
	qqdata.startIp = int64(a[0]) + (int64(a[1]) * 256) + (int64(a[2]) * 256 * 256) + (int64(a[3]) * 256 * 256 * 256)
	startIp = qqdata.startIp
	return
}

func (qqdata *qqwryData) getEndIp() (endIp int64) {
	qqdata.fp.Seek(qqdata.endIpOff, os.SEEK_SET)
	buf := make([]byte, 5)
	_, err := qqdata.fp.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	b := buf[0:4]
	a := buf[4:5]

	qqdata.endIp = int64(b[0]) + (int64(b[1]) * 256) + (int64(b[2]) * 256 * 256) + (int64(b[3]) * 256 * 256 * 256)
	qqdata.countryFlag = int(a[0])
	endIp = qqdata.endIp
	return
}

func (qqdata *qqwryData) getCountry() {
	switch qqdata.countryFlag {
	case 1, 2:
		qqdata.country = qqdata.getFlagStr(qqdata.endIpOff + 4)
		if 1 == qqdata.countryFlag {
			qqdata.local = ""
		} else {
			qqdata.local = qqdata.getFlagStr(qqdata.endIpOff + 8)
		}
	default:
		qqdata.country = qqdata.getFlagStr(qqdata.endIpOff + 4)
		//qqdata.local = qqdata.getFlagStr(ftell(qqdata.fp))
	}
}

func (qqdata *qqwryData) getFlagStr(offset int64) (local string) {
	flag := 0
	for true {
		qqdata.fp.Seek(offset, os.SEEK_SET)
		buf := make([]byte, 1)
		_, err := qqdata.fp.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		flag = int(buf[0])
		if flag == 1 || flag == 2 {
			buf = make([]byte, 3)
			_, err := qqdata.fp.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			if flag == 2 {
				qqdata.countryFlag = 2
				qqdata.endIpOff = offset - 4
			}

			offset = int64(buf[0]) + (int64(buf[1]) * 256) + (int64(buf[2]) * 256 * 256)
		} else {
			break
		}
	}

	if offset < 12 {
		return
	}

	qqdata.fp.Seek(offset, os.SEEK_SET)
	local = qqdata.getStr()
	return
}

func (qqdata *qqwryData) getStr() (str string) {
	str = ""

	br := bufio.NewReader(qqdata.fp)
	str, err := br.ReadString(0x00)
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "000A")
	if err != nil {
		panic(err)
	}

	return
}

func ip_to_int(ip string) (r int64) {
	b := net.ParseIP(ip)
	r = inet_aton(b)
	return
}

/*
 * Convert uint to net.IP
 */
func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

/*
 * Convert net.IP to int64
 */
func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func iconvStr(str string) string {
	cd, err := iconv.Open("utf-8", "gb2312")
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return ""
	}
	defer cd.Close()
	gbk := cd.ConvString(str)
	return gbk
}

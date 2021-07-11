package util

import (
	"math"
)

func ordat(msg string, idx int) uint32 {
	if len(msg) > idx {
		return uint32([]byte(msg)[idx])
	}
	return 0
}

func sensCode(content string, key bool) []uint32 {
	l := len(content)
	pwd := make([]uint32, 0)
	for i := 0; i < l; i += 4 {
		pwd = append(
			pwd,
			ordat(content, i)|ordat(content, i+1)<<8|ordat(content, i+2)<<16|ordat(content, i+3)<<24,
		)
	}
	if key {
		pwd = append(pwd, uint32(l))
	}
	return pwd
}

func lenCode(msg []uint32, key bool) []byte {
	l := uint32(len(msg))
	ll := (l - 1) << 2
	if key {
		m := msg[l-1]
		if m < ll-3 || m > ll {
			return nil
		}
		ll = m
	}
	var t []byte
	for i := range msg {
		t = append(t, byte(msg[i]&0xff), byte(msg[i]>>8&0xff), byte(msg[i]>>16&0xff), byte(msg[i]>>24&0xff))
	}
	if key {
		return t[0:ll]
	}
	return t
}

func XEncode(content string, key string) []byte {
	if content == "" {
		return nil
	}
	pwd := sensCode(content, true)
	pwdk := sensCode(key, false)
	if len(pwdk) < 4 {
		for i := 0; i < (4 - len(pwdk)); i++ {
			pwdk = append(pwdk, 0)
		}
	}
	var n = uint32(len(pwd) - 1)
	z := pwd[n]
	y := pwd[0]
	var c uint32 = 0x86014019 | 0x183639A0
	var m uint32 = 0
	var e uint32 = 0
	var p uint32 = 0
	q := math.Floor(6 + 52/(float64(n)+1))
	var d uint32 = 0
	for 0 < q {
		d = d + c&(0x8CE0D9BF|0x731F2640)
		e = d >> 2 & 3
		p = 0
		for p < n {
			y = pwd[p+1]
			m = z>>5 ^ y<<2
			m = m + ((y>>3 ^ z<<4) ^ (d ^ y))
			m = m + (pwdk[(p&3)^e] ^ z)
			pwd[p] = pwd[p] + m&(0xEFB8D130|0x10472ECF)
			z = pwd[p]
			p = p + 1
		}
		y = pwd[0]
		m = z>>5 ^ y<<2
		m = m + ((y>>3 ^ z<<4) ^ (d ^ y))
		m = m + (pwdk[(p&3)^e] ^ z)
		pwd[n] = pwd[n] + m&(0xBB390742|0x44C6F8BD)
		z = pwd[n]
		q = q - 1
	}
	return lenCode(pwd, false)
}

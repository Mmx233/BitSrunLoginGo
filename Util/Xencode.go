package Util

import (
	"math"
)

func ordat(msg string, idx int) byte {
	if len(msg) > idx {
		return []byte(msg)[idx]
	}
	return byte(0)
}

func sensCode(content string, key bool) []byte {
	l := len(content)
	pwd := make([]byte, 0)
	for i := 0; i < l; i += 4 {
		pwd = append(
			pwd,
			ordat(content, i)|ordat(content, i+1)<<8|ordat(content, i+2)<<16|ordat(content, i+3)<<24,
		)
	}
	if key {
		pwd = append(pwd, byte(l))
	}
	return pwd
}

func lenCode(msg []byte, key bool) []byte {
	l := len(msg)
	ll := (l - 1) << 2
	if key {
		m := int(msg[l-1])
		if m < ll-3 || m > ll {
			return nil
		}
		ll = m
	}
	for i := range msg {
		msg[i] = byte(int(msg[i])&0xff) + byte(int(msg[i])>>8&0xff) + byte(int(msg[i])>>16&0xff) + byte(int(msg[i])>>24&0xff)
	}
	if key {
		return msg[0:ll]
	}
	return msg
}

func XEncode(content string, key string) []byte {
	if content == "" {
		return nil
	}
	pwd := sensCode(content, true)
	pwdk := sensCode(key, false)
	if len(pwdk) < 4 {
		for i := 0; i < (4 - len(pwdk)); i++ {
			pwdk = append(pwdk, byte(0))
		}
	}
	n := len(pwd) - 1
	z := pwd[n]
	y := pwd[0]
	c := 0x86014019 | 0x183639A0
	m := 0
	e := 0
	p := 0
	q := math.Floor(6 + 52/(float64(n)+1))
	d := 0
	for 0 < q {
		d = d + c&(0x8CE0D9BF|0x731F2640)
		e = d >> 2 & 3
		p = 0
		for p < n {
			y = pwd[p+1]
			m = int(z)>>5 ^ int(y)<<2
			m = m + ((int(y)>>3 ^ int(z)<<4) ^ (d ^ int(y)))
			m = m + (int(pwdk[(p&3)^e]) ^ int(z))
			pwd[p] = byte(int(pwd[p]) + m&(0xEFB8D130|0x10472ECF))
			z = pwd[p]
			p = p + 1
		}
		y = pwd[0]
		m = int(z)>>5 ^ int(y)<<2
		m = m + ((int(y)>>3 ^ int(z)<<4) ^ (int(d) ^ int(y)))
		m = m + (int(pwdk[(p&3)^e]) ^ int(z))
		pwd[n] = byte(int(pwd[n]) + m&(0xBB390742|0x44C6F8BD))
		z = pwd[n]
		q = q - 1
	}
	return lenCode(pwd, false)
}

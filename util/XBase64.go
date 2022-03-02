package util

func getbyte(a byte) int {
	x := int(a)
	if x > 255 {
		Log.Fatal("INVALID_CHARACTER_ERR: DOM Exception 5")
	}
	return x
}

func Base64(s []byte) string {
	const ALPHA = "LVoJPiCN2R8G90yg+hmFHuacZ1OWMnrsSTXkYpUq/3dlbfKwv6xztjI7DeBE45QA"
	const PADCHAR = "="
	i := 0
	b10 := 0
	var x []byte
	imax := len(s) - len(s)%3
	if len(s) == 0 {
		return ""
	}
	for i := 0; i < imax; i += 3 {
		b10 = (getbyte(s[i]) << 16) | (getbyte(s[i+1]) << 8) | (getbyte(s[i+2]))
		x = append(x, ALPHA[(b10>>18)])
		x = append(x, ALPHA[((b10>>12)&63)])
		x = append(x, ALPHA[((b10>>6)&63)])
		x = append(x, ALPHA[(b10&63)])
	}
	i = imax
	if len(s)-imax == 1 {
		b10 = getbyte(s[i]) << 16
		x = append(x, ALPHA[(b10>>18)], ALPHA[((b10>>12)&63)], PADCHAR[0], PADCHAR[0])
	} else if len(s)-imax == 2 {
		b10 = (getbyte(s[i]) << 16) | (getbyte(s[i+1]) << 8)
		x = append(x, ALPHA[(b10>>18)], ALPHA[((b10>>12)&63)], ALPHA[((b10>>6)&63)], PADCHAR[0])
	}
	return string(x)
}

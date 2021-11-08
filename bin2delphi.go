package bin2delphi

import (
	"bytes"
	"fmt"
)

func NewUnit(name string) *Unit {
	u := &Unit{isUTF8: isUTF8(name)}
	fmt.Fprintf(u, "unit %s;\r\n\r\ninterface", name)
	return u
}

type Unit struct {
	bytes.Buffer
	isUTF8 bool
}

func isUTF8(s string) bool {
	for _, b := range s {
		if b >= 128 {
			return true
		}
	}
	return false
}

func (u *Unit) AddConstant(name string, data []byte) {
	u.isUTF8 = u.isUTF8 || isUTF8(name)

	if len(data) == 0 {
		fmt.Fprintf(u, "\r\n\r\nconst\r\n  %s: Pointer = nil;", name)
	} else {
		fmt.Fprintf(
			u,
			"\r\n\r\nconst\r\n  %s: array [0 .. %d] of Byte = (",
			name,
			len(data)-1,
		)

		const maxBytesInLine = 15
		bytesInLine := maxBytesInLine
		for i, b := range data {
			if bytesInLine >= maxBytesInLine {
				fmt.Fprint(u, "\r\n    ")
				bytesInLine = 0
			} else {
				fmt.Fprint(u, " ")
			}
			fmt.Fprintf(u, "$%02X", b)
			if i != len(data)-1 {
				fmt.Fprint(u, ",")
			}
			bytesInLine++
		}

		fmt.Fprintf(u, "\r\n  );")
	}
}

func (u *Unit) Generate() []byte {
	u.WriteString("\r\n\r\nimplementation\r\n\r\nend.\r\n")
	if u.isUTF8 {
		// Prepend the UTF-8 byte order mark.
		return append([]byte{0xEF, 0xBB, 0xBF}, u.Bytes()...)
	}
	return u.Bytes()
}

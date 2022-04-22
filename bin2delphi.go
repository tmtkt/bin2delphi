package bin2delphi

import (
	"bytes"
	"fmt"
)

// NewUnit creates a new empty unit of the given name. You can call AddConstant
// on the Unit to fill it with constants. Call Generate when you are done to get
// the code.
func NewUnit(name string) *Unit {
	u := &Unit{isUTF8: isUTF8(name)}
	fmt.Fprintf(u, "unit %s;\r\n\r\ninterface", name)
	return u
}

// Unit represents a Delphi file (.pas). Create one with NewUnit, add constants
// with AddConstant and then generate the code with Generate.
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

// AddConstant adds a byte array constant of the given name, with the given data
// to the unit.
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

// Generate returns the Delphi code for the unit. It contains the unit name, all
// the constants and an empty implementation section.
//
// You can only call Generate once. If you call it multiple times, the generated
// code will be broken.
func (u *Unit) Generate() []byte {
	u.WriteString("\r\n\r\nimplementation\r\n\r\nend.\r\n")
	if u.isUTF8 {
		// Prepend the UTF-8 byte order mark.
		return append([]byte{0xEF, 0xBB, 0xBF}, u.Bytes()...)
	}
	return u.Bytes()
}

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	unitName  = flag.String("unit", "Main", "Unit name. Empty string to omit unit boilerplate.")
	constName = flag.String("const", "", "Constant name to use. Must not be empty.")
)

func main() {
	flag.Parse()

	if *constName == "" {
		flag.Usage()
		return
	}

	wantUnitBoilerplate := *unitName != ""

	if wantUnitBoilerplate {
		if !ascii(*constName) || !ascii(*unitName) {
			// Encode this as UTF-8, start with the BOM.
			fmt.Print(0xEF, 0xBB, 0xBF)
		}
		fmt.Printf("unit %s;\r\n\r\ninterface\r\n\r\n", *unitName)
	}

	gen := generator{bytesInLine: maxBytesInLine}
	_, err := io.Copy(&gen, os.Stdin)
	if err != nil {
		panic(err)
	}
	if gen.byteCount == 0 {
		panic("Delphi constant arrays cannot be empty (size 0)")
	}

	fmt.Printf(
		"const\r\n  %s: array [0 .. %d] of Byte = (%s\r\n);",
		*constName,
		gen.byteCount-1,
		gen.buf.Bytes(),
	)

	if wantUnitBoilerplate {
		fmt.Print("\r\n\r\nimplementation\r\n\r\nend.\r\n")
	}
}

func ascii(s string) bool {
	for _, b := range s {
		if b >= 128 {
			return false
		}
	}
	return true
}

type generator struct {
	buf         bytes.Buffer
	bytesInLine int
	byteCount   int
}

const maxBytesInLine = 15

func (g *generator) Write(p []byte) (int, error) {
	if g.byteCount > 0 {
		fmt.Fprint(&g.buf, ",")
	}
	for i, b := range p {
		if g.bytesInLine >= maxBytesInLine {
			fmt.Fprint(&g.buf, "\r\n    ")
			g.bytesInLine = 0
		} else {
			fmt.Fprint(&g.buf, " ")
		}
		fmt.Fprintf(&g.buf, "$%02X", b)
		if i != len(p)-1 {
			fmt.Fprint(&g.buf, ",")
		}
		g.bytesInLine++
	}
	g.byteCount += len(p)
	return len(p), nil
}

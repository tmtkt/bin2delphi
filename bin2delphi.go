package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	packageName = flag.String("unit", "Main", "Unit name. Empty string to omit unit boilerplate.")
	varName     = flag.String("const", "", "Constant name to use. Must not be empty.")
)

func main() {
	flag.Parse()

	if *varName == "" {
		flag.Usage()
		return
	}

	if *packageName != "" {
		fmt.Print(`unit Main;

interface

`)
	}

	gen := generator{bytesInLine: maxBytesInLine}
	_, err := io.Copy(&gen, os.Stdin)
	if err != nil {
		panic(err)
	}
	if gen.byteCount == 0 {
		panic("Delphi constant arrays cannot be empty (size 0)")
	}
	fmt.Printf(`const
  %s: array [0 .. %d] of Byte = (%s
);`, *varName, gen.byteCount-1, gen.buf.Bytes())

	if *packageName != "" {
		fmt.Print(`

implementation

end.
`)
	}
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
			fmt.Fprint(&g.buf, "\n    ")
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

package main

import (
	"flag"
	"io"
	"os"

	"github.com/tmtkt/bin2delphi"
)

var (
	unitName  = flag.String("unit", "Main", "Unit name. Must not be empty.")
	constName = flag.String("const", "C", "Constant name. Must not be empty.")
)

func main() {
	flag.Parse()

	if *unitName == "" || *constName == "" {
		flag.Usage()
		return
	}

	data, err := io.ReadAll(os.Stdin)
	check(err)
	u := bin2delphi.NewUnit(*unitName)
	u.AddConstant(*constName, data)
	_, err = os.Stdout.Write(u.Generate())
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

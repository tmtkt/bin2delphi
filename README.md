# Command Line Tool

Install with

	go install github.com/tmtkt/bin2delphi/cmd/bin2delphi@latest

Usage:

	bin2delphi -const="VirusExe" -unit="VirusPayload" < virus.exe > VirusPayload.pas

This will create a file `VirusPayload.pas` and place a byte array in it,
containing the binary data from the file `virus.exe`. Now you can create your
trojan in Delphi.

The tool generates a full Delphi unit like this one:

```
unit VirusPayload;

interface

const
  VirusExe: array [0 .. 6] of Byte = (
    $48, $65, $6C, $6C, $6F, $0D, $0A
  );

implementation

end.
```

The interface relies on stdin and stdout, meaning you can use it with one input
and one output. The example above uses the `<` and `>` operators for this.

If you want to have multiple constants in a single file, you need to use the
Go API.

# Go API

This example reads all `.exe` files in a directory and adds them to a Delphi
file `ExeFiles.pas`. The constants are named with the base file names (without
.exe).

```
package main

import (
	"os"
	"strings"

	"github.com/tmtkt/bin2delphi"
)

func main() {
	unit := bin2delphi.NewUnit("ExeFiles")

	files, _ := os.ReadDir(".")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".exe") {
			data, _ := os.ReadFile(file.Name())
			unit.AddConstant(strings.TrimSuffix(file.Name(), ".exe"), data)
		}
	}

	os.WriteFile("ExeFiles.pas", unit.Generate(), 0666)
}
```

The resulting Delphi file might look like this:

```
unit ExeFiles;

interface

const
  First: array [0 .. 2] of Byte = (
    $61, $62, $63
  );

const
  TheSecond: array [0 .. 2] of Byte = (
    $61, $62, $63
  );

implementation

end.
```

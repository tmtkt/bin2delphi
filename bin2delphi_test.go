package bin2delphi_test

import (
	"strings"
	"testing"

	"github.com/tmtkt/bin2delphi"
	"github.com/gonutz/check"
)

func TestEmptyUnit(t *testing.T) {
	code := bin2delphi.NewUnit("U").Generate()
	check.Eq(t, string(code), strings.Replace(`unit U;

interface

implementation

end.
`, "\n", "\r\n", -1))
}

func TestUnitWithOneConstant(t *testing.T) {
	unit := bin2delphi.NewUnit("U")
	unit.AddConstant("C", []byte{1, 2, 3})
	code := unit.Generate()
	check.Eq(t, string(code), strings.Replace(`unit U;

interface

const
  C: array [0 .. 2] of Byte = (
    $01, $02, $03
  );

implementation

end.
`, "\n", "\r\n", -1))
}

func TestEmptyConstantsAreNilPointers(t *testing.T) {
	unit := bin2delphi.NewUnit("HasEmpty")
	unit.AddConstant("Empty", nil)
	code := unit.Generate()
	check.Eq(t, string(code), strings.Replace(`unit HasEmpty;

interface

const
  Empty: Pointer = nil;

implementation

end.
`, "\n", "\r\n", -1))
}

func TestEveryConstantGetsItsOwnBlock(t *testing.T) {
	unit := bin2delphi.NewUnit("MultipleConstants")
	unit.AddConstant("A", []byte{1})
	unit.AddConstant("Empty", nil)
	unit.AddConstant("B", []byte{2})
	code := unit.Generate()
	check.Eq(t, string(code), strings.Replace(`unit MultipleConstants;

interface

const
  A: array [0 .. 0] of Byte = (
    $01
  );

const
  Empty: Pointer = nil;

const
  B: array [0 .. 0] of Byte = (
    $02
  );

implementation

end.
`, "\n", "\r\n", -1))
}

func TestIfUnitNameIsUTF8_BOMIsPrepended(t *testing.T) {
	unit := bin2delphi.NewUnit("Üü")
	code := unit.Generate()
	want := append([]byte{0xEF, 0xBB, 0xBF}, []byte(strings.Replace(`unit Üü;

interface

implementation

end.
`, "\n", "\r\n", -1))...)
	check.Eq(t, code, want)
}

func TestIfAnyConstantNameIsUTF8_BOMIsPrepended(t *testing.T) {
	unit := bin2delphi.NewUnit("U")
	unit.AddConstant("Öö", []byte{1})
	code := unit.Generate()
	want := append([]byte{0xEF, 0xBB, 0xBF}, []byte(strings.Replace(`unit U;

interface

const
  Öö: array [0 .. 0] of Byte = (
    $01
  );

implementation

end.
`, "\n", "\r\n", -1))...)
	check.Eq(t, code, want)
}

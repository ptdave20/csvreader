package csvreader

import "testing"

var valueHeader = []string{
	"a",
	"b",
	"c",
}
var valueTest = [][]string{
	valueHeader,
	{
		"1",
		"2",
		"true",
	},
}

type headerTestAStruct struct {
	A string `csv:"a"`
	B string `csv:"b"`
	C string `csv:"c"`
}

func TestHeaderA(t *testing.T) {
	var s headerTestAStruct
	csvHeader := GetHeader(valueHeader, s)
	if csvHeader == nil {
		t.Fail()
	}
	if csvHeader.Length() != 3 {
		t.Fail()
	}
	for i := range csvHeader {
		if headers := csvHeader.HeaderValues(i); len(headers) != 1 {
			t.Fail()
		}
	}
}
func TestUnmarshallA(t *testing.T) {
	var s headerTestAStruct
	csvHeader := GetHeader(valueTest[0], s)
	if csvHeader == nil {
		t.Fail()
	}

	if err := UnmarshallRow(csvHeader, valueTest[1], nil, &s); err != nil {
		t.Logf("Failed test because: %s", err)
		t.Fail()
	}

	if s.A != valueTest[1][0] {
		t.Fail()
	}
	if s.B != valueTest[1][1] {
		t.Fail()
	}
	if s.C != valueTest[1][2] {
		t.Fail()
	}
}

type headerTestBStruct struct {
	A string `csv:"a"`
	B int    `csv:"b"`
	C bool   `csv:"c"`
}

func TestUnmarshallB(t *testing.T) {
	var s headerTestBStruct
	csvHeader := GetHeader(valueTest[0], s)
	if csvHeader == nil {
		t.Fail()
	}

	if err := UnmarshallRow(csvHeader, valueTest[1], nil, &s); err != nil {
		t.Logf("Failed test because: %s", err)
		t.Fail()
	}

	if err := UnmarshallRow(csvHeader, valueTest[1], nil, &s); err != nil {
		t.Logf("Failed test because: %s", err)
		t.Fail()
	}

	if s.A != "1" {
		t.Fail()
	}

	if s.B != 2 {
		t.Fail()
	}

	if s.C != true {
		t.Fail()
	}
}

func BenchmarkReader(b *testing.B) {
	// Build data
	var data = make([][]string, 100001)

	data[0] = valueHeader
	for i := 1; i < 100001; i++ {
		data[i] = []string{"c", "123", "true"}
	}

	var s headerTestBStruct
	csvHeader := GetHeader(data[0], s)
	b.StartTimer()
	for i := 1; i < 100001; i++ {
		if err := UnmarshallRow(csvHeader, data[i], nil, &s); err != nil {
			b.Fail()
		}
	}

	b.StopTimer()
}

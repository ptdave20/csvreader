package csvreader

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TAG represents the available header values for a field
const TAG string = "csv"

type csvField struct {
	fieldIndex   int
	kind         reflect.Kind
	headerValues []string
	required     bool
	rowIndex     int
}

var falseValues = []string{
	"f",
	"false",
	"0",
	"inactive",
}

var trueValues = []string{
	"t",
	"true",
	"1",
	"active",
}

type CSVOptions struct {
	DefaultInt    int64
	DefaultUint   uint64
	DefaultBool   bool
	DefaultString string
	TrueValues    []string
	FalseValues   []string
}

func DefaultCSVOptions() *CSVOptions {
	return &CSVOptions{
		DefaultInt:    0,
		DefaultUint:   0,
		DefaultBool:   false,
		DefaultString: "",
		TrueValues:    trueValues,
		FalseValues:   falseValues,
	}
}

type CSVHeader []*csvField

func (c CSVHeader) Length() int {
	return len(c)
}

func (c CSVHeader) HeaderValues(index int) []string {
	return c[index].headerValues
}

func GetHeader(row []string, data interface{}) CSVHeader {
	var t reflect.Type
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	} else {
		t = reflect.TypeOf(data)
	}

	var header CSVHeader

	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)

		tagValue := field.Tag.Get(TAG)
		tagArray := strings.Split(tagValue, ",")
		for i := range tagArray {
			if len(tagArray[i]) > 0 {
				tagArray[i] = strings.TrimSpace(tagArray[i])
			}
		}

		// Omit
		if tagArray[0] == "-" {
			continue
		}
		var newValue = &csvField{
			fieldIndex:   fieldIndex,
			kind:         field.Type.Kind(),
			headerValues: nil,
			required:     false,
			rowIndex:     -1,
		}
		header = append(header, newValue)

		if len(tagArray[0]) == 0 {
			tagArray[0] = field.Name
		}

		for tagArray[len(tagArray)-1] == "required" || tagArray[len(tagArray)-1] == "omitEmpty" {
			lastField := tagArray[len(tagArray)-1]
			switch lastField {
			case "required":
				newValue.required = true
				break
			}
			tagArray = tagArray[:len(tagArray)-1]
		}

		newValue.headerValues = tagArray

		// Find the column that matches
		for headerIndex := range newValue.headerValues {
			if newValue.rowIndex != -1 {
				break
			}
			for rowCol := range row {
				if newValue.headerValues[headerIndex] == row[rowCol] {
					newValue.rowIndex = rowCol
					break
				}
			}
		}
	}

	return header
}

func UnmarshallRow(header CSVHeader, row []string, options *CSVOptions, data interface{}) error {
	if options == nil {
		options = DefaultCSVOptions()
	}
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for _, col := range header {
		if col.rowIndex == -1 {
			continue
		}
		field := v.Field(col.fieldIndex)
		switch col.kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if (len(row) <= col.rowIndex || len(row[col.rowIndex]) == 0) && col.required {
				return errors.New(fmt.Sprintf("Column %d is required and does not have a value", col.rowIndex+1))
			} else if (len(row) <= col.rowIndex || len(row[col.rowIndex]) == 0) && !col.required {
				field.SetInt(options.DefaultInt)
				continue
			}
			value, err := strconv.Atoi(row[col.rowIndex])
			if err != nil {
				return errors.New(fmt.Sprintf("Failed to parse column %d to an int: %s", col.rowIndex+1, err))
			}
			field.SetInt(int64(value))
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if (len(row) <= col.rowIndex || len(row[col.rowIndex]) == 0) && col.required {
				return errors.New(fmt.Sprintf("Column %d is required and does not have a value", col.rowIndex+1))
			} else if (len(row) <= col.rowIndex || len(row[col.rowIndex]) == 0) && !col.required {
				field.SetUint(options.DefaultUint)
				continue
			}
			value, err := strconv.Atoi(row[col.rowIndex])
			if err != nil {
				return errors.New(fmt.Sprintf("Failed to parse column %d to an unsigned int: %s", col.rowIndex+1, err))
			}
			field.SetUint(uint64(value))
			break
		case reflect.String:
			if (len(row) <= col.rowIndex || len(row[col.rowIndex]) == 0) && col.required {
				return errors.New(fmt.Sprintf("Column %d is required and does not have a value", col.rowIndex+1))
			} else if len(row) <= col.rowIndex || len(row[col.rowIndex]) == 0 {
				field.SetString(options.DefaultString)
			} else {
				field.SetString(row[col.rowIndex])
			}
			break
		case reflect.Bool:
			if len(row) > col.rowIndex && isOneOfValue(row[col.rowIndex], trueValues) {
				field.SetBool(true)
			} else if len(row) > col.rowIndex && isOneOfValue(row[col.rowIndex], falseValues) {
				field.SetBool(false)
			} else if col.required {
				return errors.New(fmt.Sprintf("Failed to parse column %d to an boolean", col.rowIndex+1))
			} else {
				field.SetBool(options.DefaultBool)
			}
			break
		}
	}
	return nil
}

func isOneOfValue(value string, comparison []string) bool {
	lowerValue := strings.ToLower(value)
	for index := range comparison {
		if lowerValue == comparison[index] {
			return true
		}
	}
	return false
}

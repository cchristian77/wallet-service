package util

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func GetUnderlyingTypeAndValue(i interface{}) (reflect.Value, reflect.Type, bool) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return v.Elem(), v.Type(), v.IsNil()
		}

		return v.Elem(), v.Elem().Type(), v.IsNil()
	}

	return v, reflect.TypeOf(i), false
}

func CompareData[T any](source, target T, depth uint) error {
	if depth == 0 {
		return nil
	}

	valuesource, typesource, sourceIsNil := GetUnderlyingTypeAndValue(source)
	valuetarget, typetarget, targetIsNil := GetUnderlyingTypeAndValue(target)

	if sourceIsNil && targetIsNil {
		return nil
	}

	// Check if either one is nil
	if sourceIsNil || targetIsNil {
		return fmt.Errorf("one of the value is nil")
	}

	// Make sure both values are of the same type
	if typesource != typetarget {
		return fmt.Errorf("type of source is not matching with target. Source: %v, Target: %v", typesource, typetarget)
	}

	// Iterate through the fields of the struct
	for i := 0; i < typesource.NumField(); i++ {
		field := typesource.Field(i)
		fieldName := field.Name

		// Use reflection to get the field values for both structs
		fieldValuesource := valuesource.FieldByName(fieldName)
		fieldValuetarget := valuetarget.FieldByName(fieldName)

		// If fields are timestamp, use the time.Equal function
		if strings.Contains(fieldValuesource.String(), "time.Time") {
			valsourceTime, _, sourceIsNil := GetUnderlyingTypeAndValue(fieldValuesource.Interface())
			valtargetTime, _, targetIsNil := GetUnderlyingTypeAndValue(fieldValuetarget.Interface())
			if sourceIsNil && targetIsNil {
				continue
			}
			if sourceIsNil != targetIsNil {
				return fmt.Errorf("field %s is not matching, some are nil. Source: %v, Target: %v", fieldName, valsourceTime, valtargetTime)
			}
			if !valsourceTime.Interface().(time.Time).Equal(valtargetTime.Interface().(time.Time)) {
				return fmt.Errorf("field %s is not matching. Source: %v, Target: %v", fieldName, valsourceTime.Interface().(time.Time), valtargetTime.Interface().(time.Time))
			}

			continue
		}

		// If type is decimal, use the decimal function
		if fieldValuesource.Type().String() == "decimal.Decimal" {
			if !fieldValuesource.Interface().(decimal.Decimal).Equal(fieldValuetarget.Interface().(decimal.Decimal)) {
				return fmt.Errorf("field %s is not matching. Source: %v, Target: %v", fieldName, fieldValuesource.Interface().(decimal.Decimal), fieldValuetarget.Interface().(decimal.Decimal))
			}
			continue
		}

		// If fields are struct, recursively call the function
		if fieldValuesource.Kind() == reflect.Struct && depth > 1 {
			if err := CompareData(fieldValuesource.Interface(), fieldValuetarget.Interface(), depth-1); err != nil {
				return err
			}
			continue
		}

		// If fields are array, recursively call the function
		if fieldValuesource.Kind() == reflect.Slice && depth > 1 {
			if depth == 0 {
				return nil
			}

			src := fieldValuesource.Interface()
			tgt := fieldValuetarget.Interface()

			source := reflect.ValueOf(src)
			target := reflect.ValueOf(tgt)

			if source.Len() != target.Len() {
				return fmt.Errorf("source and target has different length. Source: %v, Target: %v", source.Len(), target.Len())
			}

			for i := 0; i < source.Len(); i++ {
				if err := CompareData(source.Index(i).Interface(), target.Index(i).Interface(), depth-1); err != nil {
					return fmt.Errorf("mismatch on array item index {%v}. Details: %v", i, err)
				}
			}
		}

		// Compare the raw field values
		if !reflect.DeepEqual(fieldValuesource.Interface(), fieldValuetarget.Interface()) {
			return fmt.Errorf("field %s is not matching. Source: %v, Target: %v", fieldName, fieldValuesource.Interface(), fieldValuetarget.Interface())
		}
	}

	return nil
}

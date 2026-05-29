package ecommon

import (
	"encoding"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"strings"
)

type paramUnmarshaler interface {
	UnmarshalParam(string) error
}

func makeRequest[TReq any](data interface{}, callback func(TReq)) {
	if data == nil {
		panic("empty request body")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	var parsedReq TReq
	if err := json.Unmarshal(jsonData, &parsedReq); err != nil {
		panic(err)
	}

	validateBoundRequest(parsedReq)
	callback(parsedReq)
}

func makeFormRequest[TReq any](formData map[string]string, callback func(TReq)) {
	values := url.Values{}
	for k, v := range formData {
		values.Set(k, v)
	}
	bindValuesAndRun(values, callback)
}

func makeGetRequest[TReq any](queryParams map[string]string, callback func(TReq)) {
	values := url.Values{}
	for k, v := range queryParams {
		values.Set(k, v)
	}
	bindValuesAndRun(values, callback)
}

func bindValuesAndRun[TReq any](values url.Values, callback func(TReq)) {
	var parsedReq TReq
	if err := bindFormValues(&parsedReq, values); err != nil {
		panic(err)
	}

	validateBoundRequest(parsedReq)
	callback(parsedReq)
}

func bindFormValues(dst interface{}, values url.Values) error {
	value := reflect.ValueOf(dst)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("expected non-nil pointer target")
	}

	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct target")
	}

	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fieldInfo := valueType.Field(i)
		if fieldInfo.PkgPath != "" {
			continue
		}

		name := formFieldName(fieldInfo)
		if name == "" {
			continue
		}

		rawValues, ok := values[name]
		if !ok || len(rawValues) == 0 {
			continue
		}

		if err := setBoundField(value.Field(i), rawValues[0]); err != nil {
			return fmt.Errorf("%s: %w", fieldInfo.Name, err)
		}
	}

	return nil
}

func formFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("form")
	if tag == "-" {
		return ""
	}
	if tag != "" {
		name, _, _ := strings.Cut(tag, ",")
		return name
	}
	return field.Name
}

func setBoundField(field reflect.Value, raw string) error {
	if ok, err := unmarshalField(field, raw); ok {
		return err
	}

	target := field
	if target.Kind() == reflect.Ptr {
		if target.IsNil() {
			target.Set(reflect.New(target.Type().Elem()))
		}
		target = target.Elem()
	}

	if ok, err := unmarshalField(target, raw); ok {
		return err
	}

	switch target.Kind() {
	case reflect.String:
		target.SetString(raw)
		return nil
	default:
		return fmt.Errorf("unsupported field type %s", field.Type())
	}
}

func unmarshalField(field reflect.Value, raw string) (bool, error) {
	if field.Kind() == reflect.Ptr && field.IsNil() {
		field.Set(reflect.New(field.Type().Elem()))
	}

	if field.CanInterface() {
		if unmarshaler, ok := field.Interface().(paramUnmarshaler); ok {
			return true, unmarshaler.UnmarshalParam(raw)
		}
		if unmarshaler, ok := field.Interface().(encoding.TextUnmarshaler); ok {
			return true, unmarshaler.UnmarshalText([]byte(raw))
		}
	}

	if field.CanAddr() {
		addr := field.Addr()
		if addr.CanInterface() {
			if unmarshaler, ok := addr.Interface().(paramUnmarshaler); ok {
				return true, unmarshaler.UnmarshalParam(raw)
			}
			if unmarshaler, ok := addr.Interface().(encoding.TextUnmarshaler); ok {
				return true, unmarshaler.UnmarshalText([]byte(raw))
			}
		}
	}

	return false, nil
}

func validateBoundRequest(req interface{}) {
	if err := validateStructValue(reflect.ValueOf(req)); err != nil {
		panic(err)
	}
}

func validateStructValue(value reflect.Value) error {
	if !value.IsValid() {
		return nil
	}
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil
	}

	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fieldInfo := valueType.Field(i)
		if fieldInfo.PkgPath != "" {
			continue
		}

		rules := fieldInfo.Tag.Get("binding")
		if rules == "" || rules == "-" {
			continue
		}

		field := value.Field(i)
		for _, rule := range strings.Split(rules, ",") {
			switch rule {
			case "required":
				if isZeroValue(field) {
					return fmt.Errorf("%s is required", fieldInfo.Name)
				}
			case "positive_bigint":
				if !validatePositiveBigIntValue(field) {
					return fmt.Errorf("%s must be a positive bigint", fieldInfo.Name)
				}
			}
		}
	}

	return nil
}

func isZeroValue(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}
	return value.IsZero()
}

func validatePositiveBigIntValue(field reflect.Value) bool {
	if !field.IsValid() || !field.CanInterface() {
		return false
	}

	for field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return false
		}
		field = field.Elem()
	}

	if field.Kind() == reflect.String {
		return validateStringAsBigInt(field.String())
	}

	if !field.CanInterface() {
		return false
	}

	switch v := field.Interface().(type) {
	case string:
		return validateStringAsBigInt(v)
	case Big:
		return v.ToInt().Sign() > 0
	default:
		return false
	}
}

func validateStringAsBigInt(s string) bool {
	if s == "" {
		return false
	}

	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}

	if s[0] != '+' && (s[0] < '0' || s[0] > '9') {
		return false
	}

	if s[0] == '+' {
		s = s[1:]
		if s == "" {
			return false
		}
	}

	if s == "0" || (len(s) > 1 && s[0] == '0') {
		return false
	}

	if len(s) <= 18 {
		for _, c := range s {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	}

	num, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return false
	}
	return num.Sign() > 0
}

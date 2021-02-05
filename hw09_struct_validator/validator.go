package hw09_struct_validator //nolint:golint,stylecheck
import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var ErrInvalidString = errors.New("invalid string")

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(iv interface{}) error {
	v := reflect.ValueOf(iv)

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("%T is not a pointer to struct", iv)
	}

	fmt.Println("v.NumField() ", v.NumField())

	z := reflect.TypeOf(iv)
	fmt.Println("reflect.TypeOf", z.Name())
	fmt.Printf("%T: %v\n", z, z)

	t := v.Type()
	fmt.Printf("%T", t)

	mp := make(map[string]interface{}, t.NumField())
	for i := 0; i < t.NumField(); i++ {

		field := t.Field(i) // reflect.StructField
		fmt.Println("\n\nnew elem:")
		fmt.Println("type of value field", field.Type.String())
		fmt.Println("field.Name", field.Name)
		fmt.Println("FieldByName value: ", v.FieldByName(field.Name))
		g := v.FieldByName(field.Name)
		fmt.Println("g ", g)

		p := v.Field(i)

		switch {
		case p.Kind() == reflect.String:
			typeSwitch(p.Interface())
			fmt.Println("KIND is string")
			tag := field.Tag.Get("validate")
			strings.HasPrefix(tag, "len:")
			strings.HasPrefix(tag, "in:")
			strings.HasPrefix(tag, "regexp:")
		case p.Kind() == reflect.Int:
			typeSwitch(p.Interface())
			fmt.Println("KIND is int")
			tag := field.Tag.Get("validate")
			strings.HasPrefix(tag, "min:")
			strings.Contains(tag, "min|max")
		case p.Kind() == reflect.Slice:
			typeSwitch(p.Interface())
			fmt.Println("KIND is slice")

		}

		j := p.Interface()

		switch j.(type) {
		case string:
			fmt.Println("value is string")
		case int:
			fmt.Println("value is int")

		}
		fmt.Println("tag", field.Tag.Get("validate"))
		fv := v.Field(i) // reflect.Value

		mp[field.Name] = fv.Interface()
	}
	return nil
}

func typeSwitch(val interface{}) {
	switch h := val.(type) {
	case int:
		fmt.Println("!!!!!!int with value", val)
	case string:
		fmt.Println("!!!!!!!string with value ", val)
	case []string:
		for i, v := range h {
			fmt.Printf("element %d: %s", i, v)
		}
		fmt.Println("!!!!!!!Slice of string with value", val)
	default:
		fmt.Println("!!!!!!!Unhandled", "with value", val)
	}
}

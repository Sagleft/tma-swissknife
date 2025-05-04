package config

import (
	"fmt"
	"reflect"

	"github.com/kelseyhightower/envconfig"
)

/*
ProcessConfig iterates through the fields
of the structure and loads envconfig for each one.
*/
func ProcessConfig(cfg any) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("cfg must be a non-nil pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("cfg must be a pointer to a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanAddr() {
			return fmt.Errorf(
				"cannot obrain address of config field %q",
				field.Name,
			)
		}

		if err := envconfig.Process("", fieldValue.Addr().Interface()); err != nil {
			return fmt.Errorf("parse config for %s: %w", field.Name, err)
		}
	}

	return nil
}

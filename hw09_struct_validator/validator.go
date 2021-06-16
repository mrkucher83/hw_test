package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %v", v.Field, v.Err.Error())
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var b strings.Builder
	for _, er := range v {
		fmt.Fprintf(&b, "%s: %v\n", er.Field, er.Err)
	}
	return b.String()
}

var (
	ErrUnknownTag = errors.New("got unknown tag of validation")
	ErrLength     = errors.New("length should be equal")
)

func Validate(v interface{}) error {
	var errs ValidationErrors

	vr := reflect.ValueOf(v)
	if vr.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, but got %T", v)
	}

	vt := vr.Type()

	for i := 0; i < vr.NumField(); i++ {
		err := ValidateField(vt.Field(i), vr.Field(i), &errs)
		if err != nil {
			return err
		}
	}

	return errs
}

func ValidateField(vt reflect.StructField, vr reflect.Value, errs *ValidationErrors) error {
	tag := vt.Tag.Get("validate")

	if len(tag) == 0 || tag == "-" {
		return nil
	}

	tagValues := strings.Split(tag, "|")
	for _, tagVal := range tagValues {
		if vr.Kind() == reflect.Slice {
			for j := 0; j < vr.Len(); j++ {
				if err := ValidateValue(vr.Index(j), vt, tagVal, errs); err != nil {
					return err
				}
			}
		} else {
			err := ValidateValue(vr, vt, tagVal, errs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidateValue(vr reflect.Value, vt reflect.StructField, val string, errs *ValidationErrors) error {
	tagValue := strings.Split(val, ":")
	num, _ := strconv.Atoi(tagValue[1])

	if vr.Kind() == reflect.String {
		switch tagValue[0] {
		case "len":
			if vr.Len() != num {
				*errs = append(*errs, ValidationError{
					vt.Name, fmt.Errorf("%w: %d", ErrLength, num),
				})
			}
		case "regexp":
			re, _ := regexp.Compile(tagValue[1])
			if !re.MatchString(vr.Interface().(string)) {
				*errs = append(*errs, ValidationError{
					vt.Name, fmt.Errorf("should contain certain symbols: %v", re),
				})
			}
		case "in":
			args := strings.Split(tagValue[1], ",")
			var isContain bool
			for _, val := range args {
				is := &isContain
				if vr.Interface() == val {
					*is = true
					break
				}
			}
			if !isContain {
				*errs = append(*errs, ValidationError{
					vt.Name, fmt.Errorf("should be equal one of: %v", tagValue[1]),
				})
			}

		default:
			return fmt.Errorf("%w: %q for: %s", ErrUnknownTag, tagValue[0], vt.Name)
		}
	}

	if vr.Kind() == reflect.Int {
		switch tagValue[0] {
		case "min":
			if vr.Interface().(int) < num {
				*errs = append(*errs, ValidationError{
					vt.Name, fmt.Errorf("should be equal or greater: %d", num),
				})
			}
		case "max":
			if vr.Interface().(int) > num {
				*errs = append(*errs, ValidationError{
					vt.Name, fmt.Errorf("should be equal or less: %d", num),
				})
			}
		case "in":
			args := strings.Split(tagValue[1], ",")
			var isContain bool
			for _, val := range args {
				is := &isContain
				num, _ := strconv.Atoi(val)
				if vr.Interface() == num {
					*is = true
					break
				}
			}
			if !isContain {
				*errs = append(*errs, ValidationError{
					vt.Name, fmt.Errorf("should be equal one of: %v", tagValue[1]),
				})
			}

		default:
			return fmt.Errorf("%w: %q for: %s", ErrUnknownTag, tagValue[0], vt.Name)
		}
	}

	return nil
}

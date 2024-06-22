package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

// more performant, because compiling it once at runtime
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func New(data url.Values) *Form {
	return &Form{
		data, errors(map[string][]string{}),
	}
}

func (form *Form) Required(fields ...string) {
	for _, field := range fields {
		if strings.TrimSpace(form.Values.Get(field)) == "" {
			form.Errors.Add(field, "this "+field+" is required.")
		}
	}
}

func (form *Form) MaxLength(max int, fields ...string) {
	for _, field := range fields {
		value := form.Values.Get(field)
		if value == "" {
			return
		}

		if utf8.RuneCountInString(value) > max {
			form.Errors.Add(field, fmt.Sprintf("%s cannot be more than %d characters.", field, max))
		}
	}
}

func (form *Form) MinLength(min int, fields ...string) {
	for _, field := range fields {
		value := form.Values.Get(field)
		if value == "" {
			return
		}

		if utf8.RuneCountInString(value) < min {
			form.Errors.Add(field, fmt.Sprintf("%s has to be atleast %d characters long.", field, min))
		}
	}
}

func (form *Form) MatchesPattern(pattern *regexp.Regexp, fields ...string) {
	for _, field := range fields {
		value := form.Values.Get(field)
		if value == "" {
			return
		}

		if !pattern.MatchString(value) {
			form.Errors.Add(field, fmt.Sprintf("%s is invalid.", field))
		}
	}
}

func (form *Form) PermittedValues(field string, opts ...string) {
	value := form.Values.Get(field)
	if value == "" {
		return
	}

	for _, opt := range opts {
		if value == opt {
			return
		}
	}

	form.Errors.Add(field, fmt.Sprintf("%s value is invalid.", field))
}

func (form *Form) IsEqual(field1, field2 string) {
	value1 := form.Values.Get(field1)
	if value1 == "" {
		return
	}

	value2 := form.Values.Get(field2)
	if value2 == "" {
		return
	}

	if value1 != value2 {
		form.Errors.Add(field2, fmt.Sprintf("%s and %s don't match.", field1, field2))
	}
}

func (form *Form) Valid() bool {
	return len(form.Errors) == 0
}

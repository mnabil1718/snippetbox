package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

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

func (form *Form) Valid() bool {
	return len(form.Errors) == 0
}

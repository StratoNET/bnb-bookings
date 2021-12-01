package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form is a custom form struct for any given form, holding associated errors (errors.go) & embeds an url.Values object
type Form struct {
	url.Values
	Errors errors
}

// NewForm initialises a Form struct, which includes values from within any form fields (data)
func NewForm(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// begin validation check functions...

// ValidForm returns true if there are no errors in form data, otherwise false
func (f *Form) ValidForm() bool {
	return len(f.Errors) == 0
}

// RequiredFields checks if form has given required fields in post request & that they're not empty
func (f *Form) RequiredFields(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.AddErrMsg(field, "this field is required and cannot be empty")
		}
	}
}

// HasField checks if form has given required field in post request & that it's not empty
func (f *Form) HasField(field string) bool {
	fld := f.Get(field)
	return fld != ""
	// if fld == "" {
	// 	return false
	// }
	// return true
}

// MinLength checks field has minimum required characters
func (f *Form) MinLength(field string, length int) bool {
	fld := f.Get(field)
	if len(fld) < length {
		f.Errors.AddErrMsg(field, fmt.Sprintf("this field must be at least %d characters in length", length))
		return false
	}
	return true
}

//
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.AddErrMsg(field, "please input a valid email address e.g. your.name@example.com")
	}

}

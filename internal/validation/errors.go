package forms

type errors map[string][]string

// AddErrMsg appends an error message to the slice of error messages, contained in errors map, for any given field
func (e errors) AddErrMsg(field, errmsg string) {
	e[field] = append(e[field], errmsg)
}

// GetErrMsg returns first error message (or "") in slice of error messages, contained in errors map, for any given field
func (e errors) GetErrMsg(field string) string {
	errmsg := e[field]
	if len(errmsg) == 0 {
		return ""
	}
	return errmsg[0]
}

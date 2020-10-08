package forms

// type errors hold validation error messagees for forms
type errors map[string][]string

// Add function adds new error field with its message into errors map object.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get method retreives first error message for given field
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

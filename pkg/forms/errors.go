package forms

type errors map[string][]string

func (err errors) Add(field string, message string) {
	err[field] = append(err[field], message)
}

func (err errors) Get(field string) string {
	es := err[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

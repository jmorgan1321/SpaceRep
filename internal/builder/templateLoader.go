package builder

type tmplLoader interface {
	LoadTemplates([]string) error
}

type tmplMapLoader struct {
}

func (tm *tmplMapLoader) LoadTemplates([]string) error {
	return nil
}

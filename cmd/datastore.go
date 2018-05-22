package cmd

import "cloud.google.com/go/datastore"

type Entity struct {
	Props map[string]interface{}
}

func (e *Entity) Load(ps []datastore.Property) error {
	err := datastore.LoadStruct(e, ps)

	if fmerr, ok := err.(*datastore.ErrFieldMismatch); ok && fmerr != nil && fmerr.Reason == "no such struct field" {
	} else if err != nil {
		return err
	}

	e.Props = map[string]interface{}{}
	for _, p := range ps {
		e.Props[p.Name] = p.Value
	}

	return nil
}

func (e *Entity) Save() ([]datastore.Property, error) {
	pr, err := datastore.SaveStruct(e)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

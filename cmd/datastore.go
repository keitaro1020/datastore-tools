package cmd

import (
	"cloud.google.com/go/datastore"
	"context"
	"google.golang.org/api/option"
)

type Entity struct {
	Props map[string]interface{}
}

type JsonKey struct {
	Kind      string
	ID        int64
	Name      string
	Namespace string
}

func NewJsonKey(key *datastore.Key) *JsonKey {
	return &JsonKey{
		Kind:      key.Kind,
		ID:        key.ID,
		Name:      key.Name,
		Namespace: key.Namespace,
	}
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

func NewDatastoreClient(c context.Context, keyfile, project string) (*datastore.Client, error) {
	var opts []option.ClientOption
	if keyfile != "" {
		opts = []option.ClientOption{
			option.WithCredentialsFile(keyfile),
		}
	}

	return datastore.NewClient(c, project, opts...)
}

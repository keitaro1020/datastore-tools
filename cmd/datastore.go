package cmd

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"strings"
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

type DatastoreClient struct {
	Client *datastore.Client
}

func NewDatastoreClient(c context.Context, keyfile, project string) (*DatastoreClient, error) {
	var opts []option.ClientOption
	if keyfile != "" {
		opts = []option.ClientOption{
			option.WithCredentialsFile(keyfile),
		}
	}
	client, err := datastore.NewClient(c, project, opts...)
	if err != nil {
		return nil, err
	}

	return &DatastoreClient{
		Client: client,
	}, nil
}

func (c *DatastoreClient) GetQuery(kind, namespace, filter string, keysOnly bool) (*datastore.Query, error) {
	query := datastore.NewQuery(kind)
	if namespace != "" {
		query = query.Namespace(namespace)
	}
	if filter != "" {
		filter := strings.SplitN(filter, "=", 2)
		if len(filter) == 2 {
			if filter[0] == "__key__" {
				key := datastore.Key{
					Kind: kind,
					Name: filter[1],
				}
				if namespace != "" {
					key.Namespace = namespace
				}
				query = query.Filter(fmt.Sprintf("%s =", filter[0]), &key)
			} else {
				query = query.Filter(fmt.Sprintf("%s =", filter[0]), filter[1])
			}
		} else {
			return nil, fmt.Errorf("error: invalid filter parameter: %s", o.OptFilter)
		}
	}
	if keysOnly {
		query = query.KeysOnly()
	}
	return query, nil
}

func (c *DatastoreClient) GetAll(ctx context.Context, query *datastore.Query, entities *[]Entity) ([]*datastore.Key, error) {
	return c.Client.GetAll(ctx, query, entities)
}

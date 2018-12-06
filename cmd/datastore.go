package cmd

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

type Entity struct {
	Props      map[string]interface{}
	Properties []datastore.Property
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

	e.Properties = ps
	e.Props = map[string]interface{}{}
	for _, p := range ps {
		e.Props[p.Name] = p.Value
	}

	return nil
}

func (e *Entity) Save() ([]datastore.Property, error) {
	return e.Properties, nil
}

func (e *Entity) SetValue(name, value string) error {
	res := false
	for i, prop := range e.Properties {
		if prop.Name == name {
			var setVal interface{}
			switch prop.Value.(type) {
			case string:
				setVal = value
			case int, int8, int16, int32, int64:
				v, _ := strconv.Atoi(value)
				setVal = v
			case float32:
				v, _ := strconv.ParseFloat(value, 32)
				setVal = v
			case float64:
				v, _ := strconv.ParseFloat(value, 64)
				setVal = v
			default:
				return fmt.Errorf("unsupported property type: %#v", reflect.ValueOf(prop.Value))
			}

			e.Properties[i] = datastore.Property{
				Name:    prop.Name,
				Value:   setVal,
				NoIndex: prop.NoIndex,
			}
			res = true
		}
	}
	if res {
		return nil
	} else {
		return fmt.Errorf("unknown property name: %s", name)
	}
}

type DatastoreClient interface {
	GetQuery(kind, namespace, filter string, keysOnly bool) (*datastore.Query, error)
	GetAll(ctx context.Context, query *datastore.Query, entities *[]Entity) ([]*datastore.Key, error)
	PutMulti(ctx context.Context, keys []*datastore.Key, entities []Entity) ([]*datastore.Key, error)
	DeleteMulti(ctx context.Context, keys []*datastore.Key) error
}

type DatastoreClientImpl struct {
	Client *datastore.Client
}

func NewDatastoreClient(c context.Context, keyfile, project string) (DatastoreClient, error) {
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

	var cli DatastoreClient = &DatastoreClientImpl{
		Client: client,
	}
	return cli, nil
}

func (c *DatastoreClientImpl) GetQuery(kind, namespace, filter string, keysOnly bool) (*datastore.Query, error) {
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

func (c *DatastoreClientImpl) GetAll(ctx context.Context, query *datastore.Query, entities *[]Entity) ([]*datastore.Key, error) {
	return c.Client.GetAll(ctx, query, entities)
}

func (c *DatastoreClientImpl) PutMulti(ctx context.Context, keys []*datastore.Key, entities []Entity) ([]*datastore.Key, error) {
	return c.Client.PutMulti(ctx, keys, entities)
}

func (c *DatastoreClientImpl) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	return c.Client.DeleteMulti(ctx, keys)
}

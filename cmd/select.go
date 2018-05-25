package cmd

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

func newSelectCmd() *cobra.Command {
	type Options struct {
		OptProject   string
		OptKind      string
		OptNamespace string
		OptKeyFile   string
		OptWhere     string
		OptCount     bool
		OptTable     bool
	}

	var (
		o = &Options{}
	)

	cmd := &cobra.Command{
		Use:   "select",
		Short: "Select Entity",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			client, err := NewDatastoreClient(ctx, o.OptKeyFile, o.OptProject)
			if err != nil {
				return err
			}

			query := datastore.NewQuery(o.OptKind)
			if o.OptNamespace != "" {
				query = query.Namespace(o.OptNamespace)
			}
			if o.OptWhere != "" {
				where := strings.SplitN(o.OptWhere, "=", 2)
				if len(where) == 2 {
					if where[0] == "__key__" {
						key := datastore.Key{
							Kind: o.OptKind,
							Name: where[1],
						}
						if o.OptNamespace != "" {
							key.Namespace = o.OptNamespace
						}
						query = query.Filter(fmt.Sprintf("%s =", where[0]), &key)
					} else {
						query = query.Filter(fmt.Sprintf("%s =", where[0]), where[1])
					}
				} else {
					return fmt.Errorf("error: invalid where parameter: %s", o.OptWhere)
				}
			}
			if o.OptCount {
				query = query.KeysOnly()
			}

			var entities []Entity

			keys, err := client.GetAll(ctx, query, &entities)
			if err != nil {
				return err
			}

			if !o.OptCount {
				if o.OptTable {
					outputTable(cmd.OutOrStdout(), keys, entities)
				} else {
					outputJson(cmd.OutOrStdout(), keys, entities)
				}
			}
			cmd.Printf("count: %d \n", len(keys))

			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  false,
	}
	cmd.Flags().StringVarP(&o.OptProject, "project", "p", "", "datastore project id [required]")
	cmd.Flags().StringVarP(&o.OptKind, "kind", "k", "", "datastore kind [required]")
	cmd.Flags().StringVarP(&o.OptNamespace, "namespace", "n", "", "datastore namespace")
	cmd.Flags().StringVarP(&o.OptKeyFile, "key-file", "f", "", "gcp service account JSON key file")
	cmd.Flags().StringVarP(&o.OptWhere, "where", "w", "", "query (Property=Value)")
	cmd.Flags().BoolVarP(&o.OptCount, "count", "c", false, "count only")
	cmd.Flags().BoolVarP(&o.OptTable, "table", "t", false, "output table view")

	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("kind")

	return cmd
}

func init() {
	RootCmd.AddCommand(newSelectCmd())
}

func outputTable(w io.Writer, keys []*datastore.Key, entities []Entity) {
	if len(keys) > 0 {
		headers := []string{"__key__"}
		for propKey := range entities[0].Props {
			headers = append(headers, propKey)
		}
		table := tablewriter.NewWriter(w)
		table.SetHeader(headers)
		table.SetRowLine(true)

		for i, key := range keys {
			entity := entities[i]
			entity.Props["__key__"] = key

			var row []string
			for _, header := range headers {
				if v, ok := entity.Props[header]; ok {
					switch tv := v.(type) {
					case *datastore.Key:
						if tv.ID != 0 {
							row = append(row, fmt.Sprintf("%v", tv.ID))
						} else {
							row = append(row, tv.Name)
						}
					default:
						row = append(row, fmt.Sprintf("%v", tv))
					}
				}
			}
			table.Append(row)
		}
		table.Render()
	}
}

func outputJson(w io.Writer, keys []*datastore.Key, entities []Entity) {
	for i, key := range keys {
		entity := entities[i]
		entity.Props["__key__"] = &JsonKey{
			Kind:      key.Kind,
			ID:        key.ID,
			Name:      key.Name,
			Namespace: key.Namespace,
		}

		for k, v := range entity.Props {
			switch vc := v.(type) {
			case *datastore.Key:
				entity.Props[k] = &JsonKey{
					Kind:      vc.Kind,
					ID:        vc.ID,
					Name:      vc.Name,
					Namespace: vc.Namespace,
				}
			}
		}

		j, _ := json.Marshal(entity.Props)
		var ij bytes.Buffer
		json.Indent(&ij, j, "", "  ")
		js := ij.String()
		if len(keys) > i+1 {
			js += ","
		}
		fmt.Fprintf(w, "%s\n", js)
	}
}

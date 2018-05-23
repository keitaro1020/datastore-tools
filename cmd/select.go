package cmd

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/olekukonko/tablewriter"
	"fmt"
)

func newSelectCmd() *cobra.Command {
	type Options struct {
		OptProject   string
		OptKind      string
		OptNamespace string
		OptKeyFile   string
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
					outputTable(cmd, keys, entities)
				} else {
					outputJson(cmd, keys, entities)
				}
			}
			cmd.Printf("count: %d \n", len(keys))

			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.Flags().StringVarP(&o.OptProject, "project", "p", "", "datastore project id [required]")
	cmd.Flags().StringVarP(&o.OptKind, "kind", "k", "", "datastore kind [required]")
	cmd.Flags().StringVarP(&o.OptNamespace, "namespace", "n", "", "datastore namespace")
	cmd.Flags().StringVarP(&o.OptKeyFile, "key-file", "f", "", "gcp service account JSON key file [required]")
	cmd.Flags().BoolVarP(&o.OptCount, "count", "c", false, "count only")
	cmd.Flags().BoolVarP(&o.OptTable, "table", "t", false, "output table view")

	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("kind")
	cmd.MarkFlagRequired("key-file")

	return cmd
}

func init() {
	RootCmd.AddCommand(newSelectCmd())
}

func outputTable(cmd *cobra.Command, keys []*datastore.Key, entities []Entity) {
	if len(keys) > 0 {
		headers := []string{"__key__"}
		for propKey := range entities[0].Props {
			headers = append(headers, propKey)
		}
		table := tablewriter.NewWriter(cmd.OutOrStdout())
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

func outputJson(cmd *cobra.Command, keys []*datastore.Key, entities []Entity) {
	for i, key := range keys {
		entity := entities[i]
		entity.Props["__key__"] = key

		j, _ := json.Marshal(entity.Props)
		var ij bytes.Buffer
		json.Indent(&ij, j, "", "  ")
		js := ij.String()
		if len(keys) > i+1 {
			js += ","
		}
		cmd.Printf("%s\n", js)
	}
}
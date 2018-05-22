package cmd

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

func newSelectCmd() *cobra.Command {
	type Options struct {
		OptProject   string
		OptKind      string
		OptNamespace string
		OptKeyFile   string
		OptCount     bool
	}

	var (
		o = &Options{}
	)

	cmd := &cobra.Command{
		Use:   "select",
		Short: "Select Entity",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			opts := []option.ClientOption{
				option.WithCredentialsFile(o.OptKeyFile),
			}

			client, err := datastore.NewClient(ctx, o.OptProject, opts...)
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
			cmd.Printf("count: %d \n", len(keys))

			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.Flags().StringVarP(&o.OptProject, "project", "p", "", "datastore project id [required]")
	cmd.Flags().StringVarP(&o.OptKind, "kind", "k", "", "datastore kind [required]")
	cmd.Flags().StringVarP(&o.OptNamespace, "namespace", "n", "", "datastore namespace")
	cmd.Flags().StringVarP(&o.OptKeyFile, "key-file", "f", "", "gcp service account JSON key file")
	cmd.Flags().BoolVarP(&o.OptCount, "count", "c", false, "count only")

	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("kind")
	cmd.MarkFlagRequired("key-file")

	return cmd
}

func init() {
	RootCmd.AddCommand(newSelectCmd())
}

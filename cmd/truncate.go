package cmd

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/spf13/cobra"
)

func newTruncateCmd() *cobra.Command {
	type Options struct {
		OptProject   string
		OptKind      string
		OptKeyFile   string
		OptNamespace string
	}

	var (
		o = &Options{}
	)

	cmd := &cobra.Command{
		Use:   "truncate",
		Short: "Truncate Entity",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			client, err := NewDatastoreClient(ctx, o.OptKeyFile, o.OptProject)
			if err != nil {
				return err
			}

			query := datastore.NewQuery(o.OptKind).KeysOnly()
			if o.OptNamespace != "" {
				query = query.Namespace(o.OptNamespace)
			}

			var entities []Entity

			keys, err := client.GetAll(ctx, query, &entities)
			if err != nil {
				return err
			}

			count := len(keys)

			tKeys, keys := slice(keys, 0, 500)
			for len(tKeys) > 0 {
				if err := client.DeleteMulti(ctx, tKeys); err != nil {
					return err
				}

				tKeys, keys = slice(keys, 0, 500)
			}
			cmd.Printf("truncate finish, count = %d \n", count)

			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  false,
	}
	cmd.Flags().StringVarP(&o.OptProject, "project", "p", "", "datastore project id [required]")
	cmd.Flags().StringVarP(&o.OptKind, "kind", "k", "", "datastore kind [required]")
	cmd.Flags().StringVarP(&o.OptNamespace, "namespace", "n", "", "datastore namespace")
	cmd.Flags().StringVarP(&o.OptKeyFile, "key-file", "f", "", "gcp service account JSON key file")

	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("kind")

	return cmd
}

func init() {
	RootCmd.AddCommand(newTruncateCmd())
}

func slice(slice []*datastore.Key, start, end int) ([]*datastore.Key, []*datastore.Key) {
	if len(slice) < start || len(slice) < end {
		return slice, nil
	}
	ans := make([]*datastore.Key, (end - start))
	copy(ans, slice[start:end])
	slice = append(slice[:start], slice[end:]...)
	return ans, slice
}

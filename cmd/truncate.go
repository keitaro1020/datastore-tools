package cmd

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/spf13/cobra"
	"sync"
)

func newTruncateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "truncate",
		Short:         "Truncate Entity",
		RunE:          truncateFunction,
		SilenceErrors: true,
		SilenceUsage:  false,
	}
	cmd.Flags().StringVarP(&o.OptProject, "project", "p", "", "datastore project id [required]")
	cmd.Flags().StringVarP(&o.OptKind, "kind", "k", "", "datastore kind [required]")
	cmd.Flags().StringVarP(&o.OptNamespace, "namespace", "n", "", "datastore namespace")
	cmd.Flags().StringVarP(&o.OptKeyFile, "key-file", "f", "", "gcp service account JSON key file")
	cmd.Flags().StringVarP(&o.OptFilter, "filter", "w", "", "filter query (Property=Value)")

	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("kind")

	return cmd
}

func init() {
	RootCmd.AddCommand(newTruncateCmd())
}

func truncateFunction(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := NewDatastoreClient(ctx, o.OptKeyFile, o.OptProject)
	if err != nil {
		return err
	}

	query, err := client.GetQuery(o.OptKind, o.OptNamespace, o.OptFilter, true)
	if err != nil {
		return err
	}

	var entities []Entity
	keys, err := client.GetAll(ctx, query, &entities)
	if err != nil {
		return err
	}

	count := len(keys)
	wg := &sync.WaitGroup{}

	for _, ks := range slice(keys, 500) {
		wg.Add(1)
		go func(keys []*datastore.Key) {
			defer wg.Done()
			if err := client.DeleteMulti(ctx, keys); err != nil {
				cmd.Printf("delete error: %+v", err)
			}
		}(ks)
	}
	wg.Wait()
	cmd.Printf("truncate finish, count = %d \n", count)

	return nil
}

func slice(slice []*datastore.Key, size int) [][]*datastore.Key {
	if len(slice) <= size {
		return [][]*datastore.Key{slice}
	}

	resSize := len(slice) / size
	if len(slice)%size > 0 {
		resSize += 1
	}

	res := make([][]*datastore.Key, resSize)

	for i := range res {
		if len(slice) > size {
			rt := make([]*datastore.Key, size)
			copy(rt, slice[0:size])
			slice = append(slice[:0], slice[size:]...)
			res[i] = rt
		} else {
			res[i] = slice
		}
	}

	return res
}

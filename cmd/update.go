package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "update",
		Short:         "Update Entity",
		RunE:          updateFunction,
		SilenceErrors: true,
		SilenceUsage:  false,
	}
	cmd.Flags().StringVarP(&o.OptProject, "project", "p", "", "datastore project id [required]")
	cmd.Flags().StringVarP(&o.OptKind, "kind", "k", "", "datastore kind [required]")
	cmd.Flags().StringVarP(&o.OptNamespace, "namespace", "n", "", "datastore namespace")
	cmd.Flags().StringVarP(&o.OptKeyFile, "key-file", "f", "", "gcp service account JSON key file")
	cmd.Flags().StringVarP(&o.OptFilter, "filter", "w", "", "filter query (Property=Value)")
	cmd.Flags().StringVarP(&o.OptSet, "set", "s", "", "set update value (Property=Value)")

	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("kind")
	cmd.MarkFlagRequired("set")

	return cmd
}

func init() {
	RootCmd.AddCommand(newUpdateCmd())
}

func updateFunction(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	set := strings.SplitN(o.OptSet, "=", 2)
	if len(set) < 2 {
		return fmt.Errorf("error: invalid set parameter: %s", o.OptSet)
	}

	client, err := NewDatastoreClient(ctx, o.OptKeyFile, o.OptProject)
	if err != nil {
		return err
	}

	query, err := client.GetQuery(o.OptKind, o.OptNamespace, o.OptFilter, false)
	if err != nil {
		return err
	}

	var entities []Entity
	keys, err := client.GetAll(ctx, query, &entities)
	if err != nil {
		return err
	}

	for _, entity := range entities {
		if err := entity.SetValue(set[0], set[1]); err != nil {
			return err
		}
	}

	if _, err := client.PutMulti(ctx, keys, entities); err != nil {
		return err
	}

	return nil
}

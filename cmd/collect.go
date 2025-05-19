package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/secberus/meraki-collector/collector"
	"github.com/secberus/meraki-collector/config"
	"github.com/secberus/meraki-collector/meraki"
	"github.com/secberus/meraki-collector/push"
	"github.com/secberus/meraki-collector/resource"
)

var collectCmd = &cobra.Command{
	Use:          "collect [-r/--dry-run]",
	Short:        "Collect resources from the Meraki Dashboard API",
	Long:         "Collect resources from the Meraki Dashboard API.",
	Args:         cobra.NoArgs,
	RunE:         collect,
	SilenceUsage: true,
}

var dryrun bool

func init() {
	collectCmd.Flags().BoolVar(&dryrun, "dry-run", false, "dry run?")
}

func collect(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %s", err)
	}

	meraki, err := meraki.Init(&cfg.Meraki)
	if err != nil {
		return fmt.Errorf("failed to initialize Meraki client: %s", err)
	}

	pushsvc, err := push.Init(cfg.Push)
	if err != nil {
		return fmt.Errorf("failed to initialize Push client: %s", err)
	}

	collector := collector.New(meraki, pushsvc, collector.WithDryRun(dryrun))

	// collect from Meraki API root (organizations)
	if err := collector.Collect(cmd.Context(), resource.Organizations); err != nil {
		return fmt.Errorf("failed to collect from Meraki API: %s", err)
	}

	return nil
}

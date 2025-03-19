package cmd

import (
	"github.com/secberus/meraki-collector/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Root = &cobra.Command{
	Use: "meraki-collector <command>",
}

func init() {
	pf := Root.PersistentFlags()

	pf.StringP("config", "c", "", "Configuration file (default: ./config.yaml, $HOME/.s6s/config.yaml)")
	viper.BindPFlag(config.ConfigFileKey, pf.Lookup("config"))

	pf.StringP("bundle", "b", "", "Datasource credentials bundle file (overrides configuration file) (default: $HOME/.s6s/bundle.json)")
	viper.BindPFlag(config.BundleFileKey, pf.Lookup("bundle"))
	viper.SetDefault(config.BundleFileKey, config.DefaultBundleFile)

	pf.StringP("endpoint", "e", "", "Secberus Push API endpoint (overrides configuration & bundle files)")
	viper.BindPFlag(config.PushEndpointKey, pf.Lookup("push-endpoint"))

	pf.StringP("meraki-base-url", "u", "", "Meraki Dashboard API base URL (default: https://api.meraki.com/)")
	viper.BindPFlag(config.MerakiBaseUrlKey, pf.Lookup("meraki-base-url"))
	viper.SetDefault(config.MerakiBaseUrlKey, config.DefaultMerakiBaseUrl)

	pf.StringP("meraki-api-key", "k", "", "Meraki Dashboard API key")
	viper.BindPFlag(config.MerakiApiKeyKey, pf.Lookup("meraki-api-key"))

	pf.Bool("meraki-debug", false, "Meraki Dashboard SDK debugging")
	viper.BindPFlag(config.MerakiDebugKey, pf.Lookup("meraki-debug"))

	Root.AddCommand(collectCmd)
}

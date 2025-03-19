/*
 * Copyright 2025 Secberus, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ConfigFileKey        = "config-file"
	BundleFileKey        = "bundle-file"
	PushEndpointKey      = "push.endpoint"
	MerakiBaseUrlKey     = "meraki.base-url"
	MerakiApiKeyKey      = "meraki.api-key"
	MerakiDebugKey       = "meraki.debug"
	DefaultBundleFile    = "$HOME/.s6s/bundle.json"
	DefaultPushEndpoint  = "push.secberus.io:7744"
	DefaultMerakiBaseUrl = "https://api.meraki.com/"
)

type PushConfig struct {
	Endpoint        string `yaml:"endpoint"`
	X509Certificate string `yaml:"x509_certificate" mapstructure:"x509_certificate"`
	PrivateKey      string `yaml:"private_key" mapstructure:"private_key"`
	CABundle        string `yaml:"ca_bundle" mapstructure:"ca_bundle"`
}

type MerakiConfig struct {
	BaseUrl string `yaml:"base_url" mapstructure:"base_url"`
	ApiKey  string `yaml:"api_key" mapstructure:"api_key"`
	Debug   bool   `yaml:"debug"`
}

type Config struct {
	Push   *PushConfig  `yaml:"push"`
	Meraki MerakiConfig `yaml:"meraki"`
}

func Load() (*Config, error) {
	viper.SetEnvPrefix("S6S_")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.s6s")

	cfgfile := viper.GetString(ConfigFileKey)
	if cfgfile != "" {
		viper.SetConfigFile(cfgfile)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && cfgfile != "" {
			return nil, fmt.Errorf("specified configuration file does not exist: %s", cfgfile)
		} else {
			return nil, fmt.Errorf("failed to read configuration file: %s", err)
		}
	}

	var cfg Config
	// XXX viper.Unmarshal doesn't treat flag bindings as nested
	if err := mapstructure.Decode(viper.AllSettings(), &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %s", err)
	}

	if cfg.Push == nil || cfg.Push.X509Certificate == "" || cfg.Push.PrivateKey == "" || cfg.Push.CABundle == "" {
		bundlefile := viper.GetString(BundleFileKey)
		bundleraw, err := os.ReadFile(bundlefile)
		if err != nil {
			return nil, fmt.Errorf("unable to read credentials bundle file %q: %s", bundlefile, bundleraw)
		}

		if err := yaml.Unmarshal(bundleraw, &cfg.Push); err != nil {
			return nil, fmt.Errorf("failed to decode credentials bundle file %q: %s", bundlefile, err)
		}
	}

	if endpoint := viper.GetString(PushEndpointKey); endpoint != "" {
		cfg.Push.Endpoint = endpoint
	}

	if cfg.Push == nil || cfg.Push.X509Certificate == "" || cfg.Push.PrivateKey == "" {
		return nil, errors.New("datasource Push API credentials are required! check config & bundle files")
	}

	if cfg.Meraki.ApiKey == "" {
		return nil, errors.New("meraki API key is required")
	}

	return &cfg, nil
}

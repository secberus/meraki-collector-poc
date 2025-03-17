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
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigFile = "$HOME/.s6s/config"
	ConfigFileEnvVar  = "S6S_CONFIG_FILE"
	DefaultEndpoint   = "push.secberus.io:7744"
	DefaultBaseUrl    = "https://api.meraki.com/"
)

type S6sConfig struct {
	Endpoint        string `yaml:"endpoint"`
	X509Certificate string `yaml:"x509_certificate"`
	PrivateKey      string `yaml:"private_key"`
	CABundle        string `yaml:"ca_bundle"`
}

type MerakiConfig struct {
	BaseUrl string `yaml:"base_url"`
	ApiKey  string `yaml:"api_key"`
	Debug   bool   `yaml:"debug"`
}

type Config struct {
	S6s    S6sConfig    `yaml:"s6s"`
	Meraki MerakiConfig `yaml:"meraki"`
}

func Load() (*Config, error) {

	var cfgFile string
	cfgFile, ok := os.LookupEnv(ConfigFileEnvVar)
	if !ok {
		cfgFile = DefaultConfigFile
	}

	raw, err := os.ReadFile(os.ExpandEnv(cfgFile))
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	cfg.S6s.Endpoint = DefaultEndpoint
	cfg.Meraki.BaseUrl = DefaultBaseUrl

	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

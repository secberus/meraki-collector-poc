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
	"encoding/json"
	"os"
)

const (
	DefaultConfigFile = "$HOME/.s6s/config"
	ConfigFileEnvVar  = "S6S_CONFIG_FILE"
	DefaultEndpoint   = "push.secberus.io:7744"
	DefaultBaseUrl    = "https://api.meraki.com/"
)

type S6sConfig struct {
	Endpoint        string `json:"endpoint"`
	X509Certificate []byte `json:"x509_certificate"`
	CaKey           []byte `json:"ca_key"`
	CaBundle        []byte `json:"ca_bundle"`
}

type MerakiConfig struct {
	BaseUrl string `json:"base_url"`
	ApiKey  string `json:"api_key"`
	Debug   bool   `json:"debug"`
}

type Config struct {
	S6s    S6sConfig    `json:"s6s"`
	Meraki MerakiConfig `json:"meraki"`
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

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

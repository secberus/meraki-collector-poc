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

package main

import (
	"fmt"
	"strconv"

	meraki "github.com/meraki/dashboard-api-go/v4/sdk"
	"github.com/secberus/meraki-collector/config"
)

func initMerakiClient(cfg *config.MerakiConfig) (*meraki.Client, error) {
	client, err := meraki.NewClientWithOptions(
		cfg.BaseUrl,
		cfg.ApiKey,
		strconv.FormatBool(cfg.Debug),
		"meraki-collector/0.0.0 Secberus (+https://secberus.com)",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Meraki client: %w", err)
	}

	return client, nil
}

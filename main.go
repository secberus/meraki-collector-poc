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
	"context"
	"log"

	"github.com/secberus/meraki-collector/config"
	"github.com/secberus/meraki-collector/resource"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}

	meraki, err := initMerakiClient(&cfg.Meraki)
	if err != nil {
		log.Fatalf("failed to initialize Meraki client: %s", err)
	}

	pushsvc, err := initPushClient(&cfg.S6s)
	if err != nil {
		log.Fatalf("failed to initialize Push client: %s", err)
	}

	collector := NewCollector(meraki, pushsvc)

	// collect from Meraki API root (organizations)
	if err := collector.Collect(context.Background(), resource.Organizations); err != nil {
		log.Fatalf("failed to collect from Meraki API: %s", err)
	}
}

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
	"fmt"
	"log"
	"strings"

	meraki "github.com/meraki/dashboard-api-go/v4/sdk"
	api "github.com/secberus/go-push-api/api/v1"
	service "github.com/secberus/go-push-api/service/v1/push"
	v1 "github.com/secberus/go-push-api/types/v1"

	"github.com/secberus/meraki-collector/resource"
)

type Collector struct {
	tables  map[string]struct{}
	meraki  *meraki.Client
	pushsvc service.PushServiceClient
}

func NewCollector(meraki *meraki.Client, pushsvc service.PushServiceClient) *Collector {
	return &Collector{
		tables:  make(map[string]struct{}),
		meraki:  meraki,
		pushsvc: pushsvc,
	}
}

func (c Collector) register(ctx context.Context, t *v1.Table) error {
	if _, ok := c.tables[t.Name]; ok {
		return nil
	}

	if _, err := c.pushsvc.GetTable(ctx, &api.GetTableInput{TableName: t.Name}); err != nil && strings.Contains(err.Error(), "not found") {
		log.Printf("table %q does not exist, creating\n", t.Name)
		if _, err := c.pushsvc.CreateTable(ctx, &api.CreateTableInput{Table: t}); err != nil {
			return fmt.Errorf("failed to CreateTable: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to GetTable: %w", err)
	}

	c.tables[t.Name] = struct{}{}
	return nil
}

func (c Collector) collect(ctx context.Context, rc *resource.Resource, parent any) error {
	log.Printf("collecting for table %q", rc.Table.Name)

	t := rc.Table
	if err := c.register(ctx, t); err != nil {
		return fmt.Errorf("failed to register table %q: %w", t.Name, err)
	}

	var recs []*v1.Record
	for v, err := range rc.Resolver(ctx, c.meraki, parent) {
		if err != nil {
			return fmt.Errorf("failed to collect for table %q: %w", t.Name, err)
		}
		if r, err := resource.RecordFor(t, v); err != nil {
			return fmt.Errorf("failed to create Record for table %q: %w", t.Name, err)
		} else {
			recs = append(recs, r)
		}
		for _, cr := range rc.Children {
			if err := c.collect(ctx, cr, v); err != nil {
				return fmt.Errorf("failed to collect child for table %q: %w", t.Name, err)
			}
		}
	}

	log.Printf("upserting %d records for table %q", len(recs), t.Name)
	if _, err := c.pushsvc.UpsertRecords(ctx, &api.UpsertRecordsInput{Records: recs}); err != nil {
		return fmt.Errorf("failed to upsert %d records: %w", len(recs), err)
	}

	return nil
}

func (c Collector) Collect(ctx context.Context, rc *resource.Resource) error {
	return c.collect(ctx, rc, nil)
}

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

package resource

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"time"

	meraki "github.com/meraki/dashboard-api-go/v4/sdk"
	v1 "github.com/secberus/go-push-api/types/v1"
)

var ConfigurationChanges = &Resource{
	Table: &v1.Table{
		Name:     "meraki_configuration_changes",
		SyncType: v1.TableSyncType_TABLE_SYNC_TYPE_TRUNCATE,
		Columns:  columnsFor[meraki.ResponseItemOrganizationsGetOrganizationConfigurationChanges]("ts"),
	},
	Resolver: getConfigurationChanges,
}

func getConfigurationChanges(ctx context.Context, client *meraki.Client, org any) iter.Seq2[any, error] {
	orgId := org.(meraki.ResponseItemOrganizationsGetOrganizations).ID
	return func(yield func(any, error) bool) {
		params := &meraki.GetOrganizationConfigurationChangesQueryParams{
			Timespan: 24 * time.Hour.Seconds(),
		}
		rsp, _, err := client.Organizations.GetOrganizationConfigurationChanges(orgId, params)
		if err != nil {
			yield(nil, fmt.Errorf("failed to GetOrganizationConfigurationChanges: %w", err))
			return
		}
		if rsp == nil {
			yield(nil, errors.New("received nil response from GetOrganizationConfigurationChanges"))
			return
		}
		for _, i := range *rsp {
			if !yield(i, nil) {
				return
			}
		}
	}
}

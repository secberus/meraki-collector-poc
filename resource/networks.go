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

	meraki "github.com/meraki/dashboard-api-go/v4/sdk"
	v1 "github.com/secberus/go-push-api/types/v1"
)

var Networks = &Resource{
	Table: &v1.Table{
		Name:     "meraki_networks",
		SyncType: v1.TableSyncType_TABLE_SYNC_TYPE_APPEND,
		Columns:  columnsFor[meraki.ResponseItemOrganizationsGetOrganizationNetworks]("id"),
	},
	Resolver: getOrganizationNetworks,
	Children: []*Resource{
		Devices,
		TopologyLinkLayer,
	},
}

func getOrganizationNetworks(ctx context.Context, client *meraki.Client, org any) iter.Seq2[any, error] {
	orgId := org.(meraki.ResponseItemOrganizationsGetOrganizations).ID
	return func(yield func(any, error) bool) {
		rsp, _, err := client.Organizations.GetOrganizationNetworks(orgId, &meraki.GetOrganizationNetworksQueryParams{PerPage: -1})
		if err != nil {
			yield(nil, fmt.Errorf("failed to GetOrganizationNetworks: %w", err))
			return
		}
		if rsp == nil {
			yield(nil, errors.New("received nil response from GetOrganizationNetworks"))
			return
		}
		for _, i := range *rsp {
			if !yield(i, nil) {
				return
			}
		}
	}
}

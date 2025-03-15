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

var TopologyLinkLayer = &Resource{
	Table: &v1.Table{
		Name:     "meraki_topology_link_layer",
		SyncType: v1.TableSyncType_TABLE_SYNC_TYPE_APPEND,
		Columns:  columnsFor[topologyLinkLayer]("network_id"),
	},
	Resolver: getTopologyLinkLayer,
}

type topologyLinkLayer struct {
	NetworkId string
	Errors    []string
	Links     *[]meraki.ResponseNetworksGetNetworkTopologyLinkLayerLinks
	Nodes     *[]meraki.ResponseNetworksGetNetworkTopologyLinkLayerNodes
}

func getTopologyLinkLayer(ctx context.Context, client *meraki.Client, network any) iter.Seq2[any, error] {
	networkId := network.(meraki.ResponseItemOrganizationsGetOrganizationNetworks).ID
	return func(yield func(any, error) bool) {
		rsp, _, err := client.Networks.GetNetworkTopologyLinkLayer(networkId)
		if err != nil {
			yield(nil, fmt.Errorf("failed to GetNetworkTopologyLinkLayer: %w", err))
			return
		}
		if rsp == nil {
			yield(nil, errors.New("received nil response from GetNetworkTopologyLinkLayer"))
			return
		}
		yield(topologyLinkLayer{
			NetworkId: networkId,
			Errors:    rsp.Errors,
			Links:     rsp.Links,
			Nodes:     rsp.Nodes,
		}, nil)
	}
}

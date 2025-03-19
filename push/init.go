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

package push

import (
	"fmt"

	service "github.com/secberus/go-push-api/service/v1/push"
	"github.com/secberus/meraki-collector/config"
	"google.golang.org/grpc"
)

func Init(cfg *config.PushConfig) (service.PushServiceClient, error) {
	tlsCreds, err := config.Credentials(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load Push credentials: %w", err)
	}

	conn, err := grpc.NewClient(cfg.Endpoint, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Push gRPC client: %w", err)
	}

	return service.NewPushServiceClient(conn), nil
}

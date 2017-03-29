// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package swarm

import (
	"context"
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"github.com/prometheus/common/log"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/util/httputil"
)

func init() {

}

// Discovery retrieves target information from a Swarm master
// and updates them via watches.
type Discovery struct {
	client          *client.Client
	refreshInterval time.Duration
}

// NewDiscovery returns a new Discovery for the given config.
func NewDiscovery(conf *config.SwarmSDConfig) (*Discovery, error) {
	tls, err := httputil.NewTLSConfig(conf.TLSConfig)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{TLSClientConfig: tls}
	wrapper := &http.Client{Transport: transport}

	cli, err := client.NewClient(conf.Server, conf.APIVersion, wrapper, map[string]string{})
	if err != nil {
		return nil, err
	}

	cd := &Discovery{
		client: cli,
	}

	return cd, nil
}

// Run implements the TargetProvider interface.
func (d *Discovery) Run(ctx context.Context, ch chan<- []*config.TargetGroup) {
	for {
		log.Info("Discovering swarm")
		select {
		case <-ctx.Done():
			return
		case <-time.After(d.refreshInterval):
			log.Info("Updated services")
		}
	}
}

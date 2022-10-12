// Copyright © 2022.  Douglas Chimento <dchimento@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/dougEfresh/tasmota-prometheus-service-discovery/pkg"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

type DiscoveryDriver struct {
	tasDiscovery    pkg.Discovery
	address         string
	refreshInterval int
	logger          log.Logger
	oldSourceList   map[string]bool
}

func (d *DiscoveryDriver) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	for c := time.Tick(time.Duration(d.refreshInterval) * time.Second); ; {
		tasmotas, err := d.tasDiscovery.Discover()
		if err != nil {
			level.Error(d.logger).Log("msg", "Error getting tasmotas", "err", err)
			time.Sleep(time.Duration(d.refreshInterval) * time.Second)
			continue
		}
		_ = level.Info(d.logger).Log("msg", fmt.Sprintf("Found %d tasmotas", len(tasmotas)))
		var tgs []*targetgroup.Group
		newSourceList := make(map[string]bool)
		for _, tasmota := range tasmotas {
			targetHost := tasmota.Hostname
			switch tasmota.Hostname {
			case "tasmota":
				newSourceList[tasmota.Address.String()] = true
				targetHost = tasmota.Address.String()
				break
			default:
				newSourceList[tasmota.Hostname] = true
				break
			}
			target := model.LabelSet{
				model.AddressLabel: model.LabelValue(targetHost),
			}
			tgs = append(tgs, &targetgroup.Group{
				Source:  targetHost,
				Targets: []model.LabelSet{target},
			})
		}
		// When targetGroup disappear, send an update with empty targetList.
		for key := range d.oldSourceList {
			if !newSourceList[key] {
				tgs = append(tgs, &targetgroup.Group{
					Source: key,
				})
			}
		}
		d.oldSourceList = newSourceList
		ch <- tgs
		// Wait for ticker or exit when ctx is closed.
		select {
		case <-c:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func New(ctx context.Context, tasDiscovery pkg.Discovery, logger log.Logger) discovery.Discoverer {
	return &DiscoveryDriver{
		tasDiscovery:    tasDiscovery,
		logger:          logger,
		oldSourceList:   map[string]bool{},
		refreshInterval: 300,
	}
}

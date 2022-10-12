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

// Package prometheus https://github.com/prometheus/prometheus/blob/main/documentation/examples/custom-sd/adapter/adapter.go
package prometheus

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"

	//	"reflect"

	//"github.com/dougEfresh/tasmota-prometheus-service-discovery/pkg/domain"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/model"

	//	"github.com/go-kit/log/level"

	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

type Adapter struct {
	ctx     context.Context
	disc    discovery.Discoverer
	manager *discovery.Manager
	output  string
	targets []string
	name    string
	logger  log.Logger
}

func (a *Adapter) refreshTargetGroups(allTargetGroups map[string][]*targetgroup.Group) {
	var allTargets []string
	for _, groups := range allTargetGroups {
		for _, group := range groups {
			for _, target := range group.Targets {
				allTargets = append(allTargets, string(target[model.AddressLabel]))
			}
		}
	}
	sort.Slice(allTargets, func(i, j int) bool {
		return allTargets[i] > allTargets[j]
	})
	if !reflect.DeepEqual(a.targets, allTargets) {
		a.targets = allTargets
		level.Info(log.With(a.logger, "component", "sd-adapter")).Log("updated targets")
		err := a.writeOutput()
		if err != nil {
			level.Error(log.With(a.logger, "component", "sd-adapter")).Log("err", err)
		}
	}
}

func (a *Adapter) runCustomSD(ctx context.Context) {
	updates := a.manager.SyncCh()
	for {
		select {
		case <-ctx.Done():
		case allTargetGroups, ok := <-updates:
			// Handle the case that a target provider exits and closes the channel
			// before the context is done.
			if !ok {
				return
			}
			a.refreshTargetGroups(allTargetGroups)
		}
	}
}

func (a *Adapter) writeOutput() error {
	type customSD struct {
		Targets []string `json:"targets"`
		Labels  map[string]string
	}
	tasmotas := customSD{
		Targets: a.targets,
		Labels:  map[string]string{"job": "tasmota"},
	}
	out := []customSD{tasmotas}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	dir, _ := filepath.Split(a.output)
	tmpfile, err := os.CreateTemp(dir, "sd-adapter")
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	_, err = tmpfile.Write(b)
	if err != nil {
		return err
	}

	// Close the file immediately for platforms (eg. Windows) that cannot move
	// a file while a process is holding a file handle.
	tmpfile.Close()
	err = os.Rename(tmpfile.Name(), a.output)
	if err != nil {
		return err
	}
	os.Chmod(a.output, 0444)
	return nil
}

// Run starts a DiscoveryDriver Manager and the custom service DiscoveryDriver implementation.
func (a *Adapter) Run() {
	//nolint:errcheck
	go a.manager.Run()
	a.manager.StartCustomProvider(a.ctx, a.name, a.disc)
	go a.runCustomSD(a.ctx)
}

func NewAdapter(ctx context.Context, file, name string, d discovery.Discoverer, logger log.Logger) *Adapter {
	return &Adapter{
		ctx:  ctx,
		disc: d,
		//groups:  make(map[string]*domain.TasmotaNet),
		manager: discovery.NewManager(ctx, logger),
		output:  file,
		name:    name,
		logger:  logger,
	}
}

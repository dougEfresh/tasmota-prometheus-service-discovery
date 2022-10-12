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

package cmd

import (
	"os"
	"time"

	"github.com/dougEfresh/tasmota-prometheus-service-discovery/pkg"
	"github.com/dougEfresh/tasmota-prometheus-service-discovery/pkg/prometheus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/spf13/cobra"
)

func newMqttCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mqtt",
		Short: "discover tasmotas from mqtt",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.NewSyncLogger(log.NewLogfmtLogger(os.Stdout))
			opts := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER"))
			opts.SetClientID("prometheus_service_discovery").SetKeepAlive(time.Minute)
			opts.SetUsername(os.Getenv("MQTT_USER")).SetPassword(os.Getenv("MQTT_PASSWORD"))
			client := mqtt.NewClient(opts)
			mqttDiscover := pkg.NewMqttDiscover(client, 5*time.Second)
			disc := prometheus.New(cmd.Context(), mqttDiscover, logger)
			sdAdapter := prometheus.NewAdapter(cmd.Context(), "/etc/prometheus/static/tasmotas.json", "tasmota-discovery", disc, logger)
			sdAdapter.Run()
			<-cmd.Context().Done()
		},
	}
}

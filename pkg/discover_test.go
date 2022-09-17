// Copyright © 2022.  Douglas Chimento <dchimento@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// nolint gochecknoglobals
package pkg_test

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"testing"
	"time"

	"github.com/dougEfresh/tasmota-prometheus-service-discovery/pkg"
	"github.com/dougEfresh/tasmota-prometheus-service-discovery/pkg/domain"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqttserver "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/stretchr/testify/require"
)

var responses = []domain.StatusNet{
	{
		StatusNet: domain.TasmotaNet{
			Hostname: "blah",
			Address: domain.Address{
				net.ParseIP("192.168.4.1"),
			},
		},
	},
	{
		StatusNet: domain.TasmotaNet{
			Hostname: "blah2",
			Address: domain.Address{
				net.ParseIP("192.168.4.2"),
			},
		},
	},
}

func TestNewMqttDiscover(t *testing.T) {
	server := mqttserver.NewServer(nil)
	tcp := listeners.NewTCP("t1", "localhost:1883")
	err := server.AddListener(tcp, nil)
	require.NoError(t, err)

	defer server.Close()

	go func() {
		_ = server.Serve()
	}()

	time.Sleep(1 * time.Second)
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("prometheus_service_discovery").SetKeepAlive(time.Minute)

	client := mqtt.NewClient(opts)
	testClient := mqtt.NewClient(opts.SetClientID("test_client"))
	token := testClient.Connect()
	token.Wait()
	require.NoError(t, token.Error())

	// nolint varnamelen
	token = testClient.Subscribe("#", 0, func(client mqtt.Client, message mqtt.Message) {
		t.Logf("Recv:  topic:%s message:%s", message.Topic(), string(message.Payload()))
		if message.Topic() == "cmnd/tasmotas/status" && string(message.Payload()) == "5" {
			for _, r := range responses {

				b, e := json.MarshalIndent(r, "", "  ")
				if e != nil {
					t.Logf("Failed to marshal responses %s", e.Error())

					return
				}
				client.Publish(fmt.Sprintf("stat/%s/status5", r.StatusNet.Hostname), 0, false, b)
			}
		}
	})
	token.Wait()
	require.NoError(t, token.Error())

	d := pkg.NewMqttDiscover(client)
	tasmotas, err := d.Discover()
	require.NoError(t, err)
	require.True(t, len(tasmotas) > 0)

	sort.Slice(tasmotas, func(i, j int) bool {
		return tasmotas[i].Hostname < tasmotas[j].Hostname
	})

	require.Equal(t, len(responses), len(tasmotas))
	for i, tasmota := range tasmotas {
		require.Equal(t, responses[i].StatusNet, tasmota)
	}
}

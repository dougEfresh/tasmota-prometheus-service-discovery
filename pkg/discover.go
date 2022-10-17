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

package pkg

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/dougEfresh/tasmota-prometheus-service-discovery/pkg/domain"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-multierror"
)

type Discovery interface {
	Discover() ([]domain.TasmotaNet, error)
}

type MqttDiscover struct {
	client  mqtt.Client
	timeout time.Duration
	sync.Mutex
}

const DISCONNECT = 1000

// nolint wrapcheck
func (m *MqttDiscover) Discover() ([]domain.TasmotaNet, error) {
	var tasmotas []domain.TasmotaNet

	var errors error

	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	defer m.client.Disconnect(DISCONNECT)

	var hdl mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
		m.Lock()
		defer m.Unlock()
		var s domain.StatusNet
		if err := json.Unmarshal(message.Payload(), &s); err != nil {
			errors = multierror.Append(errors, err)
			return
		}
		if s.StatusNet.Hostname != "" {
			tasmotas = append(tasmotas, s.StatusNet)
		}
	}

	if token := m.client.Subscribe("stat/#", 0, hdl); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	defer func() {
		token := m.client.Unsubscribe("stat/#")
		token.Wait()
	}()

	if token := m.client.Publish("cmnd/tasmotas/status", 0, false, "5"); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	time.Sleep(m.timeout)
	m.client.Disconnect(DISCONNECT)

	return tasmotas, errors

}

func NewMqttDiscover(client mqtt.Client, timeout time.Duration) *MqttDiscover {
	return &MqttDiscover{
		client:  client,
		timeout: timeout,
	}
}

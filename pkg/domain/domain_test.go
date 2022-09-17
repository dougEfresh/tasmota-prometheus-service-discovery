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

package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/dougEfresh/tasmota-prometheus-service-discovery/pkg/domain"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	jsonInput := `{
  "StatusNET": {
    "Hostname": "some-light",
    "IPAddress": "192.168.4.70",
    "Gateway": "192.168.4.4",
    "Subnetmask": "255.255.255.0",
    "DNSServer1": "192.168.4.2",
    "Mac": "C4:4F:33:D3:FE:41",
    "Webserver": 2,
    "HTTP_API": 1,
    "WifiConfig": 4,
    "WifiPower": 1
  }
}
`

	var result domain.StatusNet
	err := json.Unmarshal([]byte(jsonInput), &result)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.StatusNet)
	require.NotNil(t, result.StatusNet.Address)
	require.Equal(t, result.StatusNet.Hostname, "some-light")
	require.Equal(t, result.StatusNet.Address.String(), "192.168.4.70")
}

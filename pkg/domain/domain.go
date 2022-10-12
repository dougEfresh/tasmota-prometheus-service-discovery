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

package domain

import (
	"encoding/json"
	"fmt"
	"net"
)

type Address struct {
	net.IP
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var address string
	if err := json.Unmarshal(data, &address); err != nil {
		// nolint wrapcheck
		return err
	}

	a.IP = net.ParseIP(address)
	if a.IP == nil {
		// nolint goerr113
		return fmt.Errorf("unable to parse %s", address)
	}

	return nil
}

// nolint tagliatelle
type StatusNet struct {
	StatusNet TasmotaNet `json:"StatusNET"`
}

type TasmotaNet struct {
	Hostname string
	Address  Address `json:"IPAddress"`
}

// Copyright [2023] [Argus]
//
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

package deployedver

import (
	"github.com/release-argus/Argus/util"
)

// GetAllowInvalidCerts returns whether invalid HTTPS certs are allowed.
func (l *Lookup) GetAllowInvalidCerts() bool {
	return *util.GetFirstNonNilPtr(
		l.AllowInvalidCerts,
		l.Defaults.AllowInvalidCerts,
		l.HardDefaults.AllowInvalidCerts)
}

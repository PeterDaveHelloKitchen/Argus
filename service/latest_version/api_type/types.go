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

package types

import (
	"encoding/json"

	"github.com/coreos/go-semver/semver"
)

// Release is the format of a Release on api.github.com/repos/OWNER/REPO/releases.
type Release struct {
	URL             string          `json:"url,omitempty"`
	AssetsURL       string          `json:"assets_url,omitempty"`
	SemanticVersion *semver.Version `json:"-"`
	TagName         string          `json:"tag_name,omitempty"`
	PreRelease      bool            `json:"prerelease,omitempty"`
	Assets          []Asset         `json:"assets,omitempty"`
}

// String returns a string representation of the Release.
func (r *Release) String() string {
	jsonBytes, _ := json.Marshal(r)
	return string(jsonBytes)
}

// Asset is the format of an Asset on api.github.com/repos/OWNER/REPO/releases.
type Asset struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name,omitempty"`
	URL                string `json:"url,omitempty"`
	BrowserDownloadURL string `json:"browser_download_url,omitempty"`
}

// String returns a string representation of the Asset.
func (a *Asset) String() string {
	jsonBytes, _ := json.Marshal(a)
	return string(jsonBytes)
}

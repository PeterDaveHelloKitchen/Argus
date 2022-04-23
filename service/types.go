// Copyright [2022] [Hymenaios]
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

package service

import (
	"github.com/coreos/go-semver/semver"

	"github.com/hymenaios-io/Hymenaios/notifiers/gotify"
	"github.com/hymenaios-io/Hymenaios/notifiers/slack"
	"github.com/hymenaios-io/Hymenaios/webhook"
)

// Slice is a slice mapping of Service.
type Slice map[string]*Service

// Service is a source to be serviceed and provides everything needed to extract
// the latest version from the URL provided.
type Service struct {
	ID                 *string          `yaml:"-"`                             // service_name.
	Type               *string          `yaml:"type,omitempty"`                // "github"/"URL"
	URL                *string          `yaml:"url,omitempty"`                 // type:URL - "https://example.com", type:github - "owner/repo" or "https://github.com/owner/repo".
	AllowInvalidCerts  *bool            `yaml:"allow_invalid_certs,omitempty"` // default - false = Disallows invalid HTTPS certificates.
	AccessToken        *string          `yaml:"access_token,omitempty"`        // GitHub access token to use.
	SemanticVersioning *bool            `yaml:"semantic_versioning,omitempty"` // default - true  = Version has to follow semantic versioning (https://semver.org/) and be greater than the previous to trigger anything.
	Interval           *string          `yaml:"interval,omitempty"`            // AhBmCs = Sleep A hours, B minutes and C seconds between queries.
	URLCommands        *URLCommandSlice `yaml:"url_commands,omitempty"`        // Commands to filter the release from the URL request.
	RegexContent       *string          `yaml:"regex_content,omitempty"`       // "abc-[a-z]+-{{ version }}_amd64.deb" This regex must exist in the body of the URL to trigger new version actions.
	RegexVersion       *string          `yaml:"regex_version,omitempty"`       // "v*[0-9.]+" The version found must match this release to trigger new version actions.
	UsePreRelease      *bool            `yaml:"use_prerelease,omitempty"`      // Whether the prerelease tag should be used (prereleases are ignored by default).
	WebURL             *string          `yaml:"web_url,omitempty"`             // URL to provide on the Web UI.
	AutoApprove        *bool            `yaml:"auto_approve,omitempty"`        // default - true = Requre approval before sending WebHook(s) for new releases.
	IgnoreMisses       *bool            `yaml:"ignore_misses,omitempty"`       // Ignore URLCommands that fail (e.g. split on text that doesn't exist).
	Icon               *string          `yaml:"icon,omitempty"`                // Icon URL to use for Slack messages/Web UI.
	Gotify             *gotify.Slice    `yaml:"gotify,omitempty"`              // Service-specific Gotify vars.
	Slack              *slack.Slice     `yaml:"slack,omitempty"`               // Service-specific Slack vars.
	WebHook            *webhook.Slice   `yaml:"webhook,omitempty"`             // Service-specific WebHook vars.
	Status             *Status          `yaml:"status,omitempty"`              // Track the Status of this source (version and regex misses).
	HardDefaults       *Service         `yaml:"-"`                             // Hardcoded default values.
	Defaults           *Service         `yaml:"-"`                             // Default values.
	Announce           *chan []byte     `yaml:"-"`                             // Announce to the WebSocket.
	SaveChannel        *chan bool       `yaml:"-"`                             // Channel for triggering a save of the config.
}

// GitHubRelease is the format of a Release on api.github.com/repos/OWNER/REPO/releases.
type GitHubRelease struct {
	URL             string          `json:"url"`
	AssetsURL       string          `json:"assets_url"`
	UploadURL       string          `json:"upload_url"`
	HTMLURL         string          `json:"html_url"`
	ID              uint            `json:"id"`
	Author          GitHubAuthor    `json:"author"`
	NodeID          string          `json:"node_id"`
	SemanticVersion *semver.Version `json:"-"`
	TagName         string          `json:"tag_name"`
	TargetCommitish string          `json:"target_commitish"`
	Name            string          `json:"name"`
	Draft           bool            `json:"draft"`
	PreRelease      bool            `json:"prerelease"`
	CreatedAt       string          `json:"created_at"`
	PublishedAt     string          `json:"published_at"`
	Assets          []GitHubAsset   `json:"assets"`
}

// GitHubAuthor is the format of an Author on api.github.com/repos/OWNER/REPO/releases.
type GitHubAuthor struct {
	Login             string `json:"login"`
	ID                uint   `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received__events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

// GitHubAsset is the format of an Asset on api.github.com/repos/OWNER/REPO/releases.
type GitHubAsset struct {
	URL                string       `json:"url"`
	ID                 uint         `json:"id"`
	NodeID             string       `json:"node_id"`
	Name               string       `json:"name"`
	Label              string       `json:"label"`
	Uploader           GitHubAuthor `json:"uploader"`
	ContentType        string       `json:"content_type"`
	State              string       `json:"state"`
	Size               uint         `json:"size"`
	DownloadCount      uint         `json:"download_count"`
	CreatedAt          string       `json:"created_at"`
	UpdatedAt          string       `json:"updated_at"`
	BrowserDownloadURL string       `json:"browser_download_url"`
}

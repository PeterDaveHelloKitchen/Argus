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

//go:build unit

package testing

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/release-argus/Argus/config"
	"github.com/release-argus/Argus/notifiers/shoutrrr"
	"github.com/release-argus/Argus/service"
	"github.com/release-argus/Argus/util"
)

func TestGetAllShoutrrrNames(t *testing.T) {
	// GIVEN various Service's and Notify's
	testLogging()
	tests := map[string]struct {
		service       service.Slice
		rootNotifiers shoutrrr.Slice
		want          []string
	}{
		"nothing": {},
		"only service notifiers": {
			service: service.Slice{
				"0": {Notify: shoutrrr.Slice{"foo": {}}},
				"1": {Notify: shoutrrr.Slice{"bar": {}}},
			},
			want: []string{"bar", "foo"},
		},
		"only service notifiers with duplicates": {
			service: service.Slice{
				"0": {Notify: shoutrrr.Slice{"foo": {}}},
				"1": {Notify: shoutrrr.Slice{"foo": {}, "bar": {}}},
			},
			want: []string{"bar", "foo"},
		},
		"only root notifiers": {rootNotifiers: shoutrrr.Slice{"foo": {}, "bar": {}},
			want: []string{"bar", "foo"},
		},
		"root + service notifiers": {
			service: service.Slice{
				"0": {Notify: shoutrrr.Slice{"foo": {}}},
				"1": {Notify: shoutrrr.Slice{"foo": {}, "bar": {}}},
			},
			rootNotifiers: shoutrrr.Slice{"baz": {}},
			want:          []string{"bar", "baz", "foo"},
		},
		"root + service notifiers with duplicates": {
			service: service.Slice{
				"0": {Notify: shoutrrr.Slice{"foo": {}}},
				"1": {Notify: shoutrrr.Slice{"foo": {}, "bar": {}}},
			},
			rootNotifiers: shoutrrr.Slice{"foo": {}, "bar": {}, "baz": {}},
			want:          []string{"bar", "baz", "foo"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			cfg := config.Config{
				Service: tc.service,
				Notify:  tc.rootNotifiers,
			}

			// WHEN getAllShoutrrrNames is called on this config
			got := getAllShoutrrrNames(&cfg)

			// THEN a list of all shoutrrr's will be returned
			if len(got) != len(tc.want) {
				t.Fatalf("lists differ in length\nwant: %s\ngot:  %s",
					tc.want, got)
			}
			gotIndex := 0
			for gotIndex != 0 {
				found := false
				for wantIndex := range tc.want {
					if got[gotIndex] == tc.want[wantIndex] {
						found = true
						util.RemoveIndex(&got, gotIndex)
						util.RemoveIndex(&tc.want, wantIndex)
						break
					}
				}
				if !found {
					t.Fatalf("want: %v\ngot: %v",
						tc.want, got)
				}
				gotIndex--
			}
		})
	}
}

func TestFindShoutrrr(t *testing.T) {
	// GIVEN a Config with/without Service containing a Shoutrrr and Root Shoutrrr(s)
	jLog = util.NewJLog("INFO", false)
	testLogging()
	tests := map[string]struct {
		flag        string
		cfg         *config.Config
		outputRegex *string
		panicRegex  *string
		foundInRoot *bool
	}{
		"empty search with only Service notifiers": {
			flag:       "",
			panicRegex: stringPtr(`could not be found.*\s+.*one of these?.*\s+.* bar\s+.* baz\s+.* foo\s+`),
			cfg: &config.Config{
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{
							"foo": {},
							"bar": {},
							"baz": {}}}}},
		},
		"empty search with only Root notifiers": {
			flag:       "",
			panicRegex: stringPtr(`could not be found.*\s+.*one of these?.*\s+.* bar\s+.* baz\s+.* foo\s+`),
			cfg: &config.Config{
				Notify: shoutrrr.Slice{
					"foo": {},
					"bar": {},
					"baz": {}}},
		},
		"empty search with Root notifiers and Service notifiers and no duplicates": {
			flag:       "",
			panicRegex: stringPtr(`could not be found.*\s+.*one of these?.*\s+.* bar\s+.* baz\s+.* foo\s+`),
			cfg: &config.Config{
				Notify: shoutrrr.Slice{
					"foo": {},
					"bar": {},
					"baz": {}},
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{
							"foo": {},
							"bar": {},
							"baz": {}}}}},
		},
		"empty search with Root notifiers and Service notifiers and duplicates": {
			flag:       "",
			panicRegex: stringPtr(`could not be found.*\s+.*one of these?.*\s+.* bar\s+.* baz\s+.* foo\s+.* shazam\s+`),
			cfg: &config.Config{
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{"foo": {}, "bar": {}, "baz": {}}}},
				Notify: shoutrrr.Slice{
					"foo":    {},
					"shazam": {},
					"baz":    {}}},
		},
		"matching search of notifier in Root": {
			flag:        "bosh",
			outputRegex: stringPtr("^$"),
			cfg: &config.Config{
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{"foo": {}, "bar": {}, "baz": {}}}},
				Notify: shoutrrr.Slice{"bosh": {
					Type: "gotify",
					URLFields: map[string]string{
						"host": "example.com", "token": "example"}}}},
			foundInRoot: boolPtr(true),
		},
		"matching search of notifier in Service": {
			flag:        "baz",
			outputRegex: stringPtr("^$"),
			cfg: &config.Config{
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{"foo": {}, "bar": {}, "baz": {
							Type: "gotify",
							URLFields: map[string]string{
								"host": "example.com", "token": "example"}}}}}},
			foundInRoot: boolPtr(false),
		},
		"matching search of notifier in Root and a Service": {
			flag:        "bar",
			outputRegex: stringPtr("^$"),
			cfg: &config.Config{
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{
							"foo": {},
							"bar": {
								Type: "gotify",
								URLFields: map[string]string{
									"host": "example.com", "token": "example"}},
							"baz": {}}}},
				Notify: shoutrrr.Slice{"bar": {
					Type: "gotify",
					URLFields: map[string]string{
						"host": "example.com", "token": "example"},
				}}},
			foundInRoot: boolPtr(true),
		},
		"matching search of Service notifier with invalid config fixed with Defaults": {
			flag: "bar",
			cfg: &config.Config{
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{
							"foo": {},
							"bar": {
								Type: "smtp",
								URLFields: map[string]string{
									"host": "example.com"},
								Params: map[string]string{
									"fromaddress": "test@release-argus.io"}},
							"baz": {}}}},
				Defaults: config.Defaults{
					Notify: shoutrrr.Slice{
						"something": {
							Type: "something",
							Params: map[string]string{
								"title": "bar"}},
						"smtp": {
							Params: map[string]string{
								"toaddresses": "me@you.com"}}}}},
			foundInRoot: boolPtr(false),
		},
		"matching search of Service notifier with invalid config fixed with Root": {
			flag: "bar",
			cfg: &config.Config{
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{
							"foo": {},
							"bar": {
								Type: "smtp",
								URLFields: map[string]string{
									"host": "example.com"},
								Params: map[string]string{
									"fromaddress": "test@release-argus.io"}},
							"baz": {}}}},
				Notify: shoutrrr.Slice{
					"something": {
						Type: "something",
						Params: map[string]string{
							"title": "bar"}},
					"smtp": {
						Params: map[string]string{
							"toaddresses": "me@you.com"}}}},
			foundInRoot: boolPtr(false),
		},
		"matching search of Service notifier with invalid config fixed with Root and Defaults": {
			flag: "bosh",
			cfg: &config.Config{
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{
							"foo": {},
							"bosh": {
								Type: "smtp",
								URLFields: map[string]string{
									"host": "example.com"},
								Params: map[string]string{
									"fromaddress": "test@release-argus.io"}},
							"baz": {}}}},
				Notify: shoutrrr.Slice{
					"bosh": {
						Type: "smtp",
						URLFields: map[string]string{
							"host": "example.com"}}},
				Defaults: config.Defaults{
					Notify: shoutrrr.Slice{
						"something": {
							Type: "something",
							Params: map[string]string{
								"title": "bar"}},
						"smtp": {
							Params: map[string]string{
								"toaddresses": "me@you.com"}}}}},
			foundInRoot: boolPtr(false),
		},
		"matching search of Root notifier with invalid config": {
			flag:       "bosh",
			panicRegex: stringPtr(`bosh:\s+params:\s+toaddresses: <required>`),
			cfg: &config.Config{
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{
							"foo": {},
							"bar": {},
							"baz": {}}}},
				Notify: shoutrrr.Slice{
					"bosh": {
						Type: "smtp",
						URLFields: map[string]string{
							"host": "example.com"},
						Params: map[string]string{
							"fromaddress": "test@release-argus.io"}}}},
			foundInRoot: boolPtr(true),
		},
		"matching search of Root notifier with invalid config fixed with Defaults": {
			flag: "bosh",
			cfg: &config.Config{
				Service: service.Slice{
					"argus": {
						Notify: shoutrrr.Slice{
							"foo": {},
							"bar": {},
							"baz": {}}}},
				Notify: shoutrrr.Slice{
					"bosh": {
						Type: "smtp",
						URLFields: map[string]string{
							"host": "example.com"},
						Params: map[string]string{
							"fromaddress": "test@release-argus.io"}}},
				Defaults: config.Defaults{
					Notify: shoutrrr.Slice{
						"something": {
							Type: "something",
							Params: map[string]string{
								"title": "bar"}},
						"smtp": {
							Params: map[string]string{
								"toaddresses": "me@you.com"}}}}},
			foundInRoot: boolPtr(true),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			stdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			jLog.Testing = true
			if tc.panicRegex != nil {
				// Switch Fatal to panic and disable this panic.
				defer func() {
					r := recover()
					rStr := fmt.Sprint(r)
					re := regexp.MustCompile(*tc.panicRegex)
					match := re.MatchString(rStr)
					if !match {
						t.Errorf("expected a panic that matched %q\ngot: %q",
							*tc.panicRegex, rStr)
					}
				}()
			}

			// WHEN findShoutrrr is called with the test Config
			got := findShoutrrr(tc.flag, tc.cfg, jLog, &util.LogFrom{})

			// THEN we get the expected output
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = stdout
			output := string(out)
			if tc.outputRegex != nil {
				re := regexp.MustCompile(*tc.outputRegex)
				match := re.MatchString(output)
				if !match {
					t.Fatalf("want match for %q\nnot: %q",
						*tc.outputRegex, output)
				}
			}
			// if the notifier should have been found in the root or in a service
			if tc.foundInRoot != nil {
				if *tc.foundInRoot {
					if !identicalNotifiers(tc.cfg.Notify[tc.flag], got["test"]) {
						t.Fatalf("want: %v\ngot: %v",
							tc.cfg.Notify[tc.flag], got["test"])
					}
				} else {
					if !identicalNotifiers(tc.cfg.Service["argus"].Notify[tc.flag], got["test"]) {
						t.Fatalf("want: %v\ngot: %v",
							tc.cfg.Service["argus"].Notify[tc.flag], got["test"])
					}
					// would have been given in the Init
					got["test"].Defaults = tc.cfg.Defaults.Notify[got["test"].Type]
				}
			}
			// if there were defaults for that type
			if tc.cfg.Defaults.Notify[got["test"].Type] != nil {
				if !identicalNotifiers(tc.cfg.Defaults.Notify[got["test"].Type], got["test"].Defaults) {
					t.Fatalf("defaults were not applied\nwant: %v\ngot: %v",
						tc.cfg.Defaults.Notify[got["test"].Type], got["test"].Defaults)
				}
			}
		})
	}
}

func identicalNotifiers(a *shoutrrr.Shoutrrr, b *shoutrrr.Shoutrrr) (identical bool) {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	identical = a.Type == b.Type && a.ID == b.ID && len(a.Options) == len(b.Options) && len(a.URLFields) == len(b.URLFields) && len(a.Params) == len(b.Params)
	for i := range a.Options {
		if a.Options[i] != b.Options[i] {
			identical = false
		}
	}
	for i := range a.URLFields {
		if a.URLFields[i] != b.URLFields[i] {
			identical = false
		}
	}
	for i := range a.Params {
		if a.Params[i] != b.Params[i] {
			identical = false
		}
	}
	return
}

func TestNotifyTest(t *testing.T) {
	// GIVEN a Config with/without Service containing a Shoutrrr and Root Shoutrrr(s)
	jLog = util.NewJLog("INFO", false)
	testLogging()
	emptyShoutrrr := shoutrrr.Shoutrrr{
		Options:   map[string]string{},
		URLFields: map[string]string{},
		Params:    map[string]string{},
	}
	tests := map[string]struct {
		flag        string
		slice       service.Slice
		outputRegex *string
		panicRegex  *string
	}{
		"empty flag": {flag: "",
			outputRegex: stringPtr("^$"),
			slice: service.Slice{
				"argus": {
					Notify: shoutrrr.Slice{
						"foo": {},
						"bar": {},
						"baz": {},
					},
				}}},
		"unknown Notifier": {flag: "something",
			panicRegex: stringPtr("Notifier.* could not be found"),
			slice: service.Slice{
				"argus": {
					Notify: shoutrrr.Slice{
						"foo": {},
						"bar": {},
						"baz": {},
					},
				}}},
		"known Service Notifier with invalid Gotify token": {flag: "bar",
			panicRegex: stringPtr(`Message failed to send with "bar" config\s+invalid gotify token`),
			slice: service.Slice{
				"argus": {
					Notify: shoutrrr.Slice{
						"foo": {},
						"bar": {
							ID:   "bar",
							Type: "gotify",
							Options: map[string]string{
								"max_tries": "1",
							},
							URLFields: map[string]string{
								"host": "example.com", "token": "invalid"},
							Params:       map[string]string{},
							Main:         &emptyShoutrrr,
							Defaults:     &emptyShoutrrr,
							HardDefaults: &emptyShoutrrr,
						},
						"baz": {},
					},
				}}},
		"valid Gotify token": {flag: "bar",
			panicRegex: stringPtr(`HTTP 404 Not Found`),
			slice: service.Slice{
				"argus": {
					Notify: shoutrrr.Slice{
						"foo": {},
						"bar": {
							ID:   "bar",
							Type: "gotify",
							Options: map[string]string{
								"max_tries": "1",
							},
							URLFields: map[string]string{
								"host": "example.com", "token": "AGdjFCZugzJGhEG"},
							Params:       map[string]string{},
							Main:         &emptyShoutrrr,
							Defaults:     &emptyShoutrrr,
							HardDefaults: &emptyShoutrrr,
						},
						"baz": {},
					},
				}}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			stdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			jLog.Testing = true
			for i := range tc.slice {
				(*tc.slice[i]).Status.Init(
					len((*tc.slice[i]).Notify), 0, 0,
					&name,
					nil)
				for j := range (*tc.slice[i]).Notify {
					(*tc.slice[i]).Notify[j].Init(
						&(*tc.slice[i]).Status,
						(*tc.slice[i]).Notify[j].Main,
						(*tc.slice[i]).Notify[j].Defaults,
						(*tc.slice[i]).Notify[j].HardDefaults,
					)
				}
			}
			if tc.panicRegex != nil {
				// Switch Fatal to panic and disable this panic.
				defer func() {
					r := recover()
					rStr := fmt.Sprint(r)
					re := regexp.MustCompile(*tc.panicRegex)
					match := re.MatchString(rStr)
					if !match {
						t.Errorf("expected a panic that matched %q\ngot: %q",
							*tc.panicRegex, rStr)
					}
				}()
			}

			// WHEN NotifyTest is called with the test Config
			cfg := config.Config{
				Service: tc.slice,
			}
			NotifyTest(&tc.flag, &cfg, jLog)

			// THEN we get the expected output
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = stdout
			output := string(out)
			if tc.outputRegex != nil {
				re := regexp.MustCompile(*tc.outputRegex)
				match := re.MatchString(output)
				if !match {
					t.Errorf("want match for %q\nnot: %q",
						*tc.outputRegex, output)
				}
			}
		})
	}
}

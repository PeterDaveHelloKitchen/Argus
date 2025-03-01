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

package latestver

import (
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/release-argus/Argus/service/latest_version/filter"
	opt "github.com/release-argus/Argus/service/options"
	"github.com/release-argus/Argus/util"
)

func TestLookup_Print(t *testing.T) {
	// GIVEN a Lookup
	tests := map[string]struct {
		lookup      *Lookup
		urlCommands filter.URLCommandSlice
		require     *filter.Require
		options     opt.Options
		lines       int
	}{
		"minimal github type with no urlCommands/require": {
			lookup: &Lookup{
				Type: "github",
				URL:  "release-argus/Argus"},
			lines: 3,
		},
		"fully defined github type with no urlCommands/require": {
			lookup: testLookup(false, false),
			lines:  6,
		},
		"url type with no urlCommands/require": {
			lookup: testLookup(true, false),
			lines:  4,
		},
		"url type with urlCommands and no require": {
			lookup: testLookup(true, false),
			lines:  7,
			urlCommands: filter.URLCommandSlice{
				{Type: "regex", Regex: stringPtr("foo")}},
		},
		"github type with urlCommands and no require": {
			lookup: testLookup(false, false),
			lines:  9,
			urlCommands: filter.URLCommandSlice{
				{Type: "regex", Regex: stringPtr("foo")}},
		},
		"url type with require and no urlCommands": {
			lookup:  testLookup(true, false),
			lines:   6,
			require: &filter.Require{RegexContent: "foo"},
		},
		"github type with require and no urlCommands": {
			lookup:  testLookup(false, false),
			lines:   8,
			require: &filter.Require{RegexContent: "foo"},
		},
		"url type with urlCommands and require": {
			lookup: testLookup(true, false),
			lines:  9,
			urlCommands: filter.URLCommandSlice{
				{Type: "regex", Regex: stringPtr("foo")}},
			require: &filter.Require{
				RegexContent: "foo"},
			options: opt.Options{Active: boolPtr(false)},
		},
		"github type with urlCommands and require": {
			lookup: testLookup(false, false),
			lines:  11,
			urlCommands: filter.URLCommandSlice{
				{Type: "regex", Regex: stringPtr("foo")}},
			require: &filter.Require{RegexContent: "foo"},
			options: opt.Options{Active: boolPtr(false)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			stdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			tc.lookup.Require = tc.require
			tc.lookup.URLCommands = tc.urlCommands
			tc.lookup.Options = &tc.options

			// WHEN Print is called
			tc.lookup.Print("")

			// THEN it prints the expected number of lines
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = stdout
			got := strings.Count(string(out), "\n")
			if got != tc.lines {
				t.Errorf("Print should have given %d lines, but gave %d\n%s",
					tc.lines, got, out)
			}
		})
	}
}

func TestLookup_CheckValues(t *testing.T) {
	// GIVEN a Lookup
	tests := map[string]struct {
		lType       *string
		url         *string
		wantURL     *string
		require     *filter.Require
		urlCommands *filter.URLCommandSlice
		errRegex    string
	}{
		"valid service": {
			errRegex: `^$`,
		},
		"no type": {
			errRegex: `type: <required>`,
			lType:    stringPtr(""),
		},
		"invalid type": {
			errRegex: `type: .* <invalid>`,
			lType:    stringPtr("foo"),
		},
		"no url": {
			errRegex: `url: <required>`,
			url:      stringPtr(""),
		},
		"corrects github url": {
			errRegex: `^$`,
			url:      stringPtr("https://github.com/release-argus/Argus"),
			wantURL:  stringPtr("release-argus/Argus"),
		},
		"invalid require": {
			errRegex: `regex_content: .* <invalid>`,
			require:  &filter.Require{RegexContent: "[0-"},
		},
		"invalid urlCommands": {
			errRegex:    `type: .* <invalid>`,
			urlCommands: &filter.URLCommandSlice{{Type: "foo"}},
		},
		"all errs": {
			errRegex:    `url: <required>`,
			url:         stringPtr(""),
			require:     &filter.Require{RegexContent: "[0-"},
			urlCommands: &filter.URLCommandSlice{{Type: "foo"}},
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lookup := testLookup(false, false)
			if tc.lType != nil {
				lookup.Type = *tc.lType
			}
			if tc.url != nil {
				lookup.URL = *tc.url
			}
			if tc.require != nil {
				lookup.Require = tc.require
			}
			if tc.urlCommands != nil {
				lookup.URLCommands = *tc.urlCommands
			}

			// WHEN CheckValues is called
			err := lookup.CheckValues("")

			// THEN it err's when expected
			e := util.ErrorToString(err)
			re := regexp.MustCompile(tc.errRegex)
			match := re.MatchString(e)
			if !match {
				t.Fatalf("want match for %q\nnot: %q",
					tc.errRegex, e)
			}
		})
	}
}

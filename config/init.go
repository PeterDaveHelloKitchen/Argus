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

package config

import (
	"fmt"
	"os"

	dbtype "github.com/release-argus/Argus/db/types"
	"github.com/release-argus/Argus/service"
	deployedver "github.com/release-argus/Argus/service/deployed_version"
	svcstatus "github.com/release-argus/Argus/service/status"
	"github.com/release-argus/Argus/util"
	"gopkg.in/yaml.v3"
)

// LogInit for Argus.
func LogInit(log *util.JLog) {
	if jLog != nil {
		return
	}

	jLog = log
	service.LogInit(jLog)
}

// Init will hand out the appropriate Defaults.X and HardDefaults.X pointer(s)
func (c *Config) Init() {
	c.OrderMutex.RLock()
	defer c.OrderMutex.RUnlock()

	c.HardDefaults.SetDefaults()
	c.Settings.SetDefaults()

	if c.Defaults.Service.DeployedVersionLookup == nil {
		c.Defaults.Service.DeployedVersionLookup = &deployedver.Lookup{}
	}
	c.Defaults.Service.Convert()
	c.HardDefaults.Service.Status.SaveChannel = c.SaveChannel

	jLog.SetTimestamps(*c.Settings.GetLogTimestamps())
	jLog.SetLevel(c.Settings.GetLogLevel())

	i := 0
	for _, name := range c.Order {
		i++
		jLog.Debug(fmt.Sprintf("%d/%d %s Init", i, len(c.Service), name),
			util.LogFrom{}, true)
		c.Service[name].Init(
			&c.Defaults.Service, &c.HardDefaults.Service,
			&c.Notify, &c.Defaults.Notify, &c.HardDefaults.Notify,
			&c.WebHook, &c.Defaults.WebHook, &c.HardDefaults.WebHook)
	}

	// c.Notify
	if c.Notify != nil {
		for key := range c.Notify {
			// DefaultIfNil to handle testing. CheckValues will pick up on this nil
			c.Notify[key].Defaults = c.Defaults.Notify[c.Notify[key].Type]
			c.Notify[key].HardDefaults = c.HardDefaults.Notify[c.Notify[key].Type]
		}
	}
	// c.WebHook
	if c.WebHook != nil {
		for key := range c.WebHook {
			c.WebHook[key].Defaults = &c.Defaults.WebHook
			c.WebHook[key].HardDefaults = &c.HardDefaults.WebHook
		}
	}
}

// Load `file` as Config.
func (c *Config) Load(file string, flagset *map[string]bool, log *util.JLog) {
	c.File = file
	// Give the log to the other packages
	LogInit(log)
	c.Settings.NilUndefinedFlags(flagset)

	//#nosec G304 -- Loading the file asked for by the user
	data, err := os.ReadFile(file)
	msg := fmt.Sprintf("Error reading %q\n%s", file, err)
	jLog.Fatal(msg, util.LogFrom{}, err != nil)

	err = yaml.Unmarshal(data, c)
	msg = fmt.Sprintf("Unmarshal of %q failed\n%s", file, err)
	jLog.Fatal(msg, util.LogFrom{}, err != nil)

	c.GetOrder(data)

	databaseChannel := make(chan dbtype.Message, 32)
	c.DatabaseChannel = &databaseChannel

	saveChannel := make(chan bool, 32)
	c.SaveChannel = &saveChannel

	for key := range c.Service {
		c.Service[key].ID = key
		c.Service[key].Status = svcstatus.Status{
			DatabaseChannel: c.DatabaseChannel,
			SaveChannel:     c.SaveChannel}
	}
	c.HardDefaults.Service.Status.DatabaseChannel = c.DatabaseChannel
	c.HardDefaults.Service.Status.SaveChannel = c.SaveChannel

	// SaveHandler that listens for calls to save config changes.
	go c.SaveHandler()

	c.Init()
	c.CheckValues()
}

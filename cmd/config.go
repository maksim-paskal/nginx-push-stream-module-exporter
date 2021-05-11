/*
Copyright paskal.maksim@gmail.com
Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"flag"
)

type AppConfig struct {
	Version          string
	LogLevel         *string
	LogPretty        *bool
	WebListenAddress *string
	WebTelemetryPath *string
	NginxAddress     *string
	NginxStatsPath   *string
}

func (ac *AppConfig) String() string {
	b, err := json.MarshalIndent(ac, "", " ")
	if err != nil {
		return err.Error()
	}

	return string(b)
}

var appConfig = &AppConfig{
	Version:          gitVersion,
	LogLevel:         flag.String("log.level", "INFO", "log level"),
	LogPretty:        flag.Bool("log.pretty", false, "log in pretty format"),
	WebListenAddress: flag.String("web.listen-address", ":8102", "Address on which to expose metrics and web interface"),
	WebTelemetryPath: flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics"),
	NginxAddress:     flag.String("nginx.address", "http://127.0.0.1", "nginx address"),
	NginxStatsPath:   flag.String("nginx.stats-path", "/channels-stats", "statistics path"),
}

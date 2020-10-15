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
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

var (
	gitVersion string = "dev"
	buildTime  string
)

func main() {
	flag.Parse()

	logLevel, err := log.ParseLevel(*appConfig.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(logLevel)

	if *appConfig.LogPretty {
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if log.GetLevel() == log.DebugLevel {
		log.SetReportCaller(true)
	}

	log.Infof("Starting %s...", appConfig.Version)

	log.Debugf("using config:\n%s", appConfig.String())

	exporter := NewExporter()

	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector(moduleName))

	http.Handle(*appConfig.WebTelemetryPath, promhttp.Handler())
	log.Infoln("Listening on", *appConfig.WebListenAddress)
	log.Fatal(http.ListenAndServe(*appConfig.WebListenAddress, nil))
}

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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Exporter struct {
	up                prometheus.Gauge
	channels          prometheus.Gauge
	wildcardChannels  prometheus.Gauge
	publishedMessages prometheus.Gauge
	storedMessages    prometheus.Gauge
	messagesInTrash   prometheus.Gauge
	channelsInDelete  prometheus.Gauge
	channelsInTrash   prometheus.Gauge
	subscribers       prometheus.Gauge
}

type NginxPushModuleStatistics struct {
	Channels          int `json:"channels"`
	WildcardChannels  int `json:"wildcard_channels"`
	PublishedMessages int `json:"published_messages"`
	StoredMessages    int `json:"stored_messages"`
	MessagesInTrash   int `json:"messages_in_trash"`
	ChannelsInDelete  int `json:"channels_in_delete"`
	ChannelsInTrash   int `json:"channels_in_trash"`
	Subscribers       int `json:"subscribers"`
}

func NewExporter() *Exporter {
	return &Exporter{
		up: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: moduleName,
				Name:      "up",
				Help:      "The current health status of the server (1 = UP, 0 = DOWN).",
			},
		),
		channels: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: moduleName,
				Name:      "channels",
				Help:      "Numbers of channels",
			},
		),
		wildcardChannels: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: moduleName,
				Name:      "wildcard_channels",
				Help:      "wildcard_channels",
			},
		),
		publishedMessages: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: moduleName,
				Name:      "published_messages",
				Help:      "published_messages",
			},
		),
		storedMessages: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: moduleName,
				Name:      "stored_messages",
				Help:      "stored_messages",
			},
		),
		messagesInTrash: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: moduleName,
				Name:      "messages_in_trash",
				Help:      "messages_in_trash",
			},
		),
		channelsInDelete: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: moduleName,
				Name:      "channels_in_delete",
				Help:      "channels_in_delete",
			},
		),
		channelsInTrash: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: moduleName,
				Name:      "channels_in_trash",
				Help:      "channels_in_trash",
			},
		),
		subscribers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: moduleName,
				Name:      "subscribers",
				Help:      "subscribers",
			},
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.up.Describe(ch)
	e.channels.Describe(ch)
	e.wildcardChannels.Describe(ch)
	e.publishedMessages.Describe(ch)
	e.storedMessages.Describe(ch)
	e.messagesInTrash.Describe(ch)
	e.channelsInDelete.Describe(ch)
	e.channelsInTrash.Describe(ch)
	e.subscribers.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	up := e.Scrape()

	ch <- prometheus.MustNewConstMetric(e.up.Desc(), prometheus.GaugeValue, up)

	e.channels.Collect(ch)
	e.wildcardChannels.Collect(ch)
	e.publishedMessages.Collect(ch)
	e.storedMessages.Collect(ch)
	e.messagesInTrash.Collect(ch)
	e.channelsInDelete.Collect(ch)
	e.channelsInTrash.Collect(ch)
	e.subscribers.Collect(ch)
}

func (e *Exporter) Scrape() float64 {
	ctx := context.Background()

	client := &http.Client{}

	url := fmt.Sprintf("%s%s", *appConfig.NginxAddress, *appConfig.NginxStatsPath)

	if appConfig.IsDebugLevel {
		log.Debugf(url)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error(err)

		return 0
	}

	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)

		return 0
	}
	defer resp.Body.Close()

	if err != nil {
		log.Error(err)

		return 0
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("error code not ok %d", resp.StatusCode)

		return 0
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)

		return 0
	}

	if appConfig.IsDebugLevel {
		log.Debug(string(byteValue))
	}

	data := NginxPushModuleStatistics{}

	err = json.Unmarshal(byteValue, &data)

	if err != nil {
		log.Error(err)

		return 0
	}

	e.channels.Set(float64(data.Channels))
	e.wildcardChannels.Set(float64(data.WildcardChannels))
	e.publishedMessages.Set(float64(data.PublishedMessages))
	e.storedMessages.Set(float64(data.StoredMessages))
	e.messagesInTrash.Set(float64(data.MessagesInTrash))
	e.channelsInDelete.Set(float64(data.ChannelsInDelete))
	e.channelsInTrash.Set(float64(data.ChannelsInTrash))
	e.subscribers.Set(float64(data.Subscribers))

	return 1
}

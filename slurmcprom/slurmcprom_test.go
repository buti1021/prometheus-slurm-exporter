package slurmcprom

// SPDX-FileCopyrightText: 2023 Rivos Inc.
//
// SPDX-License-Identifier: Apache-2.0

import (
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rivosinc/prometheus-slurm-exporter/exporter"
	"github.com/stretchr/testify/assert"
)

func TestCtoGoNodeMetrics(t *testing.T) {
	assert := assert.New(t)
	collector := NewNodeFetcher(0)
	defer collector.Deinit()
	metrics, err := collector.CToGoMetricConvert()
	assert.NoError(err)
	assert.Positive(len(metrics))
}

func TestCtoGoNodeMetricsTwice(t *testing.T) {
	assert := assert.New(t)
	// force cache misses
	collector := NewNodeFetcher(0)
	defer collector.Deinit()
	metrics, err := collector.CToGoMetricConvert()
	assert.NoError(err)
	assert.Positive(len(metrics))
	// tests cached partition & node info data path
	metrics, err = collector.CToGoMetricConvert()
	assert.NoError(err)
	assert.Positive(len(metrics))
}

func TestNodeCollectorCFetcher(t *testing.T) {
	if os.Getenv("TEST_CLUSTER") != "true" {
		return
	}
	assert := assert.New(t)
	config, err := exporter.NewConfig(new(exporter.CliFlags))
	assert.Nil(err)
	config.PollLimit = 10
	nc := exporter.NewNodeCollecter(config)
	// cache miss, use our mock fetcher
	nc.SetFetcher(NewNodeFetcher(config.PollLimit))
	metricChan := make(chan prometheus.Metric)
	go func() {
		nc.Collect(metricChan)
		close(metricChan)
	}()
	metrics := make([]prometheus.Metric, 0)
	for m, ok := <-metricChan; ok; m, ok = <-metricChan {
		metrics = append(metrics, m)
		t.Logf("Received metric %s", m.Desc().String())
	}
	assert.NotEmpty(metrics)
}

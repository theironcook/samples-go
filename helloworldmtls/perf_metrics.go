package helloworldmtls

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
)

// Probably delete this after initial dev time.  Helpful to see processing rates quickly
type PerfMetricsActivities struct {
	metrics *PerfMetrics
}

type RecordMetricParams struct {
	Succeeded bool
}

func (a *PerfMetricsActivities) RecordMetric(ctx context.Context, params RecordMetricParams) error {
	if a.metrics == nil {
		a.metrics = &PerfMetrics{}
		go a.metrics.startCollecting()
	}
	if params.Succeeded {
		a.metrics.recordSucceeded()
	} else {
		a.metrics.recordFailed()
	}
	return nil
}

type PerfMetrics struct {
	mutex     sync.Mutex
	started   bool
	succeeded int
	failed    int
	// Continuous collected data (data in each tick)
	contTicks     int
	contSucceeded int
	contFailed    int
}

func (m *PerfMetrics) recordSucceeded() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.succeeded++
}

func (m *PerfMetrics) recordFailed() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.failed++
}

func (m *PerfMetrics) startCollecting() {
	m.mutex.Lock()
	if m.started {
		m.mutex.Unlock()
		return
	}
	m.started = true
	m.mutex.Unlock()

	// Ensure the ./log directory exists
	err := os.MkdirAll("./log", os.ModePerm)
	if err != nil {
		err = fmt.Errorf("failed to collect performance metrics. Failed to create log directory: %v", err)
		slog.Error(err.Error())
		return
	}

	file, err := os.OpenFile("./log/perf_metrics.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error(fmt.Errorf("failed to open metrics.log: %w", err).Error())
		return
	}
	encoder := json.NewEncoder(file)
	defer file.Close()

	ticker := time.NewTicker(time.Second * 1)

	for range ticker.C {
		m.mutex.Lock()
		// only collect ticks with data - assumes processes go as fast as they can given throttling constraints
		if m.succeeded+m.failed > 0 {
			m.contSucceeded += m.succeeded
			m.contFailed += m.failed
			m.contTicks++
			rollingAvg := float64(m.contSucceeded+m.contFailed) / float64(m.contTicks)
			slog.Info("perf_metrics", "succeeded", m.succeeded, "failed", m.failed, "rollingAvg", fmt.Sprintf("%.2f", rollingAvg))
			logEntry := map[string]interface{}{
				"timestamp":  time.Now().Format(time.RFC3339),
				"succeeded":  m.succeeded,
				"failed":     m.failed,
				"rollingAvg": fmt.Sprintf("%.2f", rollingAvg),
			}
			if err := encoder.Encode(logEntry); err != nil {
				slog.Error(fmt.Errorf("failed to write perf metric log entry: %w", err).Error())
			}
			m.failed = 0
			m.succeeded = 0
		} else {
			m.contTicks = 0
			m.contSucceeded = 0
			m.contFailed = 0
		}
		m.mutex.Unlock()
	}
}

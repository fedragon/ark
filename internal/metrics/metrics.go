package metrics

import (
	p "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totals = promauto.NewCounterVec(
		p.CounterOpts{
			Name:      "totals",
			Namespace: "ark",
			Help:      "The total number of files affected by a given operation",
		},
		[]string{"duplicate"},
	)

	TotalDuplicates = totals.With(p.Labels{"duplicate": "true"})
	TotalImported   = totals.With(p.Labels{"duplicate": "false"})

	duration = promauto.NewSummaryVec(
		p.SummaryOpts{
			Name:       "duration_ms",
			Namespace:  "ark",
			Help:       "The duration of a given operation, in milliseconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
		},
		[]string{"operation"})

	CopyFileDurationMs   = duration.With(p.Labels{"operation": "copy_file"})
	GetDurationMs        = duration.With(p.Labels{"operation": "get"})
	StoreDurationMs      = duration.With(p.Labels{"operation": "store"})
	UploadFileDurationMs = duration.With(p.Labels{"operation": "upload_file"})
)

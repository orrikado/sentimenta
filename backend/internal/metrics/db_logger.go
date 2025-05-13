package metrics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Prometheus *Prometheus
	logger.Interface
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	fmt.Printf("sql: %v\n", sql)
	queryType := strings.ToUpper(strings.Fields(sql)[0])
	duration := time.Since(begin).Seconds()

	l.Prometheus.DBQueryDuration.WithLabelValues(queryType).Observe(duration)
	if err != nil {
		l.Prometheus.DBErrorsTotal.WithLabelValues(queryType).Inc()
	}

	l.Interface.Trace(ctx, begin, fc, err)
}

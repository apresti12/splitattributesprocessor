package splitAttributesProcessor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
	"strings"
)

type splitAttrsProcessor struct {
	logger       *zap.Logger
	config       *Config
	host         component.Host
	cancel       context.CancelFunc
	nextConsumer consumer.Metrics
}

func (r *splitAttrsProcessor) shutdown(ctx context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}
	r.logger.Info("Stopping SplitAttributesProcessor")
	return nil
}

func (r *splitAttrsProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	r.logger.Info("SplitAttributesProcessor is processing metrics")
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		metrics := md.ResourceMetrics().At(i)
		for j := 0; j < metrics.ScopeMetrics().Len(); j++ {
			scopeMetrics := metrics.ScopeMetrics().At(j)
			for k := 0; k < scopeMetrics.Metrics().Len(); k++ {
				innerMetric := scopeMetrics.Metrics().At(k)
				metricLength := 0
				if innerMetric.Type() == pmetric.MetricTypeGauge {
					metricLength = innerMetric.Gauge().DataPoints().Len()
				} else {
					metricLength = innerMetric.Sum().DataPoints().Len()
				}
				for l := 0; l < metricLength; l++ {
					var datapoint pmetric.NumberDataPoint
					if innerMetric.Type() == pmetric.MetricTypeGauge {
						datapoint = innerMetric.Gauge().DataPoints().At(l)
					} else {
						datapoint = innerMetric.Sum().DataPoints().At(l)
					}
					concatenatedHashes, ok := datapoint.Attributes().Get(r.config.AttributeKey)
					if !ok {
						continue
					}
					hashList := splitHashes(concatenatedHashes.Str(), r.config.Delimiter)
					for _, hash := range hashList {
						var newDp pmetric.NumberDataPoint
						if innerMetric.Type() == pmetric.MetricTypeGauge {
							newDp = innerMetric.Gauge().DataPoints().AppendEmpty()
						} else {
							newDp = innerMetric.Sum().DataPoints().AppendEmpty()
						}
						datapoint.CopyTo(newDp)
						newDp.Attributes().PutStr("hash", hash)
						newDp.Attributes().Remove(r.config.AttributeKey)
					}
					if innerMetric.Type() == pmetric.MetricTypeGauge {
						innerMetric.Gauge().DataPoints().RemoveIf(func(dp pmetric.NumberDataPoint) bool {
							_, ok := dp.Attributes().Get(r.config.AttributeKey)
							return ok
						})
					} else {
						innerMetric.Sum().DataPoints().RemoveIf(func(dp pmetric.NumberDataPoint) bool {
							_, ok := dp.Attributes().Get(r.config.AttributeKey)
							return ok
						})
					}
				}
			}
		}
	}
	return md, nil
}

func handleSplit(HashString string, delimiter string) {

	return
}

func splitHashes(concatenatedHashes string, delimiter string) []string {
	return strings.Split(concatenatedHashes, delimiter)
}

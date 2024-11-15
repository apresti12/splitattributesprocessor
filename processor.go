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
	logger *zap.Logger
	config *Config
}

func (r *splitAttrsProcessor) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (r *splitAttrsProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (r splitAttrsProcessor) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	r.processMetrics(ctx, md)
	return nil
}

func (r splitAttrsProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{}
}

func (r splitAttrsProcessor) processMetrics(ctx context.Context, md pmetric.Metrics) {
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		metrics := md.ResourceMetrics().At(i)
		for j := 0; j < metrics.ScopeMetrics().Len(); j++ {
			scopeMetrics := metrics.ScopeMetrics().At(j)
			for k := 0; k < scopeMetrics.Metrics().Len(); k++ {
				innerMetric := scopeMetrics.Metrics().At(k)
				for l := 0; l < innerMetric.Sum().DataPoints().Len(); l++ {
					datapoint := innerMetric.Sum().DataPoints().At(l)
					concatenatedHashes, ok := datapoint.Attributes().Get(r.config.AttributeKey)
					if !ok {
						continue
					}
					hashList := splitHashes(concatenatedHashes.Str(), r.config.Delimiter)
					for _, hash := range hashList {
						newDp := innerMetric.Sum().DataPoints().AppendEmpty()
						datapoint.CopyTo(newDp)
						newDp.Attributes().PutStr("hash", hash)
						newDp.Attributes().Remove(r.config.AttributeKey)
					}
					innerMetric.Sum().DataPoints().RemoveIf(func(dp pmetric.NumberDataPoint) bool {
						_, ok := dp.Attributes().Get(r.config.AttributeKey)
						return ok
					})
				}
			}
		}
	}
}

func splitHashes(concatenatedHashes string, delimiter string) []string {
	return strings.Split(concatenatedHashes, delimiter)
}

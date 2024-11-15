package splitAttributesProcessor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

var (
	typeStr = component.MustNewType("splitattributesprocessor")
)

const (
	defaultDelimiter    = ";"
	defaultAttributeKey = "hashes"
)

func createDefaultConfig() component.Config {
	return &Config{
		Delimiter:    defaultDelimiter,
		AttributeKey: defaultAttributeKey,
	}
}

func createMetricsProcessor(ctx context.Context, params processor.Settings, baseCfg component.Config, consumer consumer.Metrics) (processor.Metrics, error) {
	logger := params.Logger
	cfg := baseCfg.(*Config)
	mProcessor := &splitAttrsProcessor{
		logger:       logger,
		config:       cfg,
		nextConsumer: consumer,
	}
	return mProcessor, nil
}

func NewFactory() processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig,
		processor.WithMetrics(createMetricsProcessor, component.StabilityLevelAlpha),
	)
}

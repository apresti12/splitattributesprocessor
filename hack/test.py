from opentelemetry import trace
from opentelemetry import metrics
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.exporter.otlp.proto.grpc.metric_exporter import OTLPMetricExporter
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.resourcedetector.gcp_resource_detector import GoogleCloudResourceDetector


def main():
    resource = Resource(attributes={
            SERVICE_NAME: "test"
        })
    exporter =  OTLPMetricExporter(endpoint="http://localhost:4317")
    reader = PeriodicExportingMetricReader(exporter=exporter, export_interval_millis=1000)
    provider = MeterProvider(resource=resource, metric_readers=[reader])
    metrics.set_meter_provider(provider)
    m = metrics.get_meter("test.meter")
    c = m.create_counter("test_counter", unit="units", description="test")
    for i in range(10):
        c.add(1, {"hashes":"item1;item2;item3;item4"})
        print("added counter")
    print("done with counter")
    print("starting gauge")
    g = m.create_gauge("test_gauge", unit="units", description="test")
    for i in range(10):
        g.set(1, {"hashes":"item_a;item_b;item_c;item_d"})
    print("done with gauge")
if __name__ == "__main__":
    main()

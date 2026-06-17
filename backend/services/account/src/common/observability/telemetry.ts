import { NodeSDK } from '@opentelemetry/sdk-node';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { OTLPMetricExporter } from '@opentelemetry/exporter-metrics-otlp-grpc';
import { OTLPLogExporter } from '@opentelemetry/exporter-logs-otlp-grpc';
import { PeriodicExportingMetricReader } from '@opentelemetry/sdk-metrics';
import {
  LoggerProvider,
  BatchLogRecordProcessor,
} from '@opentelemetry/sdk-logs';
import { logs } from '@opentelemetry/api-logs';
import { resourceFromAttributes } from '@opentelemetry/resources';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';

const resource = resourceFromAttributes({
  [ATTR_SERVICE_NAME]: process.env.SERVICE_NAME,
});

// Logs setup
const loggerProvider = new LoggerProvider({
  resource,
  processors: [
    new BatchLogRecordProcessor(
      new OTLPLogExporter({ url: process.env.OTEL_EXPORTER_ENDPOINT }),
    ),
  ],
});
logs.setGlobalLoggerProvider(loggerProvider);

const sdk = new NodeSDK({
  resource,
  traceExporter: new OTLPTraceExporter({
    url: process.env.OTEL_EXPORTER_ENDPOINT,
  }),
  metricReader: new PeriodicExportingMetricReader({
    exporter: new OTLPMetricExporter({
      url: process.env.OTEL_EXPORTER_ENDPOINT,
    }),
  }),
  instrumentations: [getNodeAutoInstrumentations()],
});

sdk.start();

process.on('SIGTERM', () => {
  Promise.all([sdk.shutdown(), loggerProvider.shutdown()]).finally(() =>
    process.exit(0),
  );
});

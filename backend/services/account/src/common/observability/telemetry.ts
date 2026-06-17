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

export interface TelemetryOptions {
  enabled: boolean;
  serviceName: string;
  exporterEndpoint: string;
}

export interface TelemetryHandle {
  shutdown: () => Promise<void>;
}

export function startTelemetry(
  options: TelemetryOptions,
): TelemetryHandle | null {
  console.log('Telemetry options:', options);

  if (!options.enabled) {
    console.log('Telemetry disabled, skipping');
    return null;
  }

  const resource = resourceFromAttributes({
    [ATTR_SERVICE_NAME]: options.serviceName,
  });

  const loggerProvider = new LoggerProvider({
    resource,
    processors: [
      new BatchLogRecordProcessor(
        new OTLPLogExporter({ url: options.exporterEndpoint }),
      ),
    ],
  });
  logs.setGlobalLoggerProvider(loggerProvider);

  const sdk = new NodeSDK({
    resource,
    traceExporter: new OTLPTraceExporter({ url: options.exporterEndpoint }),
    metricReader: new PeriodicExportingMetricReader({
      exporter: new OTLPMetricExporter({ url: options.exporterEndpoint }),
    }),
    instrumentations: [getNodeAutoInstrumentations()],
  });

  sdk.start();
  console.log('OTel SDK started');

  return {
    shutdown: async () => {
      await Promise.all([sdk.shutdown(), loggerProvider.shutdown()]);
    },
  };
}

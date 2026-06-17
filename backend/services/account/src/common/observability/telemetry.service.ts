import { Injectable, OnApplicationShutdown } from '@nestjs/common';
import { Config } from '../../config/config';
import { startTelemetry, TelemetryHandle } from './telemetry';

@Injectable()
export class TelemetryService implements OnApplicationShutdown {
  private handle: TelemetryHandle | null = null;

  constructor(private readonly config: Config) {
    this.handle = startTelemetry({
      enabled: this.config.otel.enabled,
      serviceName: this.config.serviceName,
      exporterEndpoint: this.config.otel.exporterEndpoint,
    });
  }

  async onApplicationShutdown(): Promise<void> {
    await this.handle?.shutdown();
  }
}

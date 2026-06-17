import { TelemetryService } from './telemetry.service';
import { Config } from '../../config/config';
import * as telemetryModule from './telemetry';

jest.mock('./telemetry');

const mockConfig = {
  serviceName: 'account',
  otel: {
    enabled: true,
    exporterEndpoint: 'http://localhost:4317',
  },
} as Config;

describe('TelemetryService', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('constructor', () => {
    it('should call startTelemetry with correct options', () => {
      const startTelemetrySpy = jest.spyOn(telemetryModule, 'startTelemetry');

      new TelemetryService(mockConfig);

      expect(startTelemetrySpy).toHaveBeenCalledWith({
        enabled: true,
        serviceName: 'account',
        exporterEndpoint: 'http://localhost:4317',
      });
    });
  });

  describe('onApplicationShutdown', () => {
    it('should call handle.shutdown when handle exists', async () => {
      const shutdownSpy = jest.fn().mockResolvedValue(undefined);
      jest.spyOn(telemetryModule, 'startTelemetry').mockReturnValue({
        shutdown: shutdownSpy,
      });

      const service = new TelemetryService(mockConfig);
      await service.onApplicationShutdown();

      expect(shutdownSpy).toHaveBeenCalled();
    });

    it('should not throw when handle is null', async () => {
      jest.spyOn(telemetryModule, 'startTelemetry').mockReturnValue(null);

      const service = new TelemetryService(mockConfig);

      await expect(service.onApplicationShutdown()).resolves.not.toThrow();
    });
  });
});

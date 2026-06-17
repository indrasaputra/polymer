import { validateSync } from 'class-validator';
import { plainToInstance } from 'class-transformer';
import { Config } from './config';

const validConfig = {
  env: 'development',
  serviceName: 'account',
  port: 9001,
  webhookSecret: 'webhook-secret',
  postgre: {
    url: 'postgresql://postgres:postgres@localhost:54322/postgres',
    schema: 'account',
  },
  supabase: {
    jwksUrl: 'http://localhost:54321/auth/v1/.well-known/jwks.json',
  },
  otel: {
    exporterEndpoint: 'http://localhost:4317',
  },
};

function buildConfig(overrides: Record<string, unknown> = {}): Config {
  return plainToInstance(Config, { ...validConfig, ...overrides });
}

describe('Config', () => {
  describe('validation', () => {
    it('should pass validation with valid config', () => {
      const config = buildConfig();

      const errors = validateSync(config);

      expect(errors).toHaveLength(0);
    });

    it('should fail validation when webhookSecret is missing', () => {
      const config = buildConfig({ webhookSecret: undefined });

      const errors = validateSync(config);

      expect(errors.some((e) => e.property === 'webhookSecret')).toBe(true);
    });

    it('should fail validation when env is not one of allowed values', () => {
      const config = buildConfig({ env: 'staging-invalid' });

      const errors = validateSync(config);

      expect(errors.some((e) => e.property === 'env')).toBe(true);
    });

    it('should fail validation when postgre.url is missing', () => {
      const config = buildConfig({ postgre: { schema: 'account' } });

      const errors = validateSync(config);

      expect(errors.some((e) => e.property === 'postgre')).toBe(true);
    });

    it('should fail validation when supabase.jwksUrl is missing', () => {
      const config = buildConfig({ supabase: {} });

      const errors = validateSync(config);

      expect(errors.some((e) => e.property === 'supabase')).toBe(true);
    });

    it('should fail validation when otel.exporterEndpoint is missing', () => {
      const config = buildConfig({ otel: { enabled: false } });

      const errors = validateSync(config);

      expect(errors.some((e) => e.property === 'otel')).toBe(true);
    });

    it('should use default values when optional fields are omitted', () => {
      const config = plainToInstance(Config, {
        webhookSecret: 'secret',
        postgre: { url: 'postgresql://localhost:5432/postgres' },
        supabase: { jwksUrl: 'http://localhost/jwks.json' },
        otel: { exporterEndpoint: 'http://localhost:4317' },
      });

      expect(config.env).toBe('development');
      expect(config.serviceName).toBe('account');
      expect(config.port).toBe(9001);
      expect(config.postgre.schema).toBe('public');
    });
  });

  describe('isDevelopment', () => {
    it('should return true when env is development', () => {
      const config = buildConfig({ env: 'development' });

      expect(config.isDevelopment).toBe(true);
    });

    it('should return false when env is not development', () => {
      const config = buildConfig({ env: 'production' });

      expect(config.isDevelopment).toBe(false);
    });
  });

  describe('isStaging', () => {
    it('should return true when env is staging', () => {
      const config = buildConfig({ env: 'staging' });

      expect(config.isStaging).toBe(true);
    });

    it('should return false when env is not staging', () => {
      const config = buildConfig({ env: 'development' });

      expect(config.isStaging).toBe(false);
    });
  });

  describe('isProduction', () => {
    it('should return true when env is production', () => {
      const config = buildConfig({ env: 'production' });

      expect(config.isProduction).toBe(true);
    });

    it('should return false when env is not production', () => {
      const config = buildConfig({ env: 'development' });

      expect(config.isProduction).toBe(false);
    });
  });

  describe('transform', () => {
    it.each([
      ['true', true],
      ['1', true],
      ['yes', true],
      ['y', true],
      ['false', false],
      ['0', false],
      ['no', false],
      ['n', false],
      [undefined, false],
    ])('should parse %s as %s', (input, expected) => {
      const config = plainToInstance(Config, {
        ...validConfig,
        otel: { enabled: input, exporterEndpoint: 'http://localhost:4317' },
      });

      expect(config.otel.enabled).toBe(expected);
    });
  });
});

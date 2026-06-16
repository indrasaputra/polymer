// config.spec.ts

describe('config', () => {
  const ORIGINAL_ENV = process.env.ENV;

  afterEach(() => {
    if (ORIGINAL_ENV === undefined) {
      delete process.env.ENV;
    } else {
      process.env.ENV = ORIGINAL_ENV;
    }
    jest.resetModules();
  });

  describe('ENV', () => {
    it('should default to "development" if ENV is not set', async () => {
      delete process.env.ENV;
      const { ENV } = await import('./config.js');
      expect(ENV).toBe('development');
    });

    it('should use the value of process.env.ENV when set', async () => {
      process.env.ENV = 'production';
      const { ENV } = await import('./config.js');
      expect(ENV).toBe('production');
    });
  });

  describe('isDevelopment', () => {
    it('should be true when ENV is "development"', async () => {
      process.env.ENV = 'development';
      const { isDevelopment } = await import('./config.js');
      expect(isDevelopment).toBe(true);
    });

    it('should be false when ENV is not "development"', async () => {
      process.env.ENV = 'production';
      const { isDevelopment } = await import('./config.js');
      expect(isDevelopment).toBe(false);
    });
  });

  describe('isStaging', () => {
    it('should be true when ENV is "staging"', async () => {
      process.env.ENV = 'staging';
      const { isStaging } = await import('./config.js');
      expect(isStaging).toBe(true);
    });

    it('should be false when ENV is not "staging"', async () => {
      process.env.ENV = 'production';
      const { isStaging } = await import('./config.js');
      expect(isStaging).toBe(false);
    });
  });

  describe('isProduction', () => {
    it('should be true when ENV is "production"', async () => {
      process.env.ENV = 'production';
      const { isProduction } = await import('./config.js');
      expect(isProduction).toBe(true);
    });

    it('should be false when ENV is not "production"', async () => {
      process.env.ENV = 'development';
      const { isProduction } = await import('./config.js');
      expect(isProduction).toBe(false);
    });
  });
});

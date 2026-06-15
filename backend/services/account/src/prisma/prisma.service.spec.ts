import { ConfigService } from '@nestjs/config';
import { PrismaService } from './prisma.service';

jest.mock('pg', () => ({
  Pool: jest.fn().mockImplementation(() => ({
    end: jest.fn().mockResolvedValue(undefined),
  })),
}));

jest.mock('@prisma/adapter-pg', () => ({
  PrismaPg: jest.fn().mockImplementation(() => ({ provider: 'pg' })),
}));

jest.mock('../../generated/prisma/client', () => ({
  PrismaClient: jest.fn().mockImplementation(function () {
    this.$connect = jest.fn().mockResolvedValue(undefined);
    this.$disconnect = jest.fn().mockResolvedValue(undefined);
  }),
}));

const mockConfigService = (values: Record<string, string | undefined>) =>
  ({
    get: jest.fn().mockImplementation((key: string) => values[key]),
  }) as unknown as ConfigService;

const validConfig = {
  DATABASE_URL: 'postgresql://postgres:postgres@localhost:54322/postgres',
  DATABASE_SCHEMA: 'account',
};

describe('PrismaService', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('constructor', () => {
    it('should throw if DATABASE_URL is missing', () => {
      const configService = mockConfigService({ DATABASE_URL: undefined });

      expect(() => new PrismaService(configService)).toThrow(
        'DATABASE_URL environment variable is missing!',
      );
    });

    it('should instantiate successfully with DATABASE_URL', () => {
      const configService = mockConfigService(validConfig);

      expect(() => new PrismaService(configService)).not.toThrow();
    });

    it('should default schema to public if DATABASE_SCHEMA is not set', () => {
      const { PrismaPg } = jest.requireMock('@prisma/adapter-pg');

      const configService = mockConfigService({
        DATABASE_URL: validConfig.DATABASE_URL,
        DATABASE_SCHEMA: undefined,
      });

      new PrismaService(configService);

      expect(PrismaPg).toHaveBeenCalledWith(expect.anything(), {
        schema: 'public',
      });
    });

    it('should use DATABASE_SCHEMA if set', () => {
      const { PrismaPg } = jest.requireMock('@prisma/adapter-pg');

      const configService = mockConfigService(validConfig);

      new PrismaService(configService);

      expect(PrismaPg).toHaveBeenCalledWith(expect.anything(), {
        schema: 'account',
      });
    });
  });

  describe('onModuleInit', () => {
    it('should call $connect', async () => {
      const configService = mockConfigService(validConfig);
      const service = new PrismaService(configService);
      await service.onModuleInit();

      expect(service.$connect).toHaveBeenCalled();
    });
  });

  describe('onModuleDestroy', () => {
    it('should call $disconnect and pool.end', async () => {
      const { Pool } = jest.requireMock('pg');
      const poolEndSpy = jest.fn().mockResolvedValue(undefined);
      Pool.mockImplementation(() => ({ end: poolEndSpy }));

      const configService = mockConfigService(validConfig);
      const service = new PrismaService(configService);
      await service.onModuleDestroy();

      expect(service.$disconnect).toHaveBeenCalled();
      expect(poolEndSpy).toHaveBeenCalled();
    });
  });
});

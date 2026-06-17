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

jest.mock('../config/config', () => ({
  Config: {
    DATABASE_URL: 'postgresql://postgres:postgres@localhost:54322/postgres',
    DATABASE_SCHEMA: 'account',
  },
}));

describe('PrismaService', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('constructor', () => {
    it('should instantiate successfully', () => {
      expect(() => new PrismaService()).not.toThrow();
    });

    it('should use DATABASE_SCHEMA from Config', () => {
      const { PrismaPg } = jest.requireMock('@prisma/adapter-pg');

      new PrismaService();

      expect(PrismaPg).toHaveBeenCalledWith(expect.anything(), {
        schema: 'account',
      });
    });

    it('should use DATABASE_URL from Config for pool connection', () => {
      const { Pool } = jest.requireMock('pg');

      new PrismaService();

      expect(Pool).toHaveBeenCalledWith({
        connectionString:
          'postgresql://postgres:postgres@localhost:54322/postgres',
      });
    });
  });

  describe('onModuleInit', () => {
    it('should call $connect', async () => {
      const service = new PrismaService();
      await service.onModuleInit();

      expect(service.$connect).toHaveBeenCalled();
    });
  });

  describe('onModuleDestroy', () => {
    it('should call $disconnect and pool.end', async () => {
      const { Pool } = jest.requireMock('pg');
      const poolEndSpy = jest.fn().mockResolvedValue(undefined);
      Pool.mockImplementation(() => ({ end: poolEndSpy }));

      const service = new PrismaService();
      await service.onModuleDestroy();

      expect(service.$disconnect).toHaveBeenCalled();
      expect(poolEndSpy).toHaveBeenCalled();
    });
  });
});

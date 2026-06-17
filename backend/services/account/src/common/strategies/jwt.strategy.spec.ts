import { ConfigService } from '@nestjs/config';
import { UnauthorizedException } from '@nestjs/common';

// Mock 'jwks-rsa' to completely block outbound network requests
jest.mock('jwks-rsa', () => ({
  passportJwtSecret: jest.fn().mockReturnValue(() => {
    return 'mocked-secret-or-key';
  }),
}));

import { JwtStrategy } from './jwt.strategy';
import { JwtToken } from '../dto/jwt-token.interface';

describe('JwtStrategy', () => {
  let configService: jest.Mocked<ConfigService>;

  beforeEach(() => {
    configService = {
      get: jest.fn(),
    } as unknown as jest.Mocked<ConfigService>;
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('Constructor Initialization', () => {
    it('should initialize successfully when JWKS URL is provided', () => {
      configService.get.mockReturnValue('https://supabase.co');

      const strategy = new JwtStrategy(configService);

      expect(strategy).toBeDefined();
      expect(configService.get).toHaveBeenCalledWith('SUPABASE_JWKS_URL');
    });

    it('should throw an error if SUPABASE_JWKS_URL is missing', () => {
      configService.get.mockReturnValue(undefined);

      expect(() => new JwtStrategy(configService)).toThrow(
        'SUPABASE_JWKS_URL is not configured',
      );
    });
  });

  describe('validate', () => {
    let strategy: JwtStrategy;

    beforeEach(() => {
      configService.get.mockReturnValue('https://supabase.co');
      strategy = new JwtStrategy(configService);
    });

    it('should return CurrentUser when payload is valid', () => {
      const validPayload: JwtToken = {
        sub: 'user-uuid-123',
        email: 'test@example.com',
      };

      const result = strategy.validate(validPayload);

      expect(result).toEqual({
        id: 'user-uuid-123',
        email: 'test@example.com',
      });
    });

    it('should throw UnauthorizedException if sub is missing', () => {
      const invalidPayload = {
        email: 'test@example.com',
      } as JwtToken;

      expect(() => strategy.validate(invalidPayload)).toThrow(
        UnauthorizedException,
      );
    });

    it('should throw UnauthorizedException if email is missing', () => {
      const invalidPayload = {
        sub: 'user-uuid-123',
      } as JwtToken;

      expect(() => strategy.validate(invalidPayload)).toThrow(
        UnauthorizedException,
      );
    });
  });
});

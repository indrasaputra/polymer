import { UnauthorizedException } from '@nestjs/common';
import { Config } from '../../config/config';

// Mock 'jwks-rsa' to completely block outbound network requests
jest.mock('jwks-rsa', () => ({
  passportJwtSecret: jest.fn().mockReturnValue(() => {
    return 'mocked-secret-or-key';
  }),
}));

const mockConfig = {
  supabase: {
    jwksUrl: 'http://localhost:54321/auth/v1/.well-known/jwks.json',
  },
} as Config;

import { JwtStrategy } from './jwt.strategy';
import { JwtToken } from '../dto/jwt-token.interface';

describe('JwtStrategy', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('Constructor Initialization', () => {
    it('should initialize successfully', () => {
      const strategy = new JwtStrategy(mockConfig);

      expect(strategy).toBeDefined();
    });

    it('should use jwksUri from Config', () => {
      const { passportJwtSecret } = jest.requireMock('jwks-rsa');

      new JwtStrategy(mockConfig);

      expect(passportJwtSecret).toHaveBeenCalledWith(
        expect.objectContaining({
          jwksUri: mockConfig.supabase.jwksUrl,
        }),
      );
    });
  });

  describe('validate', () => {
    let strategy: JwtStrategy;

    beforeEach(() => {
      strategy = new JwtStrategy(mockConfig);
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

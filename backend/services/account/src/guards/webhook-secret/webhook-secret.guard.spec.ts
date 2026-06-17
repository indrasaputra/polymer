import { ExecutionContext, UnauthorizedException } from '@nestjs/common';
import { WebhookSecretGuard } from './webhook-secret.guard';

jest.mock('../../config/config', () => ({
  Config: {
    WEBHOOK_SECRET: 'valid-secret',
  },
}));

const mockExecutionContext = (secret: string | undefined): ExecutionContext =>
  ({
    switchToHttp: () => ({
      getRequest: () => ({
        header: jest.fn().mockReturnValue(secret),
      }),
    }),
  }) as unknown as ExecutionContext;

describe('WebhookSecretGuard', () => {
  let guard: WebhookSecretGuard;

  beforeEach(() => {
    guard = new WebhookSecretGuard();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should be defined', () => {
    expect(guard).toBeDefined();
  });

  it('should return true when secret is valid', () => {
    const context = mockExecutionContext('valid-secret');

    expect(guard.canActivate(context)).toBe(true);
  });

  it('should throw UnauthorizedException when secret is invalid', () => {
    const context = mockExecutionContext('invalid-secret');

    expect(() => guard.canActivate(context)).toThrow(UnauthorizedException);
  });

  it('should throw UnauthorizedException with correct message when secret is invalid', () => {
    const context = mockExecutionContext('invalid-secret');

    expect(() => guard.canActivate(context)).toThrow('Secret is invalid');
  });

  it('should throw UnauthorizedException when secret is missing', () => {
    const context = mockExecutionContext(undefined);

    expect(() => guard.canActivate(context)).toThrow(UnauthorizedException);
  });
});

import { ExecutionContext } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import { AuthGuard } from '@nestjs/passport';
import { JwtAuthGuard } from './jwt.guards';
import { AuthType, AUTH_TYPE_KEY } from '../decorators/auth-type.decorator';

describe('JwtAuthGuard', () => {
  let guard: JwtAuthGuard;
  let reflector: jest.Mocked<Reflector>;
  let mockExecutionContext: jest.Mocked<ExecutionContext>;

  beforeEach(() => {
    reflector = {
      getAllAndOverride: jest.fn(),
    } as unknown as jest.Mocked<Reflector>;

    mockExecutionContext = {
      getHandler: jest.fn().mockReturnValue('mockHandler'),
      getClass: jest.fn().mockReturnValue('mockClass'),
    } as unknown as jest.Mocked<ExecutionContext>;

    guard = new JwtAuthGuard(reflector);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should be defined', () => {
    expect(guard).toBeDefined();
  });

  it('should return true immediately if AuthType is NONE', async () => {
    reflector.getAllAndOverride.mockReturnValue(AuthType.NONE);
    const superCanActivateSpy = jest.spyOn(
      AuthGuard('jwt').prototype,
      'canActivate',
    );

    const result = await guard.canActivate(mockExecutionContext);

    expect(reflector.getAllAndOverride).toHaveBeenCalledWith(AUTH_TYPE_KEY, [
      'mockHandler',
      'mockClass',
    ]);
    expect(result).toBe(true);
    expect(superCanActivateSpy).not.toHaveBeenCalled();
  });

  it('should call super.canActivate if AuthType is JWT', async () => {
    reflector.getAllAndOverride.mockReturnValue(AuthType.JWT);
    const superCanActivateSpy = jest
      .spyOn(AuthGuard('jwt').prototype, 'canActivate')
      .mockResolvedValue(true);

    const result = await guard.canActivate(mockExecutionContext);

    expect(reflector.getAllAndOverride).toHaveBeenCalled();
    expect(superCanActivateSpy).toHaveBeenCalledWith(mockExecutionContext);
    expect(result).toBe(true);
  });

  it('should default to AuthType.JWT and call super.canActivate if no metadata is found', async () => {
    reflector.getAllAndOverride.mockReturnValue(undefined);
    const superCanActivateSpy = jest
      .spyOn(AuthGuard('jwt').prototype, 'canActivate')
      .mockResolvedValue(true);

    const result = await guard.canActivate(mockExecutionContext);

    expect(reflector.getAllAndOverride).toHaveBeenCalled();
    expect(superCanActivateSpy).toHaveBeenCalledWith(mockExecutionContext);
    expect(result).toBe(true);
  });
});

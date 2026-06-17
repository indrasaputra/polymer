import { ExecutionContext } from '@nestjs/common';
import { CurrentUser } from '../dto/current-user.dto';

let capturedFactory: (_data: unknown, ctx: ExecutionContext) => CurrentUser;

jest.mock('@nestjs/common', () => {
  const originalModule = jest.requireActual('@nestjs/common');
  return {
    ...originalModule,
    createParamDecorator: (
      factory: (_data: unknown, ctx: ExecutionContext) => CurrentUser,
    ) => {
      capturedFactory = factory;
      return () => {};
    },
  };
});

// Put import here so Jest has a chance to swap out createParamDecorator with mock.
import { CurrentUserDecorator } from './current-user.decorator';

describe('CurrentUserDecorator', () => {
  it('should have initialized the decorator', () => {
    expect(CurrentUserDecorator).toBeDefined();
    expect(capturedFactory).toBeInstanceOf(Function);
  });

  it('should return the user object from the request', () => {
    const mockUser: CurrentUser = {
      id: 'user-123',
      email: 'test@example.com',
    };

    const mockRequest = {
      user: mockUser,
    };

    const mockExecutionContext = {
      switchToHttp: jest.fn().mockReturnThis(),
      getRequest: jest.fn().mockReturnValue(mockRequest),
    } as unknown as ExecutionContext;

    const result = capturedFactory(null, mockExecutionContext);

    expect(mockExecutionContext.switchToHttp).toHaveBeenCalled();
    expect(mockExecutionContext.switchToHttp().getRequest).toHaveBeenCalled();
    expect(result).toEqual(mockUser);
  });

  it('should return undefined if user is not present on the request', () => {
    const mockExecutionContext = {
      switchToHttp: jest.fn().mockReturnThis(),
      getRequest: jest.fn().mockReturnValue({}),
    } as unknown as ExecutionContext;

    const result = capturedFactory(null, mockExecutionContext);

    expect(result).toBeUndefined();
  });
});

import { AUTH_TYPE_KEY, AuthType, Public } from './auth-type.decorator';
import { SetMetadata } from '@nestjs/common';

jest.mock('@nestjs/common', () => ({
  ...jest.requireActual('@nestjs/common'),
  SetMetadata: jest.fn().mockReturnValue(() => {}),
}));

describe('Public decorator', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should call SetMetadata with correct key and value', () => {
    Public();

    expect(SetMetadata).toHaveBeenCalledWith(AUTH_TYPE_KEY, AuthType.NONE);
  });

  it('should have correct AUTH_TYPE_KEY', () => {
    expect(AUTH_TYPE_KEY).toBe('AuthTypeKey');
  });

  it('should have correct AuthType values', () => {
    expect(AuthType.NONE).toBe('None');
    expect(AuthType.BASIC).toBe('Basic');
    expect(AuthType.JWT).toBe('JWT');
  });
});

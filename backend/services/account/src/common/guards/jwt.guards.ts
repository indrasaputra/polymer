import { ExecutionContext, Injectable } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import { AuthGuard } from '@nestjs/passport';
import { AUTH_TYPE_KEY, AuthType } from '../decorators/auth-type.decorator';

@Injectable()
export class JwtAuthGuard extends AuthGuard('jwt') {
  constructor(private readonly reflector: Reflector) {
    super();
  }

  canActivate(context: ExecutionContext) {
    const authType =
      this.reflector.getAllAndOverride<AuthType>(AUTH_TYPE_KEY, [
        context.getHandler(),
        context.getClass(),
      ]) ?? AuthType.JWT;

    if (authType === AuthType.NONE) {
      return true;
    }

    return super.canActivate(context);
  }
}

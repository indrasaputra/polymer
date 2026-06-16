import { createParamDecorator, ExecutionContext } from '@nestjs/common';
import { Request } from 'express';
import { CurrentUser } from '../dto/current-user.dto';

export const CurrentUserDecorator = createParamDecorator(
  (_data: unknown, ctx: ExecutionContext): CurrentUser => {
    const request = ctx
      .switchToHttp()
      .getRequest<Request & { user: CurrentUser }>();
    return request.user;
  },
);

import {
  CanActivate,
  ExecutionContext,
  Injectable,
  UnauthorizedException,
} from '@nestjs/common';
import { Request } from 'express';
import { Config } from '../../config/config';

@Injectable()
export class WebhookSecretGuard implements CanActivate {
  canActivate(context: ExecutionContext): boolean {
    const request = context.switchToHttp().getRequest<Request>();
    const secret = request.header('x-webhook-secret');

    if (
      !secret ||
      typeof secret !== 'string' ||
      secret !== Config.WEBHOOK_SECRET
    ) {
      throw new UnauthorizedException('Secret is invalid');
    }

    return true;
  }
}

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
  constructor(private readonly config: Config) {}

  canActivate(context: ExecutionContext): boolean {
    const request = context.switchToHttp().getRequest<Request>();
    const secret = request.header('x-webhook-secret');

    if (
      !secret ||
      typeof secret !== 'string' ||
      secret !== this.config.webhookSecret
    ) {
      throw new UnauthorizedException('Secret is invalid');
    }

    return true;
  }
}

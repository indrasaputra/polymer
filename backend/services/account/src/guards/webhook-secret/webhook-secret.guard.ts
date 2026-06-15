import {
  CanActivate,
  ExecutionContext,
  Injectable,
  InternalServerErrorException,
  UnauthorizedException,
} from '@nestjs/common';
import { Request } from 'express';

@Injectable()
export class WebhookSecretGuard implements CanActivate {
  canActivate(context: ExecutionContext): boolean {
    const expectedSecret = process.env.WEBHOOK_SECRET;
    if (!expectedSecret) {
      throw new InternalServerErrorException(
        'WEBHOOK_SECRET is not configured',
      );
    }

    const request = context.switchToHttp().getRequest<Request>();
    const secret = request.header('x-webhook-secret');

    if (!secret || typeof secret !== 'string' || secret !== expectedSecret) {
      throw new UnauthorizedException('Secret is invalid');
    }

    return true;
  }
}

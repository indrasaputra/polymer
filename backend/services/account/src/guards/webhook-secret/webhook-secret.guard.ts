import {
  CanActivate,
  ExecutionContext,
  Injectable,
  UnauthorizedException,
} from '@nestjs/common';
import { Request } from 'express';
import { Observable } from 'rxjs';

@Injectable()
export class WebhookSecretGuard implements CanActivate {
  canActivate(
    context: ExecutionContext,
  ): boolean | Promise<boolean> | Observable<boolean> {
    const request = context.switchToHttp().getRequest<Request>();
    const secret = request.headers['x-webhook-secret'];

    if (secret !== process.env.WEBHOOK_SECRET) {
      throw new UnauthorizedException('Secret is invalid');
    }

    return true;
  }
}

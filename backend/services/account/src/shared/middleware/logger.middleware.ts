import { Injectable, Logger, NestMiddleware } from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';

const SENSITIVE_KEYS = [
  'password',
  'token',
  'secret',
  'authorization',
  'email',
];

@Injectable()
export class LoggerMiddleware implements NestMiddleware {
  private readonly logger = new Logger('HTTP');
  private readonly isDev = (process.env.ENV ?? 'development') === 'development';

  use(req: Request, res: Response, next: NextFunction): void {
    const { method, url } = req;
    const start = Date.now();

    res.on('finish', () => {
      const duration = Date.now() - start;
      const contentLength = String(res.getHeader('content-length') ?? '-');

      this.logger.log(
        `${method} ${url} ${res.statusCode} ${duration}ms ${contentLength}`,
      );

      if (this.isDev && req.body) {
        this.logger.debug(`Body: ${JSON.stringify(req.body)}`);
      } else if (req.body) {
        this.logger.debug(`Body: ${JSON.stringify(this.redact(req.body))}`);
      }
    });

    next();
  }

  private redact(body: Record<string, unknown>): Record<string, unknown> {
    return Object.fromEntries(
      Object.entries(body).map(([key, value]) => [
        key,
        SENSITIVE_KEYS.includes(key.toLowerCase()) ? '[REDACTED]' : value,
      ]),
    );
  }
}

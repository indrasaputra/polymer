import { Injectable, Logger, NestMiddleware } from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';

@Injectable()
export class LoggerMiddleware implements NestMiddleware {
  private readonly logger = new Logger('HTTP');

  use(req: Request, res: Response, next: NextFunction): void {
    const { method, url, body } = req;

    this.logger.log(
      `Incoming request: ${method} ${url} - body: ${JSON.stringify(body)}`,
    );

    res.on('finish', () => {
      this.logger.log(`Response: ${method} ${url} ${res.statusCode}`);
    });

    next();
  }
}

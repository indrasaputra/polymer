import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { HttpExceptionFilter } from './shared/filters/http-exception.filter';
import { Logger, LogLevel } from '@nestjs/common';

async function bootstrap() {
  const logger = new Logger('Account');
  const port = process.env.PORT ?? 9001;

  const isDev = (process.env.ENV ?? 'development') === 'development';

  const logLevels: LogLevel[] = isDev
    ? ['error', 'warn', 'log', 'debug', 'verbose']
    : ['error'];

  const app = await NestFactory.create(AppModule, {
    logger: logLevels,
  });

  app.setGlobalPrefix('api/v1');
  app.useGlobalFilters(new HttpExceptionFilter());

  app.enableShutdownHooks(); // Ensures onModuleDestroy events run when your server closes

  await app.listen(port);
  logger.log(`Accout service is running on: http://localhost:${port}`);
}

bootstrap();

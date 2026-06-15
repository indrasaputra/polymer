import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { HttpExceptionFilter } from './shared/filters/http-exception.filter';
import { Logger } from '@nestjs/common';

async function bootstrap() {
  const logger = new Logger('Account');
  const app = await NestFactory.create(AppModule);
  const port = process.env.PORT ?? 9001;

  app.setGlobalPrefix('api/v1');
  app.useGlobalFilters(new HttpExceptionFilter());

  app.enableShutdownHooks(); // Ensures onModuleDestroy events run when your server closes

  await app.listen(port);
  logger.log(`Accout service is running on: http://localhost:${port}`);
}

bootstrap();

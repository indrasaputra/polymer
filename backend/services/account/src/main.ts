import { NestFactory, Reflector } from '@nestjs/core';
import { ValidationPipe, ClassSerializerInterceptor } from '@nestjs/common';
import { Logger, LoggerErrorInterceptor } from 'nestjs-pino';
import { AppModule } from './app.module';
import { HttpExceptionFilter } from './common/filters/http-exception.filter';
import { logResponseBody } from './common/middleware/logger.middleware';
import { ResponseInterceptor } from './common/interceptors/response.interceptor';
import { Config } from './config/config';

async function bootstrap() {
  const app = await NestFactory.create(AppModule, {
    bufferLogs: true,
    bodyParser: true,
    rawBody: true,
  });
  const logger = app.get(Logger);
  const config = app.get(Config);

  app.use(logResponseBody);
  app.useLogger(logger);
  app.useGlobalInterceptors(new ClassSerializerInterceptor(app.get(Reflector)));
  app.useGlobalInterceptors(new ResponseInterceptor());
  app.useGlobalInterceptors(new LoggerErrorInterceptor());
  app.setGlobalPrefix('api/v1');
  app.useGlobalFilters(new HttpExceptionFilter());
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      transform: true,
    }),
  );

  app.enableShutdownHooks();

  await app.listen(config.port);
  logger.log(`Backend service - account - is running on port: ${config.port}`);
}

bootstrap();

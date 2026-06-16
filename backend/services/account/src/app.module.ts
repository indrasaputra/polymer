import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { SignupModule } from './signup/signup.module';
import { pinoHttpConfig } from './shared/logger/pino.logger';
import { LoggerModule } from 'nestjs-pino';

@Module({
  imports: [
    // 1. Force environment configurations to load globally across your architecture
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    LoggerModule.forRoot({
      pinoHttp: pinoHttpConfig,
      assignResponse: true, // enable propagation of `assign` fields into "request completed" logs
    }),
    // 2. Register signup module execution scope
    SignupModule,
  ],
})
export class AppModule {}

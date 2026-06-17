import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { ProfileModule } from './profile/profile.module';
import { pinoHttpConfig } from './common/logger/pino.logger';
import { LoggerModule } from 'nestjs-pino';
import { APP_GUARD } from '@nestjs/core';
import { JwtAuthGuard } from './common/guards/jwt.guards';
import { PassportModule } from '@nestjs/passport';
import { JwtStrategy } from './common/strategies/jwt.strategy';

@Module({
  imports: [
    // 1. Force environment configurations to load globally across your architecture
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    PassportModule,
    LoggerModule.forRoot({
      pinoHttp: pinoHttpConfig,
      assignResponse: true, // enable propagation of `assign` fields into "request completed" logs
    }),
    // 2. Register profile module execution scope
    ProfileModule,
  ],
  providers: [
    JwtStrategy,
    {
      provide: APP_GUARD,
      useClass: JwtAuthGuard,
    },
  ],
})
export class AppModule {}

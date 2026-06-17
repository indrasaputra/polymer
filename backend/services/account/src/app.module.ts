import { Module } from '@nestjs/common';
import { TypedConfigModule, dotenvLoader } from 'nest-typed-config';
import { ProfileModule } from './profile/profile.module';
import { pinoHttpConfig } from './common/logger/pino.logger';
import { LoggerModule } from 'nestjs-pino';
import { APP_GUARD } from '@nestjs/core';
import { JwtAuthGuard } from './common/guards/jwt.guards';
import { PassportModule } from '@nestjs/passport';
import { JwtStrategy } from './common/strategies/jwt.strategy';
import { Config } from './config/config';

@Module({
  imports: [
    TypedConfigModule.forRoot({
      isGlobal: true,
      schema: Config,
      load: dotenvLoader({
        separator: '__',
        keyTransformer: (key) =>
          key
            .toLowerCase()
            .replace(/(?<!_)_([a-z])/g, (_, p1) => p1.toUpperCase()),
      }),
    }),
    PassportModule,
    LoggerModule.forRootAsync({
      inject: [Config],
      useFactory: (config: Config) => ({
        pinoHttp: pinoHttpConfig(config),
        assignResponse: true,
      }),
    }),
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

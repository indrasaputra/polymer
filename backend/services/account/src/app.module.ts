import { Module, MiddlewareConsumer } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { SignupModule } from './signup/signup.module';
import { LoggerMiddleware } from './shared/middleware/logger.middleware';

@Module({
  imports: [
    // 1. Force environment configurations to load globally across your architecture
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    // 2. Register signup module execution scope
    SignupModule,
  ],
})
export class AppModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(LoggerMiddleware).forRoutes('*');
  }
}

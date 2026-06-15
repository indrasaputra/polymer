import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { SignupModule } from './signup/signup.module';

@Module({
  imports: [
    // 1. Force environment configurations to load globally across your architecture
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    // 2. Register your transactional wallets module execution scope
    SignupModule,
  ],
})
export class AppModule {}

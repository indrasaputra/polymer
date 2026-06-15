import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { SignupController } from './signup/signup.controller';
import { SignupService } from './signup/signup.service';

@Module({
  imports: [],
  controllers: [AppController, SignupController],
  providers: [AppService, SignupService],
})
export class AppModule {}

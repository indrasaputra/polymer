import { Module } from '@nestjs/common';
import { SignupController } from './signup.controller';
import { SignupService } from './signup.service';
import { PrismaModule } from '../prisma/prisma.module'; // Ensure your DB connection module is here

@Module({
  imports: [PrismaModule],
  controllers: [SignupController],
  providers: [SignupService],
})
export class SignupModule {}

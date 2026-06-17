import { Module } from '@nestjs/common';
import { ProfileController } from './profile.controller';
import { ProfileService } from './profile.service';
import { PrismaModule } from '../prisma/prisma.module'; // Ensure your DB connection module is here

@Module({
  imports: [PrismaModule],
  controllers: [ProfileController],
  providers: [ProfileService],
})
export class ProfileModule {}

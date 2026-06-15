import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { PrismaService } from './prisma.service';

@Module({
  imports: [ConfigModule], // Grants PrismaService access to configService
  providers: [PrismaService],
  exports: [PrismaService],
})
export class PrismaModule {}

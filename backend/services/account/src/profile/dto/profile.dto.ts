import { IsUUID, IsEmail, IsOptional } from 'class-validator';
import { Transform, Type } from 'class-transformer';

export class ProfileWebhookDto {
  @IsUUID()
  id: string;

  @IsEmail()
  email: string;

  @IsOptional()
  @Transform(({ obj }) => obj.created_at ?? obj.createdAt)
  @Type(() => Date)
  createdAt?: Date;
}

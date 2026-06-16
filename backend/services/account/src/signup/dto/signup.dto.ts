import { IsUUID, IsEmail, IsOptional } from 'class-validator';
import { Transform } from 'class-transformer';

export class SignupWebhookDto {
  @IsUUID()
  id: string;

  @IsEmail()
  email: string;

  @IsOptional()
  @Transform(({ obj }) => obj.created_at ?? obj.createdAt)
  createdAt?: Date;
}

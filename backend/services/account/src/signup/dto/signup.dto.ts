import { IsUUID, IsEmail } from 'class-validator';

export class SignupWebhookDto {
  @IsUUID()
  id: string;

  @IsEmail()
  email: string;

  createdAt?: Date;
}

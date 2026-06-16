import { Profile as ORMProfile } from '../../../generated/prisma/client';
import { IsUUID, IsEmail, IsOptional } from 'class-validator';
import { Transform, Type, Expose } from 'class-transformer';

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

export class ProfileResponseDto {
  id: string;
  email: string;

  @Expose({ name: 'first_name' })
  firstName: string | null;

  @Expose({ name: 'last_name' })
  lastName: string | null;

  @Expose({ name: 'created_at' })
  @Transform(({ value }) => value.toISOString().split('T')[0]) // Formats JSON output string
  createdAt: Date;

  constructor(obj: ProfileResponseDto) {
    Object.assign(this, obj);
  }

  static fromOrm(orm: ORMProfile): ProfileResponseDto {
    return new ProfileResponseDto({
      id: orm.id,
      email: orm.email,
      firstName: orm.firstName,
      lastName: orm.lastName,
      createdAt: orm.createdAt,
    });
  }
}

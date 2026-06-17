import {
  IsString,
  IsIn,
  IsNumber,
  IsDefined,
  ValidateNested,
} from 'class-validator';
import { Type } from 'class-transformer';

export class Postgre {
  @IsString()
  @IsDefined()
  readonly url: string;

  @IsString()
  @IsDefined()
  readonly schema: string = 'public';
}

export class Supabase {
  @IsString()
  @IsDefined()
  readonly jwksUrl: string;
}

export class OpenTelemetry {
  @IsString()
  @IsDefined()
  readonly exporterEndpoint: string;
}

export class Config {
  @IsIn(['development', 'production', 'test'])
  @IsDefined()
  readonly env: string = 'development';

  get isDevelopment(): boolean {
    return this.env === 'development';
  }

  get isStaging(): boolean {
    return this.env === 'staging';
  }

  get isProduction(): boolean {
    return this.env === 'production';
  }

  @IsString()
  @IsDefined()
  readonly serviceName: string = 'account';

  @IsNumber()
  @IsDefined()
  @Type(() => Number)
  readonly port: number = 9001;

  @IsString()
  @IsDefined()
  readonly webhookSecret: string;

  @Type(() => Postgre)
  @IsDefined()
  @ValidateNested()
  readonly postgre: Postgre;

  @Type(() => Supabase)
  @IsDefined()
  @ValidateNested()
  readonly supabase: Supabase;

  @Type(() => OpenTelemetry)
  @IsDefined()
  @ValidateNested()
  readonly otel: OpenTelemetry;
}

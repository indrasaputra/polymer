import {
  IsString,
  IsIn,
  IsNumber,
  IsBoolean,
  IsDefined,
  ValidateNested,
} from 'class-validator';
import { Type, Transform } from 'class-transformer';

function parseBoolean(value: unknown, defaultValue: boolean): boolean {
  if (value === undefined || value === null || value === '') {
    return defaultValue;
  }
  if (typeof value === 'boolean') {
    return value;
  }
  if (typeof value === 'string') {
    const v = value.trim().toLowerCase();
    if (v === 'true' || v === '1' || v === 'yes' || v === 'y') return true;
    if (v === 'false' || v === '0' || v === 'no' || v === 'n') return false;
    return defaultValue;
  }
  if (typeof value === 'number') {
    return value !== 0;
  }
  return defaultValue;
}

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
  @Transform(({ value }) => parseBoolean(value, false))
  @IsBoolean()
  @IsDefined()
  readonly enabled: boolean = false;

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

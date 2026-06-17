import * as dotenv from 'dotenv';
dotenv.config();

export const ENV = process.env.ENV ?? 'development';
export const isDevelopment = ENV === 'development';
export const isStaging = ENV === 'staging';
export const isProduction = ENV === 'production';

function requireEnv(key: string): string {
  const value = process.env[key];
  if (!value) {
    throw new Error(`Missing required environment variable: ${key}`);
  }
  return value;
}

function optionalEnv(key: string, defaultValue: any): string {
  return process.env[key] ?? defaultValue;
}

export const Config = {
  // Service
  ENV: optionalEnv('ENV', 'development'),
  PORT: optionalEnv('PORT', 9001),
  SERVICE_NAME: requireEnv('SERVICE_NAME'),

  // Webhook
  WEBHOOK_SECRET: requireEnv('WEBHOOK_SECRET'),

  // Database
  DATABASE_URL: requireEnv('DATABASE_URL'),
  DATABASE_SCHEMA: optionalEnv('DATABASE_SCHEMA', 'public'),

  // Supabase
  SUPABASE_JWKS_URL: requireEnv('SUPABASE_JWKS_URL'),

  // OpenTelemetry
  OTEL_ENABLED: optionalEnv('OTEL_ENABLED', false),
  OTEL_EXPORTER_ENDPOINT: requireEnv('OTEL_EXPORTER_ENDPOINT'),
} as const;

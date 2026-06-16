export const ENV = process.env.ENV ?? 'development';
export const isDevelopment = ENV === 'development';
export const isStaging = ENV === 'staging';
export const isProduction = ENV === 'production';

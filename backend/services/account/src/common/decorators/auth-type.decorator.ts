import { SetMetadata } from '@nestjs/common';

export const AUTH_TYPE_KEY = 'AuthTypeKey';

export enum AuthType {
  NONE = 'None',
  BASIC = 'Basic',
  JWT = 'JWT',
}

/**
 * Decorator to mark a route as public (no authentication required).
 * @example
 * ＠Get('public-endpoint')
 * ＠Public()               <-- use like this
 * async publicEndpoint() {
 *   // ...endpoint implementation...
 * }
 */
export const Public = () => {
  return SetMetadata(AUTH_TYPE_KEY, AuthType.NONE);
};

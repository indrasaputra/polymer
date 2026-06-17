import { Injectable, UnauthorizedException } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';
import { CurrentUser } from '../dto/current-user.dto';
import { JwtToken } from '../dto/jwt-token.interface';
import { passportJwtSecret } from 'jwks-rsa';
import { Config } from '../../config/config';

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  constructor(config: Config) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: false,
      secretOrKeyProvider: passportJwtSecret({
        cache: true,
        rateLimit: true,
        jwksRequestsPerMinute: 5,
        jwksUri: config.supabase.jwksUrl,
      }),
      algorithms: ['ES256'],
    });
  }

  validate(payload: JwtToken): CurrentUser {
    if (!payload.sub || !payload.email) {
      throw new UnauthorizedException('Invalid token payload');
    }

    return {
      id: payload.sub,
      email: payload.email,
    };
  }
}

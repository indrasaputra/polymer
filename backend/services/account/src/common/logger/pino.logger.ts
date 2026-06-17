// import pino, { transport } from 'pino';
import pino from 'pino';
import { Options } from 'pino-http';
import { uuidv7 } from '@kripod/uuidv7';
import { IncomingMessage } from 'http';
import { Config } from '../../config/config';

export function pinoHttpConfig(config: Config): Options {
  const pinoLogger = pino({
    name: config.serviceName,
    level: config.isDevelopment ? 'debug' : 'info',
    formatters: {
      level: (label: string) => {
        return { severity: label.toUpperCase() }; // renamed to "severity"
      },
    },
    transport: config.isDevelopment ? { target: 'pino-pretty' } : undefined,
    redact: config.isDevelopment
      ? []
      : [
          'req.headers.authorization',
          'req.body.password',
          'req.body.refreshToken',
          'req.body.name',
          'req.body.email',
          // common sso token
          'req.body.identityToken',
          'req.body.accessToken',
        ],
  });

  return {
    logger: pinoLogger,
    genReqId: function (req, res) {
      const existingID = req.id ?? req.headers['x-request-id'];
      if (existingID) return existingID;
      const id = uuidv7();
      req.id = id;
      res.setHeader('X-Request-Id', id);
      return id;
    },
    // add user.id in log
    customProps: (
      req: IncomingMessage & { user?: { id?: string; email?: string } },
    ) => ({
      context: 'HTTP',
      user: req.user
        ? {
            id: req.user.id,
          }
        : undefined,
    }),
    serializers: {
      req(req) {
        return {
          ...req,
          // by default, body is omitted from logs
          // this is just for example in local dev or low throughput system
          body: req.raw.body,
        };
      },
      res(res) {
        return {
          ...res,
          // body is only present if middleware attached it (error responses)
          body: res.raw.body,
        };
      },
    },
  };
}

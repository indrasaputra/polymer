import { Logger } from '@nestjs/common';
import { LoggerMiddleware } from './logger.middleware';
import { Request, Response, NextFunction } from 'express';

const mockRequest = (overrides?: Partial<Request>): Partial<Request> => ({
  method: 'POST',
  url: '/api/v1/signup/webhook',
  body: { id: 'uuid', email: 'john.doe@example.com' },
  ...overrides,
});

const mockResponse = (): Partial<Response> & { getHeader: jest.Mock } => ({
  statusCode: 200,
  getHeader: jest.fn().mockReturnValue('123'),
  on: jest.fn().mockImplementation((event: string, callback: () => void) => {
    if (event === 'finish') callback();
  }),
});

const mockNext: NextFunction = jest.fn();

describe('LoggerMiddleware', () => {
  let middleware: LoggerMiddleware;
  let logSpy: jest.SpyInstance;
  let debugSpy: jest.SpyInstance;

  beforeEach(() => {
    process.env.ENV = 'development';
    middleware = new LoggerMiddleware();
    logSpy = jest.spyOn(Logger.prototype, 'log').mockImplementation(() => {});
    debugSpy = jest
      .spyOn(Logger.prototype, 'debug')
      .mockImplementation(() => {});
  });

  afterEach(() => {
    jest.clearAllMocks();
    delete process.env.ENV;
  });

  it('should be defined', () => {
    expect(middleware).toBeDefined();
  });

  it('should call next()', () => {
    const next = jest.fn();
    middleware.use(
      mockRequest() as Request,
      mockResponse() as unknown as Response,
      next,
    );

    expect(next).toHaveBeenCalled();
  });

  describe('request logging', () => {
    it('should log method, url, status, duration, content-length on finish', () => {
      middleware.use(
        mockRequest() as Request,
        mockResponse() as unknown as Response,
        mockNext,
      );

      expect(logSpy).toHaveBeenCalledWith(
        expect.stringMatching(
          /^POST \/api\/v1\/signup\/webhook 200 \d+ms 123$/,
        ),
      );
    });

    it('should log - for content-length when header is missing', () => {
      const res = mockResponse();
      res.getHeader.mockReturnValue(undefined);

      middleware.use(
        mockRequest() as Request,
        res as unknown as Response,
        mockNext,
      );

      expect(logSpy).toHaveBeenCalledWith(
        expect.stringMatching(/^POST \/api\/v1\/signup\/webhook 200 \d+ms -$/),
      );
    });
  });

  describe('body logging', () => {
    it('should log full body in development without redaction', () => {
      middleware.use(
        mockRequest() as Request,
        mockResponse() as unknown as Response,
        mockNext,
      );

      expect(debugSpy).toHaveBeenCalledWith(
        'Body: {"id":"uuid","email":"john.doe@example.com"}',
      );
    });

    it('should not redact sensitive keys in development', () => {
      const req = mockRequest({
        body: {
          id: 'uuid',
          email: 'john@example.com',
          password: 'supersecret',
          token: 'abc123',
        },
      });

      middleware.use(
        req as Request,
        mockResponse() as unknown as Response,
        mockNext,
      );

      expect(debugSpy).toHaveBeenCalledWith(
        'Body: {"id":"uuid","email":"john@example.com","password":"supersecret","token":"abc123"}',
      );
    });

    it('should redact sensitive keys outside development', () => {
      process.env.ENV = 'production';
      middleware = new LoggerMiddleware();

      const req = mockRequest({
        body: {
          id: 'uuid',
          email: 'john@example.com',
          password: 'supersecret',
          token: 'abc123',
        },
      });

      middleware.use(
        req as Request,
        mockResponse() as unknown as Response,
        mockNext,
      );

      expect(debugSpy).toHaveBeenCalledWith(
        'Body: {"id":"uuid","email":"[REDACTED]","password":"[REDACTED]","token":"[REDACTED]"}',
      );
    });

    it('should not log body when body is empty', () => {
      const req = mockRequest({ body: null });

      middleware.use(
        req as Request,
        mockResponse() as unknown as Response,
        mockNext,
      );

      expect(debugSpy).not.toHaveBeenCalled();
    });

    it('should log body in production with non-sensitive keys unredacted', () => {
      process.env.ENV = 'production';
      middleware = new LoggerMiddleware();

      const req = mockRequest({
        body: { id: 'uuid', name: 'John' },
      });

      middleware.use(
        req as Request,
        mockResponse() as unknown as Response,
        mockNext,
      );

      expect(debugSpy).toHaveBeenCalledWith(
        'Body: {"id":"uuid","name":"John"}',
      );
    });
  });
});

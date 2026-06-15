import { Logger } from '@nestjs/common';
import { LoggerMiddleware } from './logger.middleware';
import { Request, Response, NextFunction } from 'express';

describe('LoggerMiddleware', () => {
  let middleware: LoggerMiddleware;
  let loggerSpy: jest.SpyInstance;

  beforeEach(() => {
    middleware = new LoggerMiddleware();
    loggerSpy = jest
      .spyOn(Logger.prototype, 'log')
      .mockImplementation(() => {});
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  const mockRequest = (overrides?: Partial<Request>): Partial<Request> => ({
    method: 'POST',
    url: '/api/v1/signup/webhook',
    body: { id: 'uuid', email: 'john.doe@example.com' },
    ...overrides,
  });

  const mockResponse = (): Partial<Response> => {
    const res: Partial<Response> = {
      statusCode: 200,
      on: jest
        .fn()
        .mockImplementation((event: string, callback: () => void) => {
          if (event === 'finish') callback();
          return res;
        }),
    };
    return res;
  };

  const mockNext: NextFunction = jest.fn();

  it('should be defined', () => {
    expect(middleware).toBeDefined();
  });

  it('should log incoming request', () => {
    const req = mockRequest();
    const res = mockResponse();

    middleware.use(req as Request, res as Response, mockNext);

    expect(loggerSpy).toHaveBeenCalledWith(
      `Incoming request: POST /api/v1/signup/webhook - body: ${JSON.stringify(req.body)}`,
    );
  });

  it('should log response', () => {
    const req = mockRequest();
    const res = mockResponse();

    middleware.use(req as Request, res as Response, mockNext);

    expect(loggerSpy).toHaveBeenCalledWith(
      `Response: POST /api/v1/signup/webhook 200`,
    );
  });

  it('should call next()', () => {
    const req = mockRequest();
    const res = mockResponse();
    const next = jest.fn();

    middleware.use(req as Request, res as Response, next);

    expect(next).toHaveBeenCalled();
  });

  it('should log correct status code', () => {
    const req = mockRequest();
    const res = mockResponse();
    res.statusCode = 500;

    middleware.use(req as Request, res as Response, mockNext);

    expect(loggerSpy).toHaveBeenCalledWith(
      `Response: POST /api/v1/signup/webhook 500`,
    );
  });
});

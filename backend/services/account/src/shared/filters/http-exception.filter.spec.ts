import {
  ArgumentsHost,
  HttpException,
  HttpStatus,
  Logger,
} from '@nestjs/common';
import { HttpExceptionFilter } from './http-exception.filter';

const mockRequest = (overrides?: object) => ({
  method: 'POST',
  url: '/api/v1/profiles/webhook',
  ...overrides,
});

const mockResponse = () => {
  const res = {
    status: jest.fn().mockReturnThis(),
    json: jest.fn().mockReturnThis(),
  };
  return res;
};

const mockHost = (req: object, res: object): ArgumentsHost =>
  ({
    switchToHttp: () => ({
      getRequest: () => req,
      getResponse: () => res,
    }),
  }) as unknown as ArgumentsHost;

describe('HttpExceptionFilter', () => {
  let filter: HttpExceptionFilter;
  let loggerSpy: jest.SpyInstance;

  beforeEach(() => {
    filter = new HttpExceptionFilter();
    loggerSpy = jest
      .spyOn(Logger.prototype, 'error')
      .mockImplementation(() => {});
    jest.useFakeTimers();
    jest.setSystemTime(new Date('2026-01-01T00:00:00.000Z'));
  });

  afterEach(() => {
    jest.clearAllMocks();
    jest.useRealTimers();
  });

  it('should be defined', () => {
    expect(filter).toBeDefined();
  });

  describe('when exception is HttpException', () => {
    it('should return correct status code', () => {
      const req = mockRequest();
      const res = mockResponse();
      const exception = new HttpException(
        'Unauthorized',
        HttpStatus.UNAUTHORIZED,
      );

      filter.catch(exception, mockHost(req, res));

      expect(res.status).toHaveBeenCalledWith(HttpStatus.UNAUTHORIZED);
    });

    it('should return correct message', () => {
      const req = mockRequest();
      const res = mockResponse();
      const exception = new HttpException(
        'Unauthorized',
        HttpStatus.UNAUTHORIZED,
      );

      filter.catch(exception, mockHost(req, res));

      expect(res.json).toHaveBeenCalledWith({
        statusCode: HttpStatus.UNAUTHORIZED,
        message: 'Unauthorized',
        path: '/api/v1/profiles/webhook',
        timestamp: '2026-01-01T00:00:00.000Z',
      });
    });
  });

  describe('when exception is unknown', () => {
    it('should return 500 status code', () => {
      const req = mockRequest();
      const res = mockResponse();

      filter.catch(new Error('something went wrong'), mockHost(req, res));

      expect(res.status).toHaveBeenCalledWith(HttpStatus.INTERNAL_SERVER_ERROR);
    });

    it('should return internal server error message', () => {
      const req = mockRequest();
      const res = mockResponse();

      filter.catch(new Error('something went wrong'), mockHost(req, res));

      expect(res.json).toHaveBeenCalledWith({
        statusCode: HttpStatus.INTERNAL_SERVER_ERROR,
        message: 'Internal server error',
        path: '/api/v1/profiles/webhook',
        timestamp: '2026-01-01T00:00:00.000Z',
      });
    });

    it('should handle non-Error unknown exceptions', () => {
      const req = mockRequest();
      const res = mockResponse();

      filter.catch({ something: 'weird' }, mockHost(req, res));

      expect(res.status).toHaveBeenCalledWith(HttpStatus.INTERNAL_SERVER_ERROR);
    });
  });

  describe('logging', () => {
    it('should log with stack trace for Error exceptions', () => {
      const req = mockRequest();
      const res = mockResponse();
      const exception = new Error('something went wrong');

      filter.catch(exception, mockHost(req, res));

      expect(loggerSpy).toHaveBeenCalledWith(
        `POST /api/v1/profiles/webhook 500 - Internal server error`,
        exception.stack,
      );
    });

    it('should log stringified exception for non-Error exceptions', () => {
      const req = mockRequest();
      const res = mockResponse();
      const exception = { something: 'weird' };

      filter.catch(exception, mockHost(req, res));

      expect(loggerSpy).toHaveBeenCalledWith(
        `POST /api/v1/profiles/webhook 500 - Internal server error`,
        JSON.stringify(exception),
      );
    });
  });
});

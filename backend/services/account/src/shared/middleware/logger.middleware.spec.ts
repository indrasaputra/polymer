// logResponseBody.middleware.spec.ts
import { Request, Response, NextFunction } from 'express';
import { logResponseBody } from './logger.middleware';

describe('logResponseBody', () => {
  let mockReq: Partial<Request>;
  let mockRes: Partial<Response> & { body?: any };
  let mockNext: NextFunction;
  let originalJson: jest.Mock;

  beforeEach(() => {
    originalJson = jest.fn().mockReturnValue(mockRes);
    mockReq = {};
    mockRes = {
      json: originalJson,
      statusCode: 200,
    };
    mockNext = jest.fn();
  });

  it('should call next()', () => {
    logResponseBody(mockReq as Request, mockRes as Response, mockNext);
    expect(mockNext).toHaveBeenCalled();
  });

  it('should intercept res.json()', () => {
    logResponseBody(mockReq as Request, mockRes as Response, mockNext);
    expect(mockRes.json).not.toBe(originalJson); // replaced
  });

  it('should call original res.json() with the body', () => {
    logResponseBody(mockReq as Request, mockRes as Response, mockNext);

    const body = { message: 'ok' };
    mockRes.json!(body);

    expect(originalJson).toHaveBeenCalledWith(body);
  });

  describe('when status code is 2xx', () => {
    it('should NOT attach body to res', () => {
      mockRes.statusCode = 200;
      logResponseBody(mockReq as Request, mockRes as Response, mockNext);

      mockRes.json!({ data: 'success' });

      expect(mockRes.body).toBeUndefined();
    });
  });

  describe('when status code is 4xx', () => {
    it('should attach body to res', () => {
      mockRes.statusCode = 400;
      logResponseBody(mockReq as Request, mockRes as Response, mockNext);

      const body = { error: 'Bad Request' };
      mockRes.json!(body);

      expect(mockRes.body).toEqual(body);
    });

    it('should attach body for 404', () => {
      mockRes.statusCode = 404;
      logResponseBody(mockReq as Request, mockRes as Response, mockNext);

      const body = { error: 'Not Found' };
      mockRes.json!(body);

      expect(mockRes.body).toEqual(body);
    });
  });

  describe('when status code is 5xx', () => {
    it('should attach body to res', () => {
      mockRes.statusCode = 500;
      logResponseBody(mockReq as Request, mockRes as Response, mockNext);

      const body = { error: 'Internal Server Error' };
      mockRes.json!(body);

      expect(mockRes.body).toEqual(body);
    });
  });

  it('should return the result of original res.json()', () => {
    const returnValue = { mocked: true };
    originalJson.mockReturnValue(returnValue);

    logResponseBody(mockReq as Request, mockRes as Response, mockNext);
    const result = mockRes.json!({ anything: true });

    expect(result).toEqual(returnValue);
  });
});

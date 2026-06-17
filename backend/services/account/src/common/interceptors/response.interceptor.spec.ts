import { ExecutionContext, HttpStatus } from '@nestjs/common';
import { CallHandler } from '@nestjs/common';
import { of } from 'rxjs';
import { Response } from 'express';
import { ResponseInterceptor } from './response.interceptor';

const mockCallHandler = (data: unknown): CallHandler => ({
  handle: jest.fn().mockReturnValue(of(data)),
});

const mockExecutionContext = (statusCode: number): ExecutionContext =>
  ({
    switchToHttp: () => ({
      getResponse: (): Partial<Response> => ({ statusCode }),
    }),
  }) as unknown as ExecutionContext;

describe('ResponseInterceptor', () => {
  let interceptor: ResponseInterceptor<unknown>;

  beforeEach(() => {
    interceptor = new ResponseInterceptor();
  });

  it('should be defined', () => {
    expect(interceptor).toBeDefined();
  });

  it('should wrap response in data envelope', (done) => {
    const ctx = mockExecutionContext(HttpStatus.OK);
    const handler = mockCallHandler({ id: 'uuid', email: 'john@example.com' });

    interceptor.intercept(ctx, handler).subscribe((result) => {
      expect(result).toEqual({
        data: { id: 'uuid', email: 'john@example.com' },
      });
      done();
    });
  });

  it('should return undefined for 204 No Content', (done) => {
    const ctx = mockExecutionContext(HttpStatus.NO_CONTENT);
    const handler = mockCallHandler(undefined);

    interceptor.intercept(ctx, handler).subscribe((result) => {
      expect(result).toBeUndefined();
      done();
    });
  });

  it('should wrap null data in envelope', (done) => {
    const ctx = mockExecutionContext(HttpStatus.OK);
    const handler = mockCallHandler(null);

    interceptor.intercept(ctx, handler).subscribe((result) => {
      expect(result).toEqual({ data: null });
      done();
    });
  });

  it('should wrap array data in envelope', (done) => {
    const ctx = mockExecutionContext(HttpStatus.OK);
    const handler = mockCallHandler([{ id: 'uuid1' }, { id: 'uuid2' }]);

    interceptor.intercept(ctx, handler).subscribe((result) => {
      expect(result).toEqual({
        data: [{ id: 'uuid1' }, { id: 'uuid2' }],
      });
      done();
    });
  });
});

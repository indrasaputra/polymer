import { Request, Response, NextFunction } from 'express';

export function logResponseBody(_: Request, res: Response, next: NextFunction) {
  const originalJson = res.json.bind(res);

  // intercept res.json()
  res.json = (body) => {
    // only attach body to res if it's an error response
    if (res.statusCode >= 400) {
      (res as any).body = body;
    }
    return originalJson(body);
  };

  next();
}

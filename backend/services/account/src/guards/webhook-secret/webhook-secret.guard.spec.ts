import { WebhookSecretGuard } from './webhook-secret.guard';

describe('WebhookSecretGuard', () => {
  it('should be defined', () => {
    expect(new WebhookSecretGuard()).toBeDefined();
  });
});

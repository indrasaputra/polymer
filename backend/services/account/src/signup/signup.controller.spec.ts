import { Test, TestingModule } from '@nestjs/testing';
import { SignupController } from './signup.controller';
import { SignupService } from './signup.service';
import { WebhookSecretGuard } from '../guards/webhook-secret/webhook-secret.guard';
import { SignupWebhookDto } from './dto/signup.dto';

const mockSignupService = {
  handleWebhook: jest.fn().mockResolvedValue(undefined),
};

const mockWebhookSecretGuard = {
  canActivate: jest.fn().mockReturnValue(true),
};

describe('SignupController', () => {
  let controller: SignupController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [SignupController],
      providers: [{ provide: SignupService, useValue: mockSignupService }],
    })
      .overrideGuard(WebhookSecretGuard)
      .useValue(mockWebhookSecretGuard)
      .compile();

    controller = module.get<SignupController>(SignupController);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });

  describe('handleWebhook', () => {
    const payload: SignupWebhookDto = {
      id: 'uuid-123',
      email: 'john.doe@example.com',
    };

    it('should call signupService.handleWebhook with payload', async () => {
      await controller.handleWebhook(payload);

      expect(mockSignupService.handleWebhook).toHaveBeenCalledWith(payload);
    });

    it('should return void', async () => {
      const result = await controller.handleWebhook(payload);

      expect(result).toBeUndefined();
    });

    it('should throw if signupService throws', async () => {
      mockSignupService.handleWebhook.mockRejectedValueOnce(
        new Error('Service error'),
      );

      await expect(controller.handleWebhook(payload)).rejects.toThrow(
        'Service error',
      );
    });
  });
});

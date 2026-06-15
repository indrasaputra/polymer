import { Test, TestingModule } from '@nestjs/testing';
import { SignupController } from './signup.controller';
import { SignupService } from './signup.service';
import { WebhookSecretGuard } from '../guards/webhook-secret/webhook-secret.guard';
import { SignupWebhookDto } from './dto/signup.dto';
import { ValidationPipe } from '@nestjs/common';

const mockSignupService = {
  handleWebhook: jest.fn().mockResolvedValue(undefined),
};

const mockWebhookSecretGuard = {
  canActivate: jest.fn().mockReturnValue(true),
};

describe('SignupController', () => {
  let controller: SignupController;
  let validationPipe: ValidationPipe;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [SignupController],
      providers: [{ provide: SignupService, useValue: mockSignupService }],
    })
      .overrideGuard(WebhookSecretGuard)
      .useValue(mockWebhookSecretGuard)
      .compile();

    controller = module.get<SignupController>(SignupController);
    validationPipe = new ValidationPipe({
      whitelist: true,
      forbidNonWhitelisted: true,
      transform: true,
    });
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });

  describe('handleWebhook', () => {
    const payload: SignupWebhookDto = {
      id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
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

    it('should reject invalid UUID', async () => {
      await expect(
        validationPipe.transform(
          { id: 'not-a-uuid', email: 'john.doe@example.com' },
          { type: 'body', metatype: SignupWebhookDto },
        ),
      ).rejects.toThrow();
    });

    it('should reject invalid email', async () => {
      await expect(
        validationPipe.transform(
          { id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', email: 'not-an-email' },
          { type: 'body', metatype: SignupWebhookDto },
        ),
      ).rejects.toThrow();
    });

    it('should reject unknown fields', async () => {
      await expect(
        validationPipe.transform(
          {
            id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
            email: 'john.doe@example.com',
            unknown: 'field',
          },
          { type: 'body', metatype: SignupWebhookDto },
        ),
      ).rejects.toThrow();
    });
  });
});

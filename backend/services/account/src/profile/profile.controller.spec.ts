import { Test, TestingModule } from '@nestjs/testing';
import { ProfileController } from './profile.controller';
import { ProfileService } from './profile.service';
import { WebhookSecretGuard } from '../guards/webhook-secret/webhook-secret.guard';
import { ProfileWebhookDto } from './dto/profile.dto';
import { ValidationPipe } from '@nestjs/common';

const mockProfileService = {
  handleCreateProfileWebhook: jest.fn().mockResolvedValue(undefined),
};

const mockWebhookSecretGuard = {
  canActivate: jest.fn().mockReturnValue(true),
};

describe('ProfileController', () => {
  let controller: ProfileController;
  let validationPipe: ValidationPipe;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [ProfileController],
      providers: [{ provide: ProfileService, useValue: mockProfileService }],
    })
      .overrideGuard(WebhookSecretGuard)
      .useValue(mockWebhookSecretGuard)
      .compile();

    controller = module.get<ProfileController>(ProfileController);
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

  describe('createProfileWebhook', () => {
    const payload: ProfileWebhookDto = {
      id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
      email: 'john.doe@example.com',
    };

    it('should call ProfileService.createProfileWebhook with payload', async () => {
      await controller.createProfileWebhook(payload);
      expect(
        mockProfileService.handleCreateProfileWebhook,
      ).toHaveBeenCalledWith(payload);
    });

    it('should return void', async () => {
      const result = await controller.createProfileWebhook(payload);

      expect(result).toBeUndefined();
    });

    it('should throw if ProfileService throws', async () => {
      mockProfileService.handleCreateProfileWebhook.mockRejectedValueOnce(
        new Error('Service error'),
      );

      await expect(controller.createProfileWebhook(payload)).rejects.toThrow(
        'Service error',
      );
    });

    it('should reject invalid UUID', async () => {
      await expect(
        validationPipe.transform(
          { id: 'not-a-uuid', email: 'john.doe@example.com' },
          { type: 'body', metatype: ProfileWebhookDto },
        ),
      ).rejects.toThrow();
    });

    it('should reject invalid email', async () => {
      await expect(
        validationPipe.transform(
          { id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', email: 'not-an-email' },
          { type: 'body', metatype: ProfileWebhookDto },
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
          { type: 'body', metatype: ProfileWebhookDto },
        ),
      ).rejects.toThrow();
    });
  });
});

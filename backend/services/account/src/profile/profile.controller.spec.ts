import { Test, TestingModule } from '@nestjs/testing';
import { ValidationPipe, NotFoundException } from '@nestjs/common';
import { ProfileController } from './profile.controller';
import { ProfileService } from './profile.service';
import { WebhookSecretGuard } from '../common/guards/webhook-secret.guard';
import { ProfileWebhookDto, ProfileResponseDto } from './dto/profile.dto';
import { CurrentUser } from '../common/dto/current-user.dto';
import { JwtAuthGuard } from '../common/guards/jwt.guards';

const mockUser: CurrentUser = {
  id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
  email: 'john.doe@example.com',
};

const mockProfile = new ProfileResponseDto({
  id: mockUser.id,
  email: mockUser.email,
  firstName: 'John',
  lastName: 'Doe',
  createdAt: new Date('2026-01-01T00:00:00.000Z'),
});

const mockProfileService = {
  handleCreateProfileWebhook: jest.fn().mockResolvedValue(undefined),
  findOne: jest.fn(),
};

const mockJwtAuthGuard = {
  canActivate: jest.fn().mockReturnValue(true),
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
      .overrideGuard(JwtAuthGuard)
      .useValue(mockJwtAuthGuard)
      .compile();

    controller = module.get<ProfileController>(ProfileController);
    validationPipe = new ValidationPipe({
      whitelist: true,
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
  });

  describe('getProfile', () => {
    it('should return profile when user exists', async () => {
      mockProfileService.findOne.mockResolvedValueOnce(mockProfile);

      const result = await controller.getProfile(mockUser);

      expect(result).toEqual(mockProfile);
    });

    it('should call profileService.findOne with user id', async () => {
      mockProfileService.findOne.mockResolvedValueOnce(mockProfile);

      await controller.getProfile(mockUser);

      expect(mockProfileService.findOne).toHaveBeenCalledWith(mockUser.id);
    });

    it('should throw NotFoundException when profile does not exist', async () => {
      mockProfileService.findOne.mockResolvedValueOnce(null);

      await expect(controller.getProfile(mockUser)).rejects.toThrow(
        NotFoundException,
      );
    });

    it('should throw NotFoundException with correct message', async () => {
      mockProfileService.findOne.mockResolvedValueOnce(null);

      await expect(controller.getProfile(mockUser)).rejects.toThrow(
        `User with id ${mockUser.id} not found`,
      );
    });

    it('should throw if profileService throws', async () => {
      mockProfileService.findOne.mockRejectedValueOnce(new Error('DB error'));

      await expect(controller.getProfile(mockUser)).rejects.toThrow('DB error');
    });
  });
});

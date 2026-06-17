import { Test, TestingModule } from '@nestjs/testing';
import { ProfileService } from './profile.service';
import { PrismaService } from '../prisma/prisma.service';
import { ProfileWebhookDto, ProfileResponseDto } from './dto/profile.dto';

const mockProfile = {
  id: 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
  email: 'john.doe@example.com',
  firstName: 'John',
  lastName: 'Doe',
  createdAt: new Date('2026-01-01T00:00:00.000Z'),
  updatedAt: new Date('2026-01-01T00:00:00.000Z'),
  deletedAt: null,
};

const mockPrismaService = {
  profile: {
    upsert: jest.fn().mockResolvedValue(undefined),
    findUnique: jest.fn(),
  },
};

describe('ProfileService', () => {
  let service: ProfileService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ProfileService,
        { provide: PrismaService, useValue: mockPrismaService },
      ],
    }).compile();

    service = module.get<ProfileService>(ProfileService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('handleCreateProfileWebhook', () => {
    it('should create profile with extracted name from email', async () => {
      const payload: ProfileWebhookDto = {
        id: 'uuid-123',
        email: 'john.doe@example.com',
      };

      await service.handleCreateProfileWebhook(payload);

      expect(mockPrismaService.profile.upsert).toHaveBeenCalledWith({
        where: { id: 'uuid-123' },
        create: {
          id: 'uuid-123',
          email: 'john.doe@example.com',
          firstName: 'John',
          lastName: 'Doe',
          updatedAt: expect.any(Date),
        },
        update: {},
      });
    });

    it('should create profile with null lastName for single name email', async () => {
      const payload: ProfileWebhookDto = {
        id: 'uuid-123',
        email: 'johndoe@example.com',
      };

      await service.handleCreateProfileWebhook(payload);

      expect(mockPrismaService.profile.upsert).toHaveBeenCalledWith({
        where: { id: 'uuid-123' },
        create: {
          id: 'uuid-123',
          email: 'johndoe@example.com',
          firstName: 'Johndoe',
          lastName: null,
          updatedAt: expect.any(Date),
        },
        update: {},
      });
    });

    it('should throw if prisma throws', async () => {
      mockPrismaService.profile.upsert.mockRejectedValueOnce(
        new Error('DB error'),
      );

      const payload: ProfileWebhookDto = {
        id: 'uuid-123',
        email: 'john.doe@example.com',
      };

      await expect(service.handleCreateProfileWebhook(payload)).rejects.toThrow(
        'DB error',
      );
    });

    it('should be idempotent on duplicate webhook', async () => {
      const payload: ProfileWebhookDto = {
        id: 'uuid-123',
        email: 'john.doe@example.com',
      };

      await service.handleCreateProfileWebhook(payload);
      await service.handleCreateProfileWebhook(payload);

      expect(mockPrismaService.profile.upsert).toHaveBeenCalledTimes(2);
    });
  });

  describe('findOne', () => {
    it('should return ProfileResponseDto when profile exists', async () => {
      mockPrismaService.profile.findUnique.mockResolvedValueOnce(mockProfile);

      const result = await service.findOne(mockProfile.id);

      expect(mockPrismaService.profile.findUnique).toHaveBeenCalledWith({
        where: { id: mockProfile.id },
      });
      expect(result).toBeInstanceOf(ProfileResponseDto);
      expect(result).toEqual(ProfileResponseDto.fromOrm(mockProfile));
    });

    it('should return null when profile does not exist', async () => {
      mockPrismaService.profile.findUnique.mockResolvedValueOnce(null);

      const result = await service.findOne(mockProfile.id);

      expect(result).toBeNull();
    });

    it('should throw if prisma throws', async () => {
      mockPrismaService.profile.findUnique.mockRejectedValueOnce(
        new Error('DB error'),
      );

      await expect(service.findOne(mockProfile.id)).rejects.toThrow('DB error');
    });
  });

  describe('extractNameFromEmail', () => {
    const cases = [
      {
        email: 'john.doe@example.com',
        expected: { firstName: 'John', lastName: 'Doe' },
      },
      {
        email: 'johndoe@example.com',
        expected: { firstName: 'Johndoe', lastName: null },
      },
      {
        email: 'john.doe.smith@example.com',
        expected: { firstName: 'John', lastName: 'Doe Smith' },
      },
      {
        email: "o'brien@example.com",
        expected: { firstName: "O'brien", lastName: null },
      },
      {
        email: 'john-doe@example.com',
        expected: { firstName: 'John-doe', lastName: null },
      },
      {
        email: '123invalid@example.com',
        expected: { firstName: '123invalid', lastName: null },
      },
    ];

    it.each(cases)(
      '$email → firstName: $expected.firstName, lastName: $expected.lastName',
      async ({ email, expected }) => {
        const payload: ProfileWebhookDto = { id: 'uuid-123', email };

        await service.handleCreateProfileWebhook(payload);

        expect(mockPrismaService.profile.upsert).toHaveBeenCalledWith({
          where: { id: 'uuid-123' },
          create: expect.objectContaining({
            firstName: expected.firstName,
            lastName: expected.lastName,
          }),
          update: {},
        });
      },
    );
  });
});

import { Test, TestingModule } from '@nestjs/testing';
import { ProfileService } from './profile.service';
import { PrismaService } from '../prisma/prisma.service';
import { ProfileWebhookDto } from './dto/profile.dto';

const mockPrismaService = {
  profile: {
    upsert: jest.fn().mockResolvedValue(undefined),
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

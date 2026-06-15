import { Test, TestingModule } from '@nestjs/testing';
import { SignupService } from './signup.service';
import { PrismaService } from '../prisma/prisma.service';
import { SignupWebhookDto } from './dto/signup.dto';

const mockPrismaService = {
  profile: {
    create: jest.fn().mockResolvedValue(undefined),
  },
};

describe('SignupService', () => {
  let service: SignupService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        SignupService,
        { provide: PrismaService, useValue: mockPrismaService },
      ],
    }).compile();

    service = module.get<SignupService>(SignupService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('handleWebhook', () => {
    it('should create profile with extracted name from email', async () => {
      const payload: SignupWebhookDto = {
        id: 'uuid-123',
        email: 'john.doe@example.com',
      };

      await service.handleWebhook(payload);

      expect(mockPrismaService.profile.create).toHaveBeenCalledWith({
        data: {
          id: 'uuid-123',
          email: 'john.doe@example.com',
          firstName: 'John',
          lastName: 'Doe',
          updatedAt: expect.any(Date),
        },
      });
    });

    it('should create profile with null lastName for single name email', async () => {
      const payload: SignupWebhookDto = {
        id: 'uuid-123',
        email: 'johndoe@example.com',
      };

      await service.handleWebhook(payload);

      expect(mockPrismaService.profile.create).toHaveBeenCalledWith({
        data: {
          id: 'uuid-123',
          email: 'johndoe@example.com',
          firstName: 'Johndoe',
          lastName: null,
          updatedAt: expect.any(Date),
        },
      });
    });

    it('should throw if prisma throws', async () => {
      mockPrismaService.profile.create.mockRejectedValueOnce(
        new Error('DB error'),
      );

      const payload: SignupWebhookDto = {
        id: 'uuid-123',
        email: 'john.doe@example.com',
      };

      await expect(service.handleWebhook(payload)).rejects.toThrow('DB error');
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
        const payload: SignupWebhookDto = { id: 'uuid-123', email };

        await service.handleWebhook(payload);

        expect(mockPrismaService.profile.create).toHaveBeenCalledWith({
          data: expect.objectContaining({
            firstName: expected.firstName,
            lastName: expected.lastName,
          }),
        });
      },
    );
  });
});

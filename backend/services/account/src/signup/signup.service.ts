import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { SignupWebhookDto } from './dto/signup.dto';

@Injectable()
export class SignupService {
  constructor(private readonly prisma: PrismaService) {}

  async handleWebhook(payload: SignupWebhookDto): Promise<void> {
    const id = payload.id;
    const email = payload.email;
    const { firstName, lastName } = this.extractNameFromEmail(payload.email);

    await this.createProfile(id, email, firstName, lastName);
  }

  private async createProfile(
    id: string,
    email: string,
    firstName: string,
    lastName: string | null,
  ): Promise<void> {
    await this.prisma.profile.upsert({
      where: { id },
      create: { id, email, firstName, lastName, updatedAt: new Date() },
      update: {},
    });
  }

  private extractNameFromEmail(email: string): {
    firstName: string;
    lastName: string | null;
  } {
    const namePart = email.split('@')[0];
    const parts = namePart.split('.');

    const nameRegex = /^[a-zA-Z\s\-']+$/;

    const firstName = parts[0];

    if (!nameRegex.test(firstName)) {
      return { firstName: this.sanitizeName(namePart), lastName: null };
    }

    if (parts.length === 1) {
      return { firstName: this.sanitizeName(firstName), lastName: null };
    }

    const lastName = parts.slice(1).join(' ');

    if (!nameRegex.test(lastName)) {
      return { firstName: this.sanitizeName(firstName), lastName: null };
    }

    return {
      firstName: this.sanitizeName(firstName),
      lastName: this.sanitizeName(lastName),
    };
  }

  private sanitizeName(name: string): string {
    return name
      .split(' ')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(' ');
  }
}

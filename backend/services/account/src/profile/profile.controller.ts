import { Body, Controller, HttpCode, Post, UseGuards } from '@nestjs/common';
import { ProfileService } from './profile.service';
import { ProfileWebhookDto } from './dto/profile.dto';
import { WebhookSecretGuard } from '../guards/webhook-secret/webhook-secret.guard';

@Controller('profiles')
export class ProfileController {
  constructor(private readonly ProfileService: ProfileService) {}

  @Post('webhook')
  @HttpCode(204)
  @UseGuards(WebhookSecretGuard)
  async createProfileWebhook(
    @Body() payload: ProfileWebhookDto,
  ): Promise<void> {
    await this.ProfileService.handleCreateProfileWebhook(payload);
  }
}

import { Body, Controller, HttpCode, Post, UseGuards } from '@nestjs/common';
import { SignupService } from './signup.service';
import { SignupWebhookDto } from './dto/signup.dto';
import { WebhookSecretGuard } from '../guards/webhook-secret/webhook-secret.guard';

@Controller('signup')
export class SignupController {
  constructor(private readonly signupService: SignupService) {}

  @Post('webhook')
  @HttpCode(204)
  @UseGuards(WebhookSecretGuard)
  async handleWebhook(@Body() payload: SignupWebhookDto): Promise<void> {
    await this.signupService.handleWebhook(payload);
  }
}

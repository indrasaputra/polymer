import { Controller } from '@nestjs/common';

@Controller('signup')
export class SignupController {
  constructor(private readonly signupService: SignupService) {}

  @Post('webhook')
  @HttpCode(204)
  async handleWebhook(
    @Headers('x-webhook-secret') secret: string,
    @Body() payload: SupabaseWebhookDto,
  ) {
    if (secret !== process.env.WEBHOOK_SECRET) {
      throw new UnauthorizedException('Secret is invalid');
    }

    await this.signupService.handleWebhook(payload);
  }
}

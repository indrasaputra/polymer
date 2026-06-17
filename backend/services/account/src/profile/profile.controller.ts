import {
  Body,
  Controller,
  HttpCode,
  Post,
  Get,
  UseGuards,
  NotFoundException,
} from '@nestjs/common';
import { ProfileService } from './profile.service';
import { ProfileResponseDto, ProfileWebhookDto } from './dto/profile.dto';
import { WebhookSecretGuard } from '../common/guards/webhook-secret.guard';
import { Public } from '../common/decorators/auth-type.decorator';
import { CurrentUserDecorator } from '../common/decorators/current-user.decorator';
import { CurrentUser } from '../common/dto/current-user.dto';

@Controller('profiles')
export class ProfileController {
  constructor(private readonly ProfileService: ProfileService) {}

  @Post('webhook')
  @HttpCode(204)
  @Public()
  @UseGuards(WebhookSecretGuard)
  async createProfileWebhook(
    @Body() payload: ProfileWebhookDto,
  ): Promise<void> {
    await this.ProfileService.handleCreateProfileWebhook(payload);
  }

  @Get('me')
  @HttpCode(200)
  async getProfile(
    @CurrentUserDecorator() user: CurrentUser,
  ): Promise<ProfileResponseDto> {
    const profile = await this.ProfileService.findOne(user.id);
    if (!profile) {
      throw new NotFoundException(`User with id ${user.id} not found`);
    }
    return profile;
  }
}

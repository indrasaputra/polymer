import { Injectable } from '@nestjs/common';

@Injectable()
export class SignupService {
  async handleWebhook(payload: SupabaseWebhookDto) {
    const { type, table, record } = payload;

    if (table === 'users' && type === 'INSERT') {
      // create profile — coming next
    }
  }
}

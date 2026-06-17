import { Injectable, OnModuleInit, OnModuleDestroy } from '@nestjs/common';
import { PrismaClient } from '../../generated/prisma/client';
import { Pool, type PoolConfig } from 'pg';
import { PrismaPg } from '@prisma/adapter-pg';
import { Config } from '../config/config';

@Injectable()
export class PrismaService
  extends PrismaClient
  implements OnModuleInit, OnModuleDestroy
{
  private readonly pool: Pool;

  constructor(config: Config) {
    const schema = config.postgre.schema;

    const poolConfig: PoolConfig = {
      connectionString: config.postgre.url,
    };
    const pool = new Pool(poolConfig);
    const adapter = new PrismaPg(pool, { schema });

    super({ adapter });
    this.pool = pool;
  }

  async onModuleInit() {
    await this.$connect();
  }

  async onModuleDestroy() {
    await this.$disconnect();
    await this.pool.end();
  }
}

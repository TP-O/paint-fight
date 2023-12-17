import { Type } from 'class-transformer';
import { ValidateNested } from 'class-validator';
import { AppConfig } from './app';
import { RedisConfig } from './redis';
import { SupabaseConfig } from './supabase';

export class RootConfig {
  @Type(() => AppConfig)
  @ValidateNested()
  public readonly app!: AppConfig;

  @Type(() => RedisConfig)
  @ValidateNested()
  public readonly redis!: RedisConfig;

  @Type(() => SupabaseConfig)
  @ValidateNested()
  public readonly supabase!: SupabaseConfig;
}

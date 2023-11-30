import { IsNotEmpty, IsString } from 'class-validator';

export class SupabaseConfig {
  @IsString()
  @IsNotEmpty()
  public readonly url!: string;

  @IsString()
  @IsNotEmpty()
  public readonly serviceRoleKey!: string;
}

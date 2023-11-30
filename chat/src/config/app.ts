import {
  IsBoolean,
  IsIn,
  IsNumber,
  IsString,
  Max,
  Min,
  MinLength,
} from 'class-validator';
import { AppEnv } from 'src/enum/app';

export class AppConfig {
  @IsIn(Object.values(AppEnv))
  public readonly env!: AppEnv;

  @IsBoolean()
  readonly debug!: boolean;

  @IsNumber()
  @Min(10)
  @Max(65535)
  public readonly port!: number;

  @IsString()
  @MinLength(20)
  public readonly secret!: string;
}

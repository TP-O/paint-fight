import { Module } from '@nestjs/common';
import { ChatModule } from './module/chat/chat.module';
import { TypedConfigModule } from 'nest-typed-config';
import { RootConfig } from './config/root';
import { loadConfig } from './utils/load-config';
import { RedisService } from './external/redis';
import { LoggerService } from './service/logger';

@Module({
  imports: [
    ChatModule,
    TypedConfigModule.forRoot({
      schema: RootConfig,
      load: loadConfig,
    }),
  ],
  providers: [RedisService, LoggerService],
})
export class AppModule {}

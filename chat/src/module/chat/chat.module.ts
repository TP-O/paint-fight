import { Module } from '@nestjs/common';
import { ChatService } from './chat.service';
import { ChatGateway } from './chat.gateway';
import { AuthService } from '@service/auth';
import { LoggerService } from '@service/logger';
import { SupabaseService } from '@external/supabase';
import { RedisService } from '@external/redis';
import { CacheModule } from '@nestjs/cache-manager';

@Module({
  imports: [CacheModule.register()],
  providers: [ChatGateway, ChatService, ChatService, AuthService, LoggerService, SupabaseService, RedisService],
  exports: [],
})
export class ChatModule {}

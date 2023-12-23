import { Module } from '@nestjs/common';
import { ChatService } from './chat.service';
import { ChatGateway } from './chat.gateway';
import { AuthService } from 'src/service/auth';
import { LoggerService } from 'src/service/logger';
import { SupabaseService } from 'src/external/supabase';
import { RedisService } from 'src/external/redis';
import { CacheModule } from '@nestjs/cache-manager';

@Module({
  imports: [CacheModule.register()],
  providers: [ChatGateway, ChatService, ChatService, AuthService, LoggerService, SupabaseService, RedisService],
  exports: [],
})
export class ChatModule {}

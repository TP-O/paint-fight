import { Module } from '@nestjs/common';
import { ChatService } from './chat.service';
import { ChatGateway } from './chat.gateway';
import { AuthService } from 'src/service/auth';
import { LoggerService } from 'src/service/logger';
import { SupabaseService } from 'src/external/supabase';

@Module({
  providers: [ChatGateway, ChatService, ChatService, AuthService, LoggerService, SupabaseService],
  exports: [ChatGateway],
})
export class ChatModule {}

import { Module } from '@nestjs/common';
import { RoomController } from './room.controller';
import { LoggerService } from '@service/logger';
import { ChatModule } from '@module/chat/chat.module';
import { CacheModule } from '@nestjs/cache-manager';
import { CacheConfig } from '@config/cache';
import { RedisService } from '@external/redis';

@Module({
  imports: [ChatModule, CacheModule.register(CacheConfig)],
  controllers: [RoomController],
  providers: [LoggerService, RedisService],
})
export class RoomModule {}

import { Module } from '@nestjs/common';
import { RoomController } from './room.controller';
import { ChatModule } from '@module/chat/chat.module';
import { CacheModule } from '@nestjs/cache-manager';
import { CacheConfig } from '@config/cache';
import { RoomService } from './room.service';
import { LoggerService } from '@service/logger';
import { RedisService } from '@external/redis.service';

@Module({
  imports: [ChatModule, CacheModule.register(CacheConfig)],
  controllers: [RoomController],
  providers: [RoomService, LoggerService, RedisService],
})
export class RoomModule {}

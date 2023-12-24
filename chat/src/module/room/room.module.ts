import { Module } from '@nestjs/common';
import { RoomController } from './room.controller';
import { LoggerService } from 'src/service/logger';

@Module({
  controllers: [RoomController],
  providers: [LoggerService],
})
export class RoomModule {}

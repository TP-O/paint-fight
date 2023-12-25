import { Controller, UseFilters, UsePipes, ValidationPipe } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { AddPlayerToRoomRequest, AddPlayerToRoomResponse } from './dto/add-player-to-room';
import { RemovePlayerFromRoomRequest, RemovePlayerFromRoomResponse } from './dto/remove-player-from-room';
import { AllExceptionFilter } from '@exception/all-exception.filter';
import { GrpcExceptionFilter } from '@exception/grpc-exception.filter';
import { Code } from '@enum/code';
import { RoomService } from './room.service';

@Controller()
@UseFilters(AllExceptionFilter, GrpcExceptionFilter)
@UsePipes(
  new ValidationPipe({
    whitelist: true,
  }),
)
export class RoomController {
  constructor(private roomService: RoomService) {}

  @GrpcMethod('RoomService', 'AddPlayerToRoom')
  async addPlayerToRoom(payload: AddPlayerToRoomRequest): Promise<AddPlayerToRoomResponse> {
    await this.roomService.addPlayer(payload);
    return {
      ok: true,
      code: Code.Ok,
      data: undefined,
    };
  }

  @GrpcMethod('RoomService', 'RemovePlayerFromRoom')
  async removePlayerFromRoom(payload: RemovePlayerFromRoomRequest): Promise<RemovePlayerFromRoomResponse> {
    await this.roomService.removePlayer(payload);
    return {
      ok: true,
      code: Code.Ok,
      data: undefined,
    };
  }
}

import { Controller, UseFilters, UsePipes, ValidationPipe } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { AddPlayerToRoomRequest, AddPlayerToRoomResponse } from './dto/add-player-to-room';
import { RemovePlayerFromRoomRequest, RemovePlayerFromRoomResponse } from './dto/remove-player-from-room';
import { AllExceptionFilter } from '@filter/all-exception.filter';
import { GrpcExceptionFilter } from '@filter/grpc-exception.filter';

@Controller()
@UseFilters(AllExceptionFilter, GrpcExceptionFilter)
@UsePipes(
  new ValidationPipe({
    whitelist: true,
  }),
)
export class RoomController {
  @GrpcMethod('RoomService', 'AddPlayerToRoom')
  addPlayerToRoom(data: AddPlayerToRoomRequest): AddPlayerToRoomResponse {
    console.log(data);

    return {
      ok: true,
      code: '',
      data: undefined,
    };
  }

  @GrpcMethod('RoomService', 'RemovePlayerFromRoom')
  removePlayerFromRoom(data: RemovePlayerFromRoomRequest): RemovePlayerFromRoomResponse {
    console.log(data);

    return {
      ok: true,
      code: '',
      data: undefined,
    };
  }
}

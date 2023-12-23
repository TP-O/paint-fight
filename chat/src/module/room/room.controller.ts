import { Controller } from '@nestjs/common';
import { GrpcMethod, RpcException } from '@nestjs/microservices';
import { AddPlayerToRoomRequest, AddPlayerToRoomResponse } from './dto/add-player-to-room';
import { RemovePlayerFromRoomRequest, RemovePlayerFromRoomResponse } from './dto/remove-player-from-room';
import { status } from '@grpc/grpc-js';

@Controller()
export class RoomController {
  @GrpcMethod('RoomService', 'AddPlayerToRoom')
  addPlayerToRoom(data: AddPlayerToRoomRequest): AddPlayerToRoomResponse {
    console.log(data);

    // throw new RpcException({
    //   code: 'YOUR_CUSTOM_ERROR_CODE',
    //   message: 'Your custom error message',
    //   // You can also include additional details in the argument object
    //   details: {
    //     // additional properties
    //     a: 1,
    //     b: 2,
    //   },
    // });

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

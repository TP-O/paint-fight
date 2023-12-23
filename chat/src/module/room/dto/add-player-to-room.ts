import { IsNotEmpty, IsString } from 'class-validator';
import { OkResponse } from 'src/type';

export class AddPlayerToRoomRequest {
  @IsString()
  @IsNotEmpty()
  roomId!: string;

  @IsString()
  @IsNotEmpty()
  playerId!: string;
}

export type AddPlayerToRoomResponse = OkResponse;

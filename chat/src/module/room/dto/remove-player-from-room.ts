import { IsNotEmpty, IsString } from 'class-validator';
import { OkResponse } from '@types';

export class RemovePlayerFromRoomRequest {
  @IsString()
  @IsNotEmpty()
  roomId!: string;

  @IsString()
  @IsNotEmpty()
  playerId!: string;
}

export type RemovePlayerFromRoomResponse = OkResponse;

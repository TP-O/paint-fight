import { ErrResponse } from '@types';
import { EmitEvent, ListenEvent } from './chat.enum';
import { Server, Socket } from 'socket.io';

type SuccessResponse = {
  message: string;
};

type PrivateMessageData = {
  senderId: string;
  content: string;
};

type RoomMessageData = PrivateMessageData & {
  roomId: string;
};

export type EmitEventMap = {
  [EmitEvent.Error]: (response: ErrResponse) => void;
  [EmitEvent.Success]: (response: SuccessResponse) => void;
  [EmitEvent.PrivateMessage]: (data: PrivateMessageData) => void;
  [EmitEvent.RoomMessage]: (data: RoomMessageData) => void;
};

type ChatSocketData = {
  playerId: string;
};

export type ChatSocket = Socket<any, EmitEventMap, any, ChatSocketData> & {
  event: ListenEvent;
};

export type ChatSocketServer = Server<any, EmitEventMap, any, ChatSocketData>;

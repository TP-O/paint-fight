import { Server, Socket } from 'socket.io';
import { ListenEvent } from './event.enum';
import { EmitEventMap } from './event.type';

type ChatSocketData = {
  playerId: string;
};

export type ChatSocket = Socket<any, EmitEventMap, any, ChatSocketData> & {
  event: ListenEvent;
};

export type ChatSocketServer = Server<any, EmitEventMap, any, ChatSocketData>;

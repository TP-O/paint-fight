import { ErrResponse } from '@types';
import { EmitEvent } from './event.enum';

type PrivateMessageData = {
  senderId: string;
  content: string;
};

type RoomMessageData = PrivateMessageData & {
  roomId: string;
};

export type EmitEventMap = {
  [EmitEvent.Error]: (response: ErrResponse) => void;
  [EmitEvent.PrivateMessage]: (data: PrivateMessageData) => void;
  [EmitEvent.RoomMessage]: (data: RoomMessageData) => void;
};

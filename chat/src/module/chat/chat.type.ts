import { EmitEvent } from './chat.enum';

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
  [EmitEvent.Error]: (response: any) => void;
  [EmitEvent.Success]: (response: SuccessResponse) => void;
  [EmitEvent.PrivateMessage]: (data: PrivateMessageData) => void;
  [EmitEvent.RoomMessage]: (data: RoomMessageData) => void;
};

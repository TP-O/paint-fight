import { Injectable, UseFilters, UseInterceptors, UsePipes, ValidationPipe } from '@nestjs/common';
import {
  ConnectedSocket,
  GatewayMetadata,
  MessageBody,
  OnGatewayConnection,
  OnGatewayDisconnect,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';
import { ChatService } from './chat.service';
import { ListenEvent } from './event.enum';
import { AllExceptionFilter } from '@exception/all-exception.filter';
import { EventInterceptor } from './event.interceptor';
import { SendPrivateMessageRequest } from './dto/send-private-message';
import { SendRoomMessageRequest } from './dto/send-room-message';
import { WsExceptionFilter } from '@exception/ws-exception.filter';
import { ChatSocket, ChatSocketServer } from './socketio.type';

@Injectable()
@UseFilters(AllExceptionFilter, WsExceptionFilter)
@UsePipes(
  new ValidationPipe({
    whitelist: true,
  }),
)
@WebSocketGateway<GatewayMetadata>({
  namespace: '/',
})
export class ChatGateway implements OnGatewayConnection, OnGatewayDisconnect {
  @WebSocketServer()
  readonly server!: ChatSocketServer;

  constructor(private chatService: ChatService) {}

  /**
   * Store player state before connection.
   */
  async handleConnection(client: ChatSocket): Promise<void> {
    await this.chatService.connect(client);
  }

  /**
   * Clear player state after disconnection.
   */
  async handleDisconnect(client: ChatSocket): Promise<void> {
    await this.chatService.disconnect(client);
  }

  /**
   * Send private message.
   */
  @UseInterceptors(new EventInterceptor(ListenEvent.PrivateMessage))
  @SubscribeMessage(ListenEvent.PrivateMessage)
  async sendPrivateMessage(
    @ConnectedSocket() client: ChatSocket,
    @MessageBody() payload: SendPrivateMessageRequest,
  ): Promise<void> {
    await this.chatService.sendPrivateMessage(client, payload);
  }

  /**
   * Send room message.
   */
  @UseInterceptors(new EventInterceptor(ListenEvent.RoomMessage))
  @SubscribeMessage(ListenEvent.RoomMessage)
  async sendRoomMesage(
    @ConnectedSocket() client: ChatSocket,
    @MessageBody() payload: SendRoomMessageRequest,
  ): Promise<void> {
    await this.chatService.sendRoomMessage(client, payload);
  }
}

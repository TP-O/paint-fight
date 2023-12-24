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
import { ListenEvent } from './chat.enum';
import { ChatSocket, ChatSocketServer } from './chat.type';
import { AllExceptionFilter } from '@filter/all-exception.filter';
import { EventBindingInterceptor } from './interceptor/event-binding';
import { SendPrivateMessageDto } from './dto/send-private-message';
import { SendRoomMessageDto } from './dto/send-room-message';
import { WsExceptionFilter } from '@filter/ws-exception.filter';

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
   *
   * @param client
   */
  async handleConnection(client: ChatSocket): Promise<void> {
    await this.chatService.connect(client);
  }

  /**
   * Clear player state after disconnection.
   *
   * @param client
   */
  async handleDisconnect(client: ChatSocket): Promise<void> {
    await this.chatService.disconnect(client);
  }

  /**
   * Send private message.
   *
   * @param client
   * @param payload
   */
  @UseInterceptors(new EventBindingInterceptor(ListenEvent.PrivateMessage))
  @SubscribeMessage(ListenEvent.PrivateMessage)
  async sendPrivateMessage(
    @ConnectedSocket() client: ChatSocket,
    @MessageBody() payload: SendPrivateMessageDto,
  ): Promise<void> {
    await this.chatService.sendPrivateMessage(client, payload);
  }

  /**
   * Send room message.
   *
   * @param client
   * @param payload
   */
  @UseInterceptors(new EventBindingInterceptor(ListenEvent.RoomMessage))
  @SubscribeMessage(ListenEvent.RoomMessage)
  async sendRoomMesage(
    @ConnectedSocket() client: ChatSocket,
    @MessageBody() payload: SendRoomMessageDto,
  ): Promise<void> {
    await this.chatService.sendRoomMessage(client, payload);
  }
}

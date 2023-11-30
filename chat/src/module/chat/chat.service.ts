import { Injectable, UnauthorizedException } from '@nestjs/common';
import { Server, Socket } from 'socket.io';
import { AuthService } from '../../service/auth';
import { EmitEvent, ListenEvent } from './chat.enum';
import { EmitEventMap } from './chat.type';
import { LoggerService } from 'src/service/logger';
import { User } from '@supabase/supabase-js';
import { SendPrivateMessageDto } from './dto/send-private-message';
import { SendRoomMessageDto } from './dto/send-room-message';

@Injectable()
export class ChatService {
  constructor(
    private authService: AuthService,
    private logger: LoggerService,
  ) {
    this.logger.setContext(ChatService.name);
  }

  /**
   * Connect the client.
   *
   * @param server websocket server.
   * @param client socket client.
   */
  async connect(
    server: Server<EmitEventMap>,
    client: Socket<EmitEventMap, EmitEventMap, EmitEventMap, { id: string }>,
  ): Promise<void> {
    try {
      const user = await this._validateConnection(server, client);
      client.data.id = user.id;
    } catch (error: any) {
      client.emit(EmitEvent.Error, {
        event: ListenEvent.Connect,
        message: error.message,
      });
      client.disconnect();
    }
  }

  /**
   * Check if the connection satisfies some sepecific conditions
   * before allowing the connection.
   *
   * @param server websocket server.
   * @param client socket client.
   */
  private async _validateConnection(
    server: Server<EmitEventMap>,
    client: Socket,
  ): Promise<User> {
    const user = await this._validateAuthorization(
      client.handshake.headers.authorization ?? '',
    );

    // TODO: check duplicate login

    // const sid = await this.playerService.getSocketId(player.id);
    // if (sid) {
    //   server.to(sid).emit(EmitEvent.Error, {
    //     event: ListenEvent.Connect,
    //     message: 'This account is being connected by someone else!',
    //   });
    //   server.to(sid).disconnectSockets();
    //   this.disconnect(server, client);
    // }

    return user;
  }

  /**
   * Verify token.
   *
   * @param headerAuthorization
   */
  private async _validateAuthorization(
    headerAuthorization: string,
  ): Promise<User> {
    const token = String(headerAuthorization).replace('Bearer ', '');
    if (!token) {
      throw new UnauthorizedException('Missing access token!');
    }

    return this.authService.getUser(token);
  }

  /**
   * Disconnect the client.
   *
   * @param server websocket server.
   * @param client socket client.
   */
  async disconnect(
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    server: Server<EmitEventMap>,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    client: Socket<EmitEventMap, EmitEventMap, EmitEventMap, { id: string }>,
  ): Promise<void> {
    //
  }

  /**
   * Send a private message to friend.
   *
   * @param server
   * @param client
   * @param payload
   */
  async sendPrivateMessage(
    server: Server<EmitEventMap>,
    client: Socket<EmitEventMap, EmitEventMap, EmitEventMap, { id: string }>,
    payload: SendPrivateMessageDto,
  ): Promise<void> {
    if (!client.data.id) {
      return;
    }

    // const sid = await this.playerService.getSocketId(payload.receiverId);
    // if (!sid) {
    //   throw new BadRequestException('This player is offline!');
    // }
    const sid = '';

    server.to(sid).emit(EmitEvent.PrivateMessage, {
      ...payload,
      senderId: client.data.id,
    });
  }

  /**
   * Send a message to joined room.
   *
   * @param server
   * @param client
   * @param payload
   */
  async sendRoomMessage(
    server: Server<EmitEventMap>,
    client: Socket<EmitEventMap, EmitEventMap, EmitEventMap, { id: string }>,
    payload: SendRoomMessageDto,
  ): Promise<void> {
    if (!client.data.id) {
      return;
    }

    // const room = await this.roomService.get(payload.roomId);
    // if (!room) {
    //   throw new NotFoundException("Room doesn't exist!");
    // }

    // if (!room.memberIds.includes(client.data.id)) {
    //   throw new ForbiddenException('You are not in this room!');
    // }

    server.to(payload.roomId).emit(EmitEvent.RoomMessage, {
      ...payload,
      senderId: client.data.id,
    });
  }
}

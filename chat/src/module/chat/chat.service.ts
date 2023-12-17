import { Injectable, UnauthorizedException } from '@nestjs/common';
import { Server, Socket } from 'socket.io';
import { AuthService } from '../../service/auth';
import { EmitEvent, ListenEvent } from './chat.enum';
import { EmitEventMap } from './chat.type';
import { LoggerService } from 'src/service/logger';
import { User } from '@supabase/supabase-js';
import { SendPrivateMessageDto } from './dto/send-private-message';
import { SendRoomMessageDto } from './dto/send-room-message';
import { RedisService } from 'src/external/redis';
import { REDIS } from 'src/constant/redis';

@Injectable()
export class ChatService {
  constructor(
    private authService: AuthService,
    private logger: LoggerService,
    private redis: RedisService,
  ) {
    this.logger.setContext(ChatService.name);
  }

  /**
   * Connect the client.
   */
  async connect(client: Socket<EmitEventMap, EmitEventMap, EmitEventMap, { userId: string }>): Promise<void> {
    try {
      const user = await this._verifyClient(client);
      client.data.userId = user.id;
    } catch (error: any) {
      client.emit(EmitEvent.Error, {
        event: ListenEvent.Connect,
        message: error.message,
      });
      client.disconnect();
    }
  }

  /**
   * Check some rules to verify the client. Return user data if successful.
   */
  private async _verifyClient(client: Socket): Promise<User> {
    const token = String(client.handshake.headers.authorization).replace('Bearer ', '');
    if (!token) {
      throw new UnauthorizedException('Token is required!');
    }
    const user = await this.authService.getUser(token);

    if (await this.redis.client.getset(`${REDIS.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${user.id}`, user.id)) {
      throw new Error('This account is being connected by someone else!');
    }

    return user;
  }

  /**
   * Disconnect the client.
   */
  async disconnect(
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    client: Socket<EmitEventMap, EmitEventMap, EmitEventMap, { userId: string }>,
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

import { Inject, Injectable } from '@nestjs/common';
import { AuthService } from '../../service/auth';
import { EmitEvent } from './chat.enum';
import { ChatSocket } from './chat.type';
import { LoggerService } from '@service/logger';
import { User } from '@supabase/supabase-js';
import { SendPrivateMessageDto } from './dto/send-private-message';
import { SendRoomMessageDto } from './dto/send-room-message';
import { RedisService } from '@external/redis';
import { REDIS } from '@constant/redis';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import { Cache } from 'cache-manager';
import { CACHE } from '@constant/cache';
import { WsException } from '@nestjs/websockets';
import { PublicError } from '@filter/public-error.error';
import { Code } from '@enum/code';
import { MillisecondsTime } from '@enum/time';

@Injectable()
export class ChatService {
  constructor(
    private authService: AuthService,
    private logger: LoggerService,
    private redis: RedisService,
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
  ) {
    this.logger.setContext(ChatService.name);
  }

  /**
   * Connect the client.
   */
  async connect(client: ChatSocket): Promise<void> {
    const player = await this._verifyClient(client);
    client.data.playerId = player.id;
  }

  /**
   * Check some rules to verify the client. Return player data if successful.
   */
  private async _verifyClient(client: ChatSocket): Promise<User> {
    const token = String(client.handshake.headers.authorization).replace('Bearer ', '');
    if (!token) {
      throw new WsException('Token is required!');
    }
    const player = await this.authService.getUser(token);

    if (await this.redis.client.getset(`${REDIS.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${player.id}`, client.id)) {
      throw new WsException('This account is being connected by someone else!');
    }

    await this.cacheManager.set(
      `${CACHE.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${player.id}`,
      client.id,
      5 * MillisecondsTime.Miniute,
    );
    return player;
  }

  /**
   * Disconnect the client.
   */
  async disconnect(
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    client: ChatSocket,
  ): Promise<void> {
    //
  }

  /**
   * Send a private message to friend.
   */
  async sendPrivateMessage(client: ChatSocket, payload: SendPrivateMessageDto): Promise<void> {
    if (!client.data.playerId) {
      client.disconnect();
      return;
    }

    let sid = await this.cacheManager.get<string | null>(
      `${CACHE.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${payload.receiverId}`,
    );
    if (!sid) {
      sid = await this.redis.client.get(`${REDIS.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${payload.receiverId}`);
      if (!sid) {
        return;
      } else {
        await this.cacheManager.set(
          `${CACHE.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${payload.receiverId}`,
          client.id,
          5 * MillisecondsTime.Miniute,
        );
      }
    }

    client.to(sid).emit(EmitEvent.PrivateMessage, {
      ...payload,
      senderId: client.data.playerId,
    });
  }

  /**
   * Send a message to joined room.
   */
  async sendRoomMessage(client: ChatSocket, payload: SendRoomMessageDto): Promise<void> {
    if (!client.data.playerId) {
      client.disconnect();
      return;
    }

    if (client.rooms.has(payload.roomId)) {
      throw new PublicError(Code.RoomDoesNotExist);
    }

    client.to(payload.roomId).emit(EmitEvent.RoomMessage, {
      ...payload,
      senderId: client.data.playerId,
    });
  }
}

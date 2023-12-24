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

    if (await this.redis.client.get(`${REDIS.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${player.id}`)) {
      throw new WsException('This account is being connected by someone else!');
    }

    this.redis.client
      .set(`${REDIS.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${player.id}`, client.id, 'EX', MillisecondsTime.Forever)
      .then(() =>
        this.cacheManager.set(
          `${CACHE.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${player.id}`,
          client.id,
          MillisecondsTime.Forever,
        ),
      );

    return player;
  }

  /**
   * Disconnect the client.
   */
  async disconnect(client: ChatSocket): Promise<void> {
    await this.redis.client.del(`${REDIS.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${client.data.playerId}`);
    await this.cacheManager.del(`${CACHE.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${client.data.playerId}`);
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
      }
    }

    this.cacheManager.set(
      `${CACHE.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${payload.receiverId}`,
      client.id,
      MillisecondsTime.Forever,
    );

    client.to(sid).emit(EmitEvent.PrivateMessage, {
      senderId: client.data.playerId,
      content: payload.content,
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
      senderId: client.data.playerId,
      roomId: payload.roomId,
      content: payload.content,
    });
  }
}

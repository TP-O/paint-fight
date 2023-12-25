import { Inject, Injectable } from '@nestjs/common';
import { AuthService } from '../../service/auth';
import { EmitEvent } from './event.enum';
import { LoggerService } from '@service/logger';
import { User } from '@supabase/supabase-js';
import { RedisService } from '@external/redis.service';
import { REDIS } from '@constant/redis';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import { Cache } from 'cache-manager';
import { CACHE } from '@constant/cache';
import { WsException } from '@nestjs/websockets';
import { PublicError } from '@exception/public-error.error';
import { Code } from '@enum/code';
import { MillisecondsTime } from '@enum/time';
import { ChatSocket } from './socketio.type';
import { SendPrivateMessageRequest } from './dto/send-private-message';
import { SendRoomMessageRequest } from './dto/send-room-message';

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

    this.logger.debug(`Player [${player.id}] connected with socket id [${client.id}]`);
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

    // Check in current node and other nodes
    if (
      (await this.cacheManager.get<string | null>(`${CACHE.KEY.PLAYER_ID_TO_SOCKET_ID}${player.id}`)) ||
      (await this.redis.client.get(`${REDIS.NAMESAPCE.PLAYER_ID_TO_SOCKET_ID}${player.id}`))
    ) {
      throw new WsException('This account is being connected by someone else!');
    }

    this.redis.client
      .set(`${REDIS.NAMESAPCE.PLAYER_ID_TO_SOCKET_ID}${player.id}`, client.id, 'EX', MillisecondsTime.Forever)
      .then(() =>
        this.cacheManager.set(`${CACHE.KEY.PLAYER_ID_TO_SOCKET_ID}${player.id}`, client.id, MillisecondsTime.Forever),
      );

    return player;
  }

  /**
   * Disconnect the client.
   */
  async disconnect(client: ChatSocket): Promise<void> {
    /**
     * TODO: check if other nodes have client data or not
     * This assumes that only node that established the connection to the disconnected client has client data.
     */
    if (client.data.playerId) {
      await this.redis.client.del(`${REDIS.NAMESAPCE.PLAYER_ID_TO_SOCKET_ID}${client.data.playerId}`);
    }

    await this.cacheManager.del(`${CACHE.KEY.PLAYER_ID_TO_SOCKET_ID}${client.data.playerId}`);

    this.logger.debug(`Player [${client.data.playerId}] with socket id [${client.id}] disconnected`);
  }

  /**
   * Send a private message to friend.
   */
  async sendPrivateMessage(client: ChatSocket, payload: SendPrivateMessageRequest): Promise<void> {
    if (!client.data.playerId) {
      client.disconnect();
      return;
    }

    let sid = await this.cacheManager.get<string | null>(`${CACHE.KEY.PLAYER_ID_TO_SOCKET_ID}${payload.receiverId}`);
    if (!sid) {
      sid = await this.redis.client.get(`${REDIS.NAMESAPCE.PLAYER_ID_TO_SOCKET_ID}${payload.receiverId}`);
      if (!sid) {
        return;
      } else {
        this.cacheManager.set(
          `${CACHE.KEY.PLAYER_ID_TO_SOCKET_ID}${payload.receiverId}`,
          client.id,
          MillisecondsTime.Forever,
        );
      }
    }

    client.to(sid).emit(EmitEvent.PrivateMessage, {
      senderId: client.data.playerId,
      content: payload.content,
    });
  }

  /**
   * Send a message to joined room.
   */
  async sendRoomMessage(client: ChatSocket, payload: SendRoomMessageRequest): Promise<void> {
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

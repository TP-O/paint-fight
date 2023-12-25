import { RedisService } from '@external/redis.service';
import { ChatGateway } from '@module/chat/chat.gateway';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import { Inject, Injectable } from '@nestjs/common';
import { Cache } from 'cache-manager';
import { AddPlayerToRoomRequest } from './dto/add-player-to-room';
import { CACHE } from '@constant/cache';
import { REDIS } from '@constant/redis';
import { MillisecondsTime } from '@enum/time';
import { RemovePlayerFromRoomRequest } from './dto/remove-player-from-room';
import { PublicError } from '@exception/public-error.error';
import { Code } from '@enum/code';
import { LoggerService } from '@service/logger';

@Injectable()
export class RoomService {
  constructor(
    private chatGateway: ChatGateway,
    private redis: RedisService,
    private logger: LoggerService,
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
  ) {
    this.logger.setContext(RoomService.name);
  }

  async addPlayer(payload: AddPlayerToRoomRequest): Promise<void> {
    const sid = await this.getPlayerSocketId(payload.playerId);
    if (!sid) {
      throw new PublicError(Code.PlayerDoesNotExist);
    }

    this.chatGateway.server.to(sid).socketsJoin(payload.roomId);
  }

  async removePlayer(payload: RemovePlayerFromRoomRequest): Promise<void> {
    const sid = await this.getPlayerSocketId(payload.playerId);
    if (!sid) {
      return;
    }

    this.chatGateway.server.to(sid).socketsLeave(payload.roomId);
  }

  async getPlayerSocketId(playerId: string): Promise<string | null> {
    let sid = await this.cacheManager.get<string | null>(`${CACHE.KEY.PLAYER_ID_TO_SOCKET_ID}${playerId}`);
    if (!sid) {
      this.logger.debug(`Player [${playerId}] does not exist in cache`);

      sid = await this.redis.client.get(`${REDIS.NAMESAPCE.PLAYER_ID_TO_SOCKET_ID}${playerId}`);
      if (!sid) {
        return null;
      } else {
        this.logger.debug(`Player [${playerId}] is added to cache`);

        this.cacheManager.set(`${CACHE.KEY.PLAYER_ID_TO_SOCKET_ID}${playerId}`, sid, MillisecondsTime.Forever);
      }
    }

    return sid;
  }
}

import { BadRequestException, Controller, Inject, UseFilters, UsePipes, ValidationPipe } from '@nestjs/common';
import { GrpcMethod } from '@nestjs/microservices';
import { AddPlayerToRoomRequest, AddPlayerToRoomResponse } from './dto/add-player-to-room';
import { RemovePlayerFromRoomRequest, RemovePlayerFromRoomResponse } from './dto/remove-player-from-room';
import { AllExceptionFilter } from '@filter/all-exception.filter';
import { GrpcExceptionFilter } from '@filter/grpc-exception.filter';
import { ChatGateway } from '@module/chat/chat.gateway';
import { RedisService } from '@external/redis';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import { Cache } from 'cache-manager';
import { CACHE } from '@constant/cache';
import { REDIS } from '@constant/redis';
import { MillisecondsTime } from '@enum/time';
import { Code } from '@enum/code';

@Controller()
@UseFilters(AllExceptionFilter, GrpcExceptionFilter)
@UsePipes(
  new ValidationPipe({
    whitelist: true,
  }),
)
export class RoomController {
  constructor(
    private chatGateway: ChatGateway,
    private redis: RedisService,
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
  ) {}

  @GrpcMethod('RoomService', 'AddPlayerToRoom')
  async addPlayerToRoom(payload: AddPlayerToRoomRequest): Promise<AddPlayerToRoomResponse> {
    let sid = await this.cacheManager.get<string | null>(
      `${CACHE.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${payload.playerId}`,
    );
    if (!sid) {
      sid = await this.redis.client.get(`${REDIS.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${payload.playerId}`);
      if (!sid) {
        throw new BadRequestException('Player is not online');
      }
    }

    this.cacheManager.set(
      `${CACHE.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${payload.playerId}`,
      sid,
      MillisecondsTime.Forever,
    );
    this.chatGateway.server.to(sid).socketsJoin(payload.roomId);

    return {
      ok: true,
      code: Code.Ok,
      data: undefined,
    };
  }

  @GrpcMethod('RoomService', 'RemovePlayerFromRoom')
  async removePlayerFromRoom(payload: RemovePlayerFromRoomRequest): Promise<RemovePlayerFromRoomResponse> {
    let sid = await this.cacheManager.get<string | null>(
      `${CACHE.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${payload.playerId}`,
    );
    if (!sid) {
      sid = await this.redis.client.get(`${REDIS.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${payload.playerId}`);
      if (!sid) {
        throw new BadRequestException('Player is not online');
      }
    }

    this.cacheManager.set(
      `${CACHE.PLAYER_ID_TO_SOCKET_ID_NAMESPACE}${payload.playerId}`,
      sid,
      MillisecondsTime.Forever,
    );
    this.chatGateway.server.to(sid).socketsLeave(payload.roomId);

    return {
      ok: true,
      code: Code.Ok,
      data: undefined,
    };
  }
}

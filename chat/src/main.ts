import { NestFactory } from '@nestjs/core';
import { FastifyAdapter, NestFastifyApplication } from '@nestjs/platform-fastify';
import { AppModule } from './app.module';
import { AppConfig } from './config/app';
import { SocketIoAdapter } from './module/chat/socketio.adapter';
import { RedisService } from './external/redis.service';
import { MicroserviceOptions } from '@nestjs/microservices';
import { addReflectionToGrpcConfig } from 'nestjs-grpc-reflection';
import { RoomServiceConfig } from '@config/microservices';

async function bootstrap() {
  // TODO: remove fastify (dont use http)
  const app = await NestFactory.create<NestFastifyApplication>(AppModule, new FastifyAdapter());
  // TODO: secure grpc connection
  app.connectMicroservice<MicroserviceOptions>(addReflectionToGrpcConfig(RoomServiceConfig));

  const appConfig = app.get(AppConfig);

  app.enableCors({
    origin: true,
    methods: ['GET', 'POST', 'PUT', 'PATCH'],
  });

  const redisService = app.get(RedisService);
  const chatAdapter = new SocketIoAdapter(redisService.client, app);
  await chatAdapter.connectToRedis();
  app.useWebSocketAdapter(chatAdapter);

  await app.startAllMicroservices();
  await app.listen(appConfig.port, '0.0.0.0');
}

bootstrap();

import { NestFactory } from '@nestjs/core';
import { FastifyAdapter, NestFastifyApplication } from '@nestjs/platform-fastify';
import { AppModule } from './app.module';
import { ValidationPipe } from '@nestjs/common';
import { AppConfig } from './config/app';
import { ChatAdapter } from './module/chat/chat.adapter';
import { AllExceptionFilter } from './filter/all-exception';
import { RpcExceptionFilter } from './filter/rpc-exception';
import { RedisService } from './external/redis';
import { LoggerService } from './service/logger';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';
import { join } from 'path';
import { addReflectionToGrpcConfig } from 'nestjs-grpc-reflection';

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(AppModule, new FastifyAdapter());
  // TODO: secure grpc connection
  app.connectMicroservice<MicroserviceOptions>(
    addReflectionToGrpcConfig({
      transport: Transport.GRPC,
      options: {
        package: ['room'],
        protoPath: [join(__dirname, './module/room/proto/service.proto')],
        gracefulShutdown: true,
      },
    }),
    {
      inheritAppConfig: true,
    },
  );

  const appConfig = app.get(AppConfig);

  app.enableCors({
    origin: true,
    methods: ['GET', 'POST', 'PUT', 'PATCH'],
  });

  const redisService = app.get(RedisService);
  const chatAdapter = new ChatAdapter(redisService.client, app);
  await chatAdapter.connectToRedis();
  app.useWebSocketAdapter(chatAdapter);

  app.useGlobalFilters(new AllExceptionFilter(new LoggerService(appConfig)), new RpcExceptionFilter());
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      stopAtFirstError: false,
    }),
  );

  await app.startAllMicroservices();
  await app.listen(appConfig.port, '0.0.0.0');
}

bootstrap();

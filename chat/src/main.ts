import { NestFactory } from '@nestjs/core';
import { FastifyAdapter, NestFastifyApplication } from '@nestjs/platform-fastify';
import { AppModule } from './app.module';
import { ValidationPipe } from '@nestjs/common';
import { AppConfig } from './config/app';
import { ChatAdapter } from './module/chat/chat.adapter';
import { AllExceptionFilter } from './filter/all-exception';
import { HttpExceptionFilter } from './filter/http-exception';
import { RedisService } from './external/redis';
import { LoggerService } from './service/logger';

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(AppModule, new FastifyAdapter());
  const appConfig = app.get(AppConfig);

  app.enableCors({
    origin: true,
    methods: ['GET', 'POST', 'PUT', 'PATCH'],
  });

  const redisService = app.get(RedisService);
  const chatAdapter = new ChatAdapter(redisService.client, app);
  await chatAdapter.connectToRedis();
  app.useWebSocketAdapter(chatAdapter);

  app.setGlobalPrefix('api/v1');
  app.useGlobalFilters(new AllExceptionFilter(new LoggerService(appConfig)), new HttpExceptionFilter());
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      stopAtFirstError: false,
    }),
  );

  await app.listen(appConfig.port, '0.0.0.0');
}

bootstrap();

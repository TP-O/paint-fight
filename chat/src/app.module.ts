import { Module } from '@nestjs/common';
import { ChatModule } from './module/chat/chat.module';
import { TypedConfigModule } from 'nest-typed-config';
import { RootConfig } from './config/root';
import { loadConfig } from './utils/load-config';
import { RoomModule } from './module/room/room.module';
import { GrpcReflectionModule, addReflectionToGrpcConfig } from 'nestjs-grpc-reflection';
import { Transport } from '@nestjs/microservices';
import { join } from 'path';

@Module({
  imports: [
    ChatModule,
    RoomModule,
    TypedConfigModule.forRoot({
      schema: RootConfig,
      load: loadConfig,
    }),
    GrpcReflectionModule.register(
      addReflectionToGrpcConfig({
        transport: Transport.GRPC,
        options: {
          package: ['room'],
          protoPath: [join(__dirname, './module/room/proto/service.proto')],
          gracefulShutdown: true,
        },
      }),
    ),
  ],
})
export class AppModule {}

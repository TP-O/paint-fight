import { Module } from '@nestjs/common';
import { ChatModule } from './module/chat/chat.module';
import { TypedConfigModule } from 'nest-typed-config';
import { RootConfig } from './config/root';
import { loadConfig } from './utils/load-config';
import { RoomModule } from './module/room/room.module';
import { GrpcReflectionModule, addReflectionToGrpcConfig } from 'nestjs-grpc-reflection';
import { RoomServiceConfig } from '@config/microservices';

@Module({
  imports: [
    ChatModule,
    RoomModule,
    TypedConfigModule.forRoot({
      schema: RootConfig,
      load: loadConfig,
    }),
    // TODO: reflect grpc in dev env
    GrpcReflectionModule.register(addReflectionToGrpcConfig(RoomServiceConfig)),
  ],
})
export class AppModule {}

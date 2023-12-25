import { GrpcOptions, Transport } from '@nestjs/microservices';
import { join } from 'path';

export const RoomServiceConfig = Object.freeze<GrpcOptions>({
  transport: Transport.GRPC,
  options: {
    package: ['room'],
    protoPath: [join(__dirname, '../module/room/proto/service.proto')],
    gracefulShutdown: true,
  },
});

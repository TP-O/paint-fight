import { Module } from '@nestjs/common';
import { ChatModule } from './module/chat/chat.module';
import { TypedConfigModule } from 'nest-typed-config';
import { RootConfig } from './config/root';
import { loadConfig } from './utils/load-config';

@Module({
  imports: [
    ChatModule,
    TypedConfigModule.forRoot({
      schema: RootConfig,
      load: loadConfig,
    }),
  ],
})
export class AppModule {}

import { ListenEvent } from './module/chat/chat.enum';
import { User } from '@supabase/supabase-js';

declare module 'socket.io' {
  class Socket {
    event: ListenEvent;
  }
}

declare module 'fastify' {
  export class FastifyRequest {
    user?: User;
  }
}

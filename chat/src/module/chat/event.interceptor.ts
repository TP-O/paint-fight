import { CallHandler, ExecutionContext, Injectable, NestInterceptor } from '@nestjs/common';
import { Observable } from 'rxjs';
import { ListenEvent } from './event.enum';
import { ChatSocket } from './socketio.type';

@Injectable()
export class EventInterceptor implements NestInterceptor {
  constructor(private event: ListenEvent) {}

  async intercept(context: ExecutionContext, next: CallHandler): Promise<Observable<any>> {
    const client = context.switchToWs().getClient<ChatSocket>();
    client.event = this.event;

    // TODO: map unique code sent from client to response error message
    return next.handle();
  }
}

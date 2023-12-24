import { CallHandler, ExecutionContext, Injectable, NestInterceptor } from '@nestjs/common';
import { Observable } from 'rxjs';
import { ListenEvent } from '../chat.enum';
import { ChatSocket } from '../chat.type';

@Injectable()
export class EventBindingInterceptor implements NestInterceptor {
  constructor(private readonly event: ListenEvent) {}

  async intercept(context: ExecutionContext, next: CallHandler): Promise<Observable<any>> {
    const client = context.switchToWs().getClient<ChatSocket>();
    client.event = this.event;
    return next.handle();
  }
}

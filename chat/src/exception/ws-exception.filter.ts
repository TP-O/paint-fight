import { Catch, ArgumentsHost, BadRequestException } from '@nestjs/common';
import { BaseWsExceptionFilter } from '@nestjs/websockets';
import { EmitEvent } from '../module/chat/event.enum';
import { Code } from '@enum/code';
import { PublicError } from './public-error.error';
import { ChatSocket } from '@module/chat/socketio.type';

/**
 * Handle validation errors for websocket.
 */
@Catch(BadRequestException, PublicError)
export class WsExceptionFilter extends BaseWsExceptionFilter {
  catch(exception: BadRequestException | PublicError, host: ArgumentsHost): void {
    let code: Code;
    let error: string | string[];
    if (exception instanceof BadRequestException) {
      code = Code.InvalidArgument;
      error = (exception.getResponse() as any).message;
    } else {
      code = exception.code;
      error = exception.message;
    }

    const client = host.switchToWs().getClient() as ChatSocket;
    client.emit(EmitEvent.Error, {
      ok: false,
      code: `${client.event}.${code}`,
      error,
    });
  }
}

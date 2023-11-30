import {
  Catch,
  ArgumentsHost,
  BadRequestException,
  HttpException,
} from '@nestjs/common';
import { BaseWsExceptionFilter, WsException } from '@nestjs/websockets';
import { Socket } from 'socket.io';
import { EmitEventMap } from '../module/chat/chat.type';
import { ErrorMessage } from 'src/type';
import { EmitEvent } from '../module/chat/chat.enum';

/**
 * Filter the exceptions from gateway.
 *
 * HttpException is caught here because the gateway can use
 * services throwing http exception.
 */
@Catch(WsException, HttpException)
export class WsExceptionFilter extends BaseWsExceptionFilter {
  catch(exception: Error, host: ArgumentsHost): void {
    const client = host.switchToWs().getClient() as Socket<EmitEventMap>;
    let message: ErrorMessage;

    if (exception instanceof BadRequestException) {
      message = (exception.getResponse() as Error).message;
    } else {
      message = exception.message;
    }

    client.emit(EmitEvent.Error, {
      event: client.event,
      message: message,
    });
  }
}

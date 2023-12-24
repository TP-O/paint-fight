import { ServerUnaryCall } from '@grpc/grpc-js';
import { ExceptionFilter, Catch, ArgumentsHost } from '@nestjs/common';
import { WsArgumentsHost } from '@nestjs/common/interfaces';
import { RpcException } from '@nestjs/microservices';
import { WsException } from '@nestjs/websockets';
import { Observable, throwError } from 'rxjs';
import { Socket } from 'socket.io';
import { EmitEvent } from 'src/module/chat/chat.enum';
import { EmitEventMap } from 'src/module/chat/chat.type';
import { LoggerService } from 'src/service/logger';
import { status as grpcStatus } from '@grpc/grpc-js';
import { Code } from 'src/enum/code';
import { ErrResponse } from 'src/type';

/**
 * Filter all unexpected exceptions. All exceptions handled by this filter
 * will be hidden from clients.
 */
@Catch()
export class AllExceptionFilter implements ExceptionFilter {
  constructor(private readonly logger: LoggerService) {}

  catch(exception: Error, host: ArgumentsHost): void | Observable<never> {
    if (host.getType() === 'ws' && !(exception instanceof WsException)) {
      this.logger.error(exception.message, exception.stack, `${AllExceptionFilter.name} - ws`);
      return this._handleWsException(host.switchToWs());
    } else if (host.getType() === 'rpc' && !(exception instanceof RpcException)) {
      this.logger.error(
        exception.message,
        exception.stack,
        `${AllExceptionFilter.name} - rpc(${(host.getArgByIndex(2) as ServerUnaryCall<any, any>).getPath()})`,
      );
      return this._handleRpcException();
    }
  }

  private _handleWsException(host: WsArgumentsHost): void {
    const client = host.getClient() as Socket<EmitEventMap>;
    client.emit(EmitEvent.Error, {
      ok: false,
      code: `${client.event}.${Code.Unknown}`,
      error: 'Unknown error!',
    });
  }

  private _handleRpcException(): Observable<never> {
    return throwError(() => ({
      code: grpcStatus.UNKNOWN,
      message: 'UNKNOWN',
      details: {
        ok: false,
        code: Code.Unknown,
        error: 'Unknown error!',
      } as ErrResponse,
    }));
  }
}

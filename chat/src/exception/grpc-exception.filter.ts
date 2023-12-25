import { ExceptionFilter, Catch, BadRequestException } from '@nestjs/common';
import { Observable, throwError } from 'rxjs';
import { status } from '@grpc/grpc-js';
import { Code } from '@enum/code';
import { GrpcErrResponse } from '@types';
import { PublicError } from './public-error.error';

/**
 * Handle validation errors for grpc.
 */
@Catch(BadRequestException, PublicError)
export class GrpcExceptionFilter implements ExceptionFilter {
  catch(exception: BadRequestException | PublicError): Observable<never> {
    return throwError((): GrpcErrResponse => {
      if (exception instanceof BadRequestException) {
        return {
          code: status.INVALID_ARGUMENT,
          message: Code.InvalidArgument,
          details: {
            ok: false,
            code: Code.InvalidArgument,
            error: (exception.getResponse() as any).message,
          },
        };
      } else {
        return {
          code: status.FAILED_PRECONDITION,
          message: exception.code,
          details: {
            ok: false,
            code: exception.code,
            error: exception.message,
          },
        };
      }
    });
  }
}

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
      let code: Code;
      let error: string | string[];
      if (exception instanceof BadRequestException) {
        code = Code.InvalidArgument;
        error = (exception.getResponse() as any).message;
      } else {
        code = exception.code;
        error = exception.message;
      }

      return {
        code: status.INVALID_ARGUMENT,
        message: Code.InvalidArgument,
        details: {
          ok: false,
          code,
          error,
        },
      };
    });
  }
}

import { status as GrpcStatus } from '@grpc/grpc-js';

export type Response = {
  ok: boolean;
  code: string;
};

export type OkResponse<T = undefined> = Response & {
  ok: true;
  data: T;
};

export type ErrResponse = Response & {
  ok: false;
  error?: string | string[];
};

export type GrpcErrResponse = {
  code: GrpcStatus;
  message: string;
  details: ErrResponse;
};

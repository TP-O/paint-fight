import { status as GrpcStatus } from '@grpc/grpc-js';

type MyResponse = {
  ok: boolean;
  code: string;
};

export declare type OkResponse<T = undefined> = MyResponse & {
  ok: true;
  data: T;
};

export type ErrResponse = MyResponse & {
  ok: false;
  error?: string | string[];
};

export type GrpcErrResponse = {
  code: GrpcStatus;
  message: string;
  details: ErrResponse;
};

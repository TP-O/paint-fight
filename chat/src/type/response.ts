import { RequireOnlyOne } from './util';

export type Response = {
  ok: boolean;
  code: string;
};

export type OkResponse<T = undefined> = Response & {
  ok: true;
  data: T;
};

export type ErrResponse = Response &
  RequireOnlyOne<{
    error: string;
    fieldErrors: string[];
  }> & {
    ok: false;
  };

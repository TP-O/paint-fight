import { Code } from 'src/enum/code';

/**
 * The error is disclosed to clients.
 */
export class PublicError extends Error {
  public readonly name = 'PublicError';

  constructor(public readonly code: Code) {
    // TODO: convert code to message error
    super(undefined);
  }
}

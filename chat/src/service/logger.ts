import { ConsoleLogger, Injectable, Scope } from '@nestjs/common';
import { AppConfig } from 'src/config/app';
import { AppEnv } from 'src/enum/app';

@Injectable({
  scope: Scope.TRANSIENT,
})
export class LoggerService extends ConsoleLogger {
  /**
   * In development environment or not.
   */
  private readonly _isDev: boolean;

  constructor(config: AppConfig) {
    super();
    this._isDev = config.env === AppEnv.Development;
  }

  debug(message: any, context?: string) {
    if (this._isDev) {
      super.debug(message, context);
    }
  }
}

import { ConsoleLogger, Injectable, Scope } from '@nestjs/common';
import { resolve } from 'path';
import { AppConfig } from 'src/config/app';
import { AppEnv } from 'src/enum/app';
import { createLogger, format, transports } from 'winston';

// Consider whether to write log to file or not

const writer = createLogger({
  format: format.combine(
    format.errors({ stack: true }),
    format.timestamp(),
    format.prettyPrint(),
  ),
  transports: [
    new transports.File({
      filename: resolve(process.env.LOG_PATH || 'log', 'error.log'),
      level: 'error',
      maxsize: 10_000_000,
    }),
  ],
});

@Injectable({
  scope: Scope.TRANSIENT,
})
export class LoggerService extends ConsoleLogger {
  /**
   * Enable error logs writing.
   */
  private readonly _logError: boolean;

  /**
   * Logger instance.
   */
  private readonly _logWriter = writer;

  constructor(config: AppConfig) {
    super();
    this._logError = config.env === AppEnv.Production;
  }

  error(message: any, stack?: string, context?: string) {
    if (!this.isLevelEnabled('error')) {
      return;
    }

    super.error(message, stack, context);

    if (this._logError) {
      this._logWriter.error({
        context: this.context,
        message: message,
        stack,
      });
    }
  }
}

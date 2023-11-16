import logging
import structlog
from constant import DEV_ENV, PROD_ENV
from config import load_config

cfg = load_config()

def disable_development_log(logger, log_method, event_dict):
  if event_dict.get("env") == DEV_ENV and cfg["app"]["env"] == PROD_ENV:
    raise structlog.DropEvent
  else:
    return event_dict

structlog.configure(
  processors=[
    structlog.contextvars.merge_contextvars,
    disable_development_log,
    structlog.processors.add_log_level,
    structlog.processors.StackInfoRenderer(),
    structlog.dev.set_exc_info,
    structlog.processors.TimeStamper(fmt="%Y-%m-%d %H:%M:%S", utc=False),
    structlog.dev.ConsoleRenderer(),
  ],
  wrapper_class=structlog.make_filtering_bound_logger(logging.NOTSET),
  context_class=dict,
  logger_factory=structlog.PrintLoggerFactory(),
  cache_logger_on_first_use=False
)

logger = structlog.get_logger()

from typing import TypedDict
import yaml
from constant import CONFIG_FILE

class RabbitmqConfig(TypedDict):
  uri: str

class AppConfig(TypedDict):
  env: str

class EmailConfig(TypedDict):
  username: str
  hostname: str
  password: str

class Config(TypedDict):
  app: AppConfig
  rabbitmq: RabbitmqConfig
  email: EmailConfig

_cfg: Config = None

def load_config() -> Config:
  global _cfg
  if _cfg is None:
    with open(CONFIG_FILE, "r") as yamlfile:
      data = yaml.load(yamlfile, Loader=yaml.FullLoader)

      _cfg = Config()
      for key, value in data.items():
        _cfg[key] = value

  return _cfg

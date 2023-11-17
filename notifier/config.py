from typing import TypedDict
import yaml, os
from constant import CONFIG_FILE

class RabbitmqConfig(TypedDict):
  uri: str

class AppConfig(TypedDict):
  env: str
  host: str

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
    config_path = os.getenv("CONFIG_PATH")
    if config_path == None or config_path == "":
      config_path = CONFIG_FILE

    with open(config_path, "r") as yamlfile:
      data = yaml.load(yamlfile, Loader=yaml.FullLoader)

      _cfg = Config()
      for key, value in data.items():
        _cfg[key] = value

  return _cfg

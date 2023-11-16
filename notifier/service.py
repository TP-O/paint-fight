from typing import TypedDict, List
from enum import Enum
import apprise
from jinja2 import Environment, FileSystemLoader
from constant import TEMPLATE_DIR, DEV_ENV
from config import load_config
from logger import logger

class AppriseSerivce(Enum):
  EMAIL = 1

class NotifyTemplate(Enum):
  RESET_PASSWORD_EMAIL = 1
  VERIFY_EMAIL_EMAIL = 1

class Message(TypedDict):
  service_id: AppriseSerivce
  title: str
  template_id: NotifyTemplate
  params: dict
  data: dict

cfg = load_config()

def get_template_path(template: NotifyTemplate) -> str | None:
  if template == NotifyTemplate.VERIFY_EMAIL_EMAIL.value:
    return f"{TEMPLATE_DIR}/email/verify_email.html"
  elif template == NotifyTemplate.RESET_PASSWORD_EMAIL.value:
    return f"{TEMPLATE_DIR}/email/reset_password.html"
  else:
    return None

def handle_message(msg: Message):
  logger.info(f"Received message: {msg}", env=DEV_ENV)

  if msg["service_id"] == AppriseSerivce.EMAIL.value:
    send_email(msg["data"]["to"], msg["title"],msg["template_id"], msg["params"])

file_loader = FileSystemLoader('.')
env = Environment(loader=file_loader)

def send_email(to: str | List[str], title: str, template_id: NotifyTemplate, params: dict):
  apobj = apprise.Apprise()
  uri = f"mailto://{cfg["email"]["username"]}:{cfg["email"]["password"]}@{cfg["email"]["hostname"]}?from={cfg["email"]["username"]}@{cfg["email"]["hostname"]}"

  if isinstance(to, str):
    apobj.add(f"{uri}&to={to}")
  elif isinstance(to, List[str]):
    for t in to:
      apobj.add(f"{uri}&to={t}")
  else:
    return

  template_path = get_template_path(template_id)
  if template_id is None:
    return

  body = env.get_template(template_path).render(params)
  apobj.notify(
    title=title,
    body=body,
    body_format=apprise.NotifyFormat.HTML,
  )

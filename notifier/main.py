from ast import literal_eval
import pika, threading
from pika.adapters.blocking_connection import BlockingChannel
from config import load_config
from constant import QUEUE, EXCHANGE, ROUTING_KEY
from service import handle_message
from logger import logger
from keep_alive import keep_alive, start_http

def callback(ch: BlockingChannel, method, properties, body: bytes):
  try:
    handle_message(literal_eval(body.decode('utf-8')))
  except Exception as e:
    print(f"Error occurs: {e}")

if __name__ == "__main__":
  cfg = load_config()

  threading.Thread(target=start_http, name="Start HTTP server").start()
  threading.Thread(target=keep_alive, name="Keep service alive").start()

  params = pika.URLParameters(cfg["rabbitmq"]["uri"])
  params.socket_timeout = 15

  connection = pika.BlockingConnection(params)
  channel = connection.channel()

  channel.exchange_declare(exchange=EXCHANGE)
  channel.queue_declare(queue=QUEUE)
  channel.queue_bind(queue=QUEUE, exchange=EXCHANGE, routing_key=ROUTING_KEY)

  channel.basic_consume(queue=QUEUE, on_message_callback=callback, auto_ack=True)

  logger.info("Start consuming...")
  channel.start_consuming()

  connection.close()

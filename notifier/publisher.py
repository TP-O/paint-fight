import pika
from config import load_config
from constant import QUEUE, EXCHANGE, ROUTING_KEY
from logger import logger

if __name__ == "__main__":
  cfg = load_config()

  params = pika.URLParameters(cfg["rabbitmq"]["uri"])
  params.socket_timeout = 15

  connection = pika.BlockingConnection(params)
  channel = connection.channel()

  channel.exchange_declare(exchange=EXCHANGE)
  channel.queue_declare(queue=QUEUE)
  channel.queue_bind(queue=QUEUE, exchange=EXCHANGE, routing_key=ROUTING_KEY)

  channel.basic_publish(exchange=EXCHANGE, routing_key=ROUTING_KEY, body='''{
                        "service_id": 1,
                        "title": "Verify email",
                        "template_id": 1,
                        "params": {
                          "username": "Phong",
                          "link": "http://localhost.com"
                        },
                        "data": {
                          "to": "letranphong2k1@gmail.com"
                        }
                      }''')
  logger.info("Message sent to broker")
  connection.close()

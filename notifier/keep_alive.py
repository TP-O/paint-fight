from http.server import BaseHTTPRequestHandler, HTTPServer
import time, requests
from config import load_config

hostName = "0.0.0.0"
serverPort = 8080

cfg = load_config()

class MyServer(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()
        self.wfile.write(bytes("<p>Ok</p>", "utf-8"))

def start_http():
    webServer = HTTPServer((hostName, serverPort), MyServer)
    print("Server started http://%s:%s" % (hostName, serverPort))

    try:
        webServer.serve_forever()
    except KeyboardInterrupt:
        pass

    webServer.server_close()
    print("Server stopped.")

def keep_alive():
    while True:
        time.sleep(60 * 14)
        requests.get(cfg["app"]["host"])

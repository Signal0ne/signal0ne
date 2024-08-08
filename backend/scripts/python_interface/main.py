from dotenv import load_dotenv
import socket
import logging
import struct
import os
import json

from get_log_occurrences import log_occurrences

logging.basicConfig(
    filename="python_interface.log",
    level=logging.DEBUG,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)

logger = logging.getLogger(__name__)

bufferSizePrefix = 4

def main():

    load_dotenv(dotenv_path='.default.env')
    socket_path = os.getenv('IPC_SOCKET')

    try:
        os.unlink(socket_path)
    except OSError:
        if os.path.exists(socket_path):
            raise

    server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)

    server.bind(socket_path)

    server.listen(1)

    print('Server is listening for incoming connections...')
    connection, client_address = server.accept()

    try:

        logger.info('Connection from', str(connection).split(", ")[0][-4:])

        while True:
            header = connection.recv(bufferSizePrefix)
            if not header:
                break
            
            payloadSize = struct.unpack('>I', header)[0]
            payload = connection.recv(payloadSize)

            data = json.loads(payload)

            command = data["command"]
            params = data["params"]
            
            try:
                if command == "get_log_occurrences":
                    log_occurrences(params["collectedLogs"], params["comparedFields"])
            except Exception as e:
                    response = (1, str(e))
                    connection.sendall(response.encode())
        
            response = 0
            connection.sendall(response.encode())

    finally:

        connection.close()
        os.unlink(socket_path)

if __name__ == '__main__':
    main()
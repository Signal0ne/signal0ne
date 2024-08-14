from dotenv import load_dotenv
import socket
import logging
import struct
import os
import json
import traceback

from get_log_occurrences import log_occurrences
from correlate_ongoing_alerts import correlate_ongoing_alerts

logging.basicConfig(
    filename="python_interface.log",
    level=logging.DEBUG,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)

logger = logging.getLogger(__name__)

bufferSizePrefix = 4

def main():

    load_dotenv(dotenv_path='.default.env')
    print("loading...")
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

        print('Connection from', str(connection).split(", ")[0][-8:])
        payload = b''
        payloadBatchBuffer = float('-inf')

        while True:

            if payloadBatchBuffer < 0 :
                batchSizeHeader = connection.recv(bufferSizePrefix)
                if not batchSizeHeader:
                    break
            
                payloadSize = struct.unpack('>I', batchSizeHeader)[0]
                payloadBatchBuffer = float(payloadSize)

            payload += connection.recv(payloadSize)

            print("Payload size", len(payload), "Overall payload size",str(payloadBatchBuffer))
            if len(payload) >= int(payloadBatchBuffer):
                print("PYTHON SIZE: ",len(payload))
                data = json.loads(payload)
                command = data["command"]
                params = data["params"]
                print(command)
                print(params)
            
                try:
                    if command == "get_log_occurrences":
                        result = log_occurrences(params["collectedLogs"], params["comparedFields"])
                        parsedResult = json.dumps(result)
                        responseTemplate = json.dumps({"status":"0", "result":parsedResult})
                        response = len(responseTemplate).to_bytes(4, 'big') + bytes(responseTemplate, encoding="utf-8")
                        print("Success!!!")
                        connection.sendall(response)
                    if command == "correlate_ongoing_alerts":
                        result = correlate_ongoing_alerts(params["collectedEntities"], params["comparedFields"])
                        parsedResult = json.dumps(result)
                        responseTemplate = json.dumps({"status":"0", "result":parsedResult})
                        response = len(responseTemplate).to_bytes(4, 'big') + bytes(responseTemplate, encoding="utf-8")
                        print("Success!!!")
                        connection.sendall(response)    
                except Exception:
                        traceback.print_exc()
                        responseTemplate = json.dumps({"status":"1", "error":traceback.format_exc()})
                        response = len(responseTemplate).to_bytes(4, 'big') + bytes(responseTemplate, encoding="utf-8")
                        connection.sendall(response)

                payload = b''
                payloadBatchBuffer = float('-inf')
            responseTemplate = json.dumps({"status":"0"})
            response = len(responseTemplate).to_bytes(4, 'big') + bytes(responseTemplate, encoding="utf-8")
            connection.sendall(response)

    finally:

        connection.close()
        os.unlink(socket_path)

if __name__ == '__main__':
    main()
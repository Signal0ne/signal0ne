import socket
import os

socket_path = '/net/socket'

def main():

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

        print('Connection from', str(connection).split(", ")[0][-4:])

        while True:
            data = connection.recv(1024)
            if not data:
                break
            print(data.decode())

            response = 'Hi Go! I am python!'
            connection.sendall(response.encode())

    finally:

        connection.close()
        os.unlink(socket_path)

if __name__ == '__main__':
    main()
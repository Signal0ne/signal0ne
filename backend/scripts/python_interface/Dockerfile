FROM python:3.9-slim

WORKDIR /app
COPY ./scripts/python_interface/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY ./scripts/python_interface ./python_interface/
COPY .default.env .default.env

# Start both applications
CMD ["python", "-u", "./python_interface/main.py"]

FROM python:3.9-slim

WORKDIR /app
COPY ./scripts/python_interface/requirements.txt .
COPY ./scripts/healthcheck.sh .
RUN pip install --no-cache-dir -r requirements.txt && \
apt-get update && apt-get install -y netcat-openbsd && rm -rf /var/lib/apt/lists/* && \
chmod +x /app/healthcheck.sh
COPY ./scripts/python_interface ./python_interface/
COPY .default.env .default.env

HEALTHCHECK --interval=15s --timeout=5s --retries=3 CMD /app/healthcheck.sh || exit 1

# Start both applications
CMD ["python", "-u", "./python_interface/main.py"]
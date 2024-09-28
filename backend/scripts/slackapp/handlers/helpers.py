import os
import requests
from slack_sdk import WebClient
from slack_sdk.errors import SlackApiError
from dotenv import load_dotenv

load_dotenv()
BACKEND_URL = os.getenv("BACKEND_URL")
BACKEND_AUTH_TOKEN = os.getenv("BACKEND_AUTH_TOKEN")
    
def get_enriched_alert_by_id(alert_id):
    url = f"{BACKEND_URL}/api/alert/{alert_id}/correlations"
    headers = {
        "Authorization ": f"Bearer {BACKEND_AUTH_TOKEN}"
    }
    response = requests.get(url, headers=headers)
    return response.json()

def create_incident(incident_destination, alert_ids):
    url = f"{BACKEND_URL}/api/incident"
    headers = {
        "Authorization ": f"Bearer {BACKEND_AUTH_TOKEN}"
    }
    data = {
        "destination": incident_destination,
        "alert_ids": alert_ids
    }
    response = requests.post(url, headers=headers, json=data)
    return response.json()
    

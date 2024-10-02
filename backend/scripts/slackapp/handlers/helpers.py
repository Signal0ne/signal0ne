import os
import json
import requests
from slack_sdk import WebClient
from slack_sdk.errors import SlackApiError
from dotenv import load_dotenv

load_dotenv()
BACKEND_URL = os.getenv("BACKEND_URL")
BACKEND_AUTH_TOKEN = os.getenv("BACKEND_AUTH_TOKEN")
NAMESPACE_ID = os.getenv("ORG_NAMESPACE_ID")
    
def get_enriched_alert_by_id(alert_id):
    url = f"{BACKEND_URL}/api/{NAMESPACE_ID}/alert/{alert_id}"
    headers = {
        "Authorization": f"Bearer {BACKEND_AUTH_TOKEN}"
    }
    response = requests.get(url, headers=headers)
    if response.status_code != 200:
        raise requests.RequestException(f"""An error from 
                                        api(code: {response.status_code}) 
                                        occurred: {response.json()}""")
    return response.json()

def create_incident(incident_destination: str, alert_ids: list):
    url = f"{BACKEND_URL}/api/{NAMESPACE_ID}/incident"
    headers = {
        "Authorization": f"Bearer {BACKEND_AUTH_TOKEN}"
    }
    data = {
        "integration": incident_destination,
        "baseAlertId": alert_ids[0],
    }
    if len(alert_ids) > 1:
        data["manuallyCorrelatedAlertIds"] = alert_ids[1:]
        
    response = requests.post(url, headers=headers, json=data)
    if response.status_code != 200:
        raise requests.RequestException(f"""An error from 
                                        api(code: {response.status_code}) 
                                        occurred: {response.json()}""")
    return response.json()
    

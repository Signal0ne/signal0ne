import os
import requests
from slack_sdk import WebClient
from slack_sdk.errors import SlackApiError
from dotenv import load_dotenv
    
def get_enriched_alert_by_id(alert_id):
    load_dotenv()
    host = os.getenv("BACKEND_HOST")
    url = f"{url}/api/alert/{alert_id}/correlations"
    response = requests.get(url)
    return response.json()

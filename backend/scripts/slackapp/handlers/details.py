from slack_bolt import Ack, Respond
import requests
from handlers.helpers import get_enriched_alert_by_id

def handle(ack: Ack, respond: Respond, command):
    ack()
    blocks = []
    
    # Extract the alert_id from the command text
    alert_id = command['text'].strip()

    if not alert_id:
        respond("Please provide a valid alert ID.")
        return
    
    try:
        data = get_enriched_alert_by_id(alert_id)
        for k,v in dict(data[0]).items():
            blocks.append({
                "type": "section",
                "text": {
                    "type": "mrkdwn",
                    "text": f"*{k}:*\n```{v}```"
                }
            })
    except requests.RequestException as e:
        print(f"An error occurred: {e}")


    
    respond({
        "response_type": "in_channel",
        "blocks": blocks
    })

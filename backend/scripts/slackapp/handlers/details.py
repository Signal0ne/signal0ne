from slack_bolt import Ack, Respond
import requests
from handlers.helpers import get_enriched_alert_by_id

def handle(ack: Ack, respond: Respond, command):
    ack()
    blocks = []
    
    command_params = command['text'].split(" ")

    if len(command_params) < 2:
        respond("Please use command in the following format: `/details <alert_id> <[]tags...>`")
        return
    
    alert_id = command_params[0]
    tags = command_params[1:]
    
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

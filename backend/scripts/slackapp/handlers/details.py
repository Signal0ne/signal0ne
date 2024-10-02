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
    
    print(command_params)
    
    alert_id = command_params[0]
    tags = command_params[1:]
    
    try:
        data = get_enriched_alert_by_id(alert_id)
        for k,v in dict(data)["additionalContext"].items():
            if not v:
                continue
            for item in v:
                if any(tag in list(item["tags"]) for tag in tags):
                    blocks.append({
                        "type": "section",
                        "text": {
                            "type": "mrkdwn",
                            "text": f"*{k}:*\n```{item}```"
                        }
                    })
    except requests.RequestException as e:
        print(f"An error occurred: {e}")
        respond("An error occurred while fetching the alert details. Please try again later.")


    
    respond({
        "response_type": "in_channel",
        "blocks": blocks
    })

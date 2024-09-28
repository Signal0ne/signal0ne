from slack_bolt import Ack, Respond
import requests
from handlers.helpers import get_enriched_alert_by_id, create_incident

def handle_create(ack: Ack, respond: Respond, command):
    ack()
    blocks = []

    command_params = command['text'].split(" ")

    if len(command_params) < 2:
        respond("Please use command in the following format: `/create-incident <incident_destination> <[]alert_ids...>`")
        return
    
    incident_destination = command_params[0]
    alert_ids = command_params[1:]

    try:
        incident = create_incident(incident_destination, alert_ids)
        blocks.append({
            "type": "section",
            "text": {
                "type": "mrkdwn",
                "text": f"Link to incident: https://{incident_destination}/incident/{incident['id']}"
            }
        })
    except requests.RequestException as e:
        print(f"An error occurred: {e}")
        respond("Failed to create incident")

    respond({
        "response_type": "in_channel",
        "blocks": blocks
    })
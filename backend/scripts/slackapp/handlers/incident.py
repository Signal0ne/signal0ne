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

    print("ALERT IDS", alert_ids)
    print("INCIDENT DESTINATION", incident_destination)

    try:
        incident_creation_response = create_incident(incident_destination, alert_ids)
        incident = incident_creation_response[0]
        blocks.append({
            "type": "section",
            "text": {
                "type": "mrkdwn",
                "text": f"*New incident created*"
            }
        })
        blocks.append({
            "type": "section",
            "text": {
                "type": "mrkdwn",
                "text": f"""*Id:* {incident['id']}\n
                        *Name:* {incident['name']}\n
                        *Status:* {incident['status']}\n
                        *Severity:* {incident['severity']}\n
                        *Link:* {incident['url']}"""
            }
        })
    except requests.RequestException as e:
        print(f"An error occurred: {e}")
        respond("Failed to create incident")

    respond({
        "response_type": "in_channel",
        "blocks": blocks
    })
    
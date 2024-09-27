from flask import Blueprint, request, jsonify
from slack_sdk import WebClient
from slack_sdk.errors import SlackApiError
from dotenv import load_dotenv
import os

# Load environment variables
load_dotenv()

workspace_access_token = os.getenv("SLACK_BOT_TOKEN")

post_message_route = Blueprint('post_message_route', __name__)

@post_message_route.route("/post_message", methods=["POST"])
def post_message():
    data = request.json
    # Schema:
    # channelName: string
    # Title: string
    # Data: dict
    # Id: string
    channel_name = data.get("channelName")
    message_data = data.get("data")
    alert_id = data.get("id")
    title = data.get("title") 

    # Ensure channel_name and message_data are provided
    if not channel_name or not message_data:
        return jsonify({"error": "Missing channelName or data"}), 400

    client = WebClient(token=workspace_access_token)

    try:
        # Creating blocks for the message
        blocks = []

        if title:
            blocks.append({
                "type": "header",
                "text": {
                    "type": "plain_text",
                    "text": title,
                    "emoji": True
                }
            })

        # Add the alert ID as the first block if it is provided
        if alert_id:
            blocks.append({
                "type": "section",
                "text": {
                    "type": "mrkdwn",
                    "text": f"*alert id:*\n```{alert_id}```"
                }
            })

        # Add other message data blocks
        for key, value in message_data.items():
            blocks.append({
                "type": "section",
                "text": {
                    "type": "mrkdwn",
                    "text": f"*{key}:*\n```{value}```"
                }
            })

        # Post the message to the Slack channel
        response = client.chat_postMessage(
            channel=channel_name,
            blocks=blocks,
            text="New Alert"
        )
        return jsonify({"message": "Message posted successfully", "response": response.data}), 200
    except SlackApiError as e:
        return jsonify({"error": str(e)}), 500
    
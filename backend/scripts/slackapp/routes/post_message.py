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

    channel_name = data.get("channelName")

    title = data.get("title") 
    alert_id = data.get("id")
    source = data.get("source")
    time = data.get("time")
    status = data.get("status")
    original_url = data.get("originalUrl")
    tags = data.get("tags")

    try:
        additional_context = dict(data.get("additionalContext"))
    except:
        additional_context = None

    if not channel_name:
        return jsonify({"error": "Missing channelName"}), 400

    client = WebClient(token=workspace_access_token)

    try:
        blocks = []

        if alert_id:
            blocks.append({
                "type": "section",
                "text": {
                    "type": "mrkdwn",
                    "text": f"""
                    
                    *{title}*\n*ID:* {alert_id}
*Source:* {source}\n*Time:* {time}\n*Status:* {status}
*Original URL:* {original_url}
*Produced output tags:* {", ".join(tags)}"""
                }
            })

        if additional_context:
            for ctxKey, ctx in additional_context.items():
                blocks.append({
                    "type": "section",
                    "text": {
                        "type": "mrkdwn",
                        "text": f"*{ctxKey}:*```{ctx}```\n"
                    }
                })

        response = client.chat_postMessage(
            channel=channel_name,
            blocks=blocks,
            text="New Alert"
        )
        return jsonify({"message": "Message posted successfully", "response": response.data}), 200
    except SlackApiError as e:
        return jsonify({"error": str(e)}), 500
    
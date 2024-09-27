from flask import Blueprint, request, jsonify
from slack_sdk import WebClient
from slack_sdk.errors import SlackApiError
from dotenv import load_dotenv
import os

# Load environment variables
load_dotenv()

workspace_access_token = os.getenv("SLACK_BOT_TOKEN")

# Initialize Slack WebClient
slack_client = WebClient(token=workspace_access_token)

create_channel_route = Blueprint('create_channel_route', __name__)

@create_channel_route.route("/create_channel", methods=["POST"])
def create_channel():
    data = request.json
    channel_name = data.get("channel_name")
    is_private = data.get("is_private", False)  # Default to public channel if not specified

    if not channel_name:
        return jsonify({"error": "Channel name is required"}), 400

    try:
        # Create the channel
        response = slack_client.conversations_create(
            name=channel_name,
            is_private=is_private
        )
        channel_id = response['channel']['id']

        return jsonify({
            "message": "Channel created successfully",
            "channel_id": channel_id,
            "channel_name": channel_name
        }), 200

    except SlackApiError as e:
        error_message = e.response['error']
        return jsonify({"error": f"Failed to create channel: {error_message}"}), e.response.status_code

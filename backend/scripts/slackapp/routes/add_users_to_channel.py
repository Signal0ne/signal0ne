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

add_users_to_channel_route = Blueprint('add_users_to_channel_route', __name__)

@add_users_to_channel_route.route("/add_users_to_channel", methods=["POST"])
def add_users_to_channel():
    data = request.json
    channels = slack_client.conversations_list()
    channel_name = data.get("channel_name")
    for channel in channels['channels']:
        if channel['name'] == channel_name:
            channel_id = channel['id']
            break
    user_handles = data.get("user_handles")  # List of user handles (e.g., ["@user1", "@user2"])

    if not channel_id:
        return jsonify({"error": "Channel ID is required"}), 400

    if not user_handles or not isinstance(user_handles, list):
        return jsonify({"error": "A list of user handles is required"}), 400

    user_ids = []
    for handle in user_handles:
        try:
            # Get the user's ID by their handle
            user_info = slack_client.users_lookupByEmail(email=handle)
            user_id = user_info['user']['id']
            user_ids.append(user_id)
        except SlackApiError as e:
            error_message = e.response['error']
            return jsonify({"error": f"Failed to find user {handle}: {error_message}"}), e.response.status_code

    try:
        # Invite users to the channel
        response = slack_client.conversations_invite(
            channel=channel_id,
            users=",".join(user_ids)
        )
        return jsonify({
            "message": "Users added to channel successfully",
            "channel_id": channel_id,
            "user_ids": user_ids
        }), 200

    except SlackApiError as e:
        error_message = e.response['error']
        return jsonify({"error": f"Failed to add users to channel: {error_message}"}), e.response.status_code

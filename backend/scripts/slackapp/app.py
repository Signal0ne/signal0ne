import os
from flask import Flask, request
from slack_bolt import App
from slack_bolt.adapter.flask import SlackRequestHandler
from dotenv import load_dotenv
from handlers import details
from routes.post_message import post_message_route
from routes.add_users_to_channel import add_users_to_channel_route
from routes.create_channel import create_channel_route

# Load environment variables
load_dotenv()

token = os.getenv("SLACK_BOT_TOKEN")
signing_secret = os.getenv("SLACK_SIGNING_SECRET")

print("TOKEN", token)

# Initialize Slack app
slack_app = App(
    token=token,
    signing_secret=signing_secret
)

# Used just in cloud
slack_app.command("/details")(details.handle)

# Initialize Flask app
flask_app = Flask(__name__)
handler = SlackRequestHandler(slack_app)

@flask_app.route("/slack/events", methods=["POST"])
def slack_events():
    return handler.handle(request)

# Used in both cloud and local selfhosted
flask_app.register_blueprint(post_message_route, url_prefix='/api')
flask_app.register_blueprint(create_channel_route, url_prefix='/api')
flask_app.register_blueprint(add_users_to_channel_route, url_prefix='/api')

if __name__ == "__main__":
    flask_app.run(host='0.0.0.0',port=3000)

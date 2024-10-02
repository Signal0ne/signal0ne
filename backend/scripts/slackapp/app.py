import os
from flask import Flask, request
from slack_bolt import App
from slack_bolt.adapter.flask import SlackRequestHandler
from dotenv import load_dotenv
from handlers import details, incident
from routes.post_message import post_message_route
from routes.add_users_to_channel import add_users_to_channel_route
from routes.create_channel import create_channel_route

load_dotenv()

token = os.getenv("SLACK_BOT_TOKEN")
signing_secret = os.getenv("SLACK_SIGNING_SECRET")

slack_app = App(
    token=token,
    signing_secret=signing_secret
)

slack_app.command("/details")(details.handle)
slack_app.command("/create-incident")(incident.handle_create)

flask_app = Flask(__name__)
handler = SlackRequestHandler(slack_app)

@flask_app.route("/slack/events", methods=["POST"])
def slack_events():
    return handler.handle(request)

flask_app.register_blueprint(post_message_route, url_prefix='/api')
flask_app.register_blueprint(create_channel_route, url_prefix='/api')
flask_app.register_blueprint(add_users_to_channel_route, url_prefix='/api')

if __name__ == "__main__":
    flask_app.run(host='0.0.0.0',port=3000)

"""Server for flask app."""

from flask import Flask, render_template

app = Flask(__name__)


@app.route("/")
@app.route("/index")
def index():
    """Home page for app."""
    user = {'username': 'Sam'}
    return render_template('index.html', title='Home', user=user)

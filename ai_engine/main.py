from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route("/generate", methods=["POST"])
def generate():
    


if __name__ == "__main__":
    app.run(port=5000, debug=True)
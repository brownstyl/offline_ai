from flask import Flask, request, jsonify

app = Flask(__name__)

@app.route("/generate", methods=["POST"])
def generate():
    data = request.get_json()
    incoming_message = data.get("message")

    print(f"Message from Go: {incoming_message}")

    ai_reply = "Ai hasn't been integrated yet, please chill..."
    return jsonify({"reply": ai_reply})



if __name__ == "__main__":
    app.run(port=5000, debug=True)
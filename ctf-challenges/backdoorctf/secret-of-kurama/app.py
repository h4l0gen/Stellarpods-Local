from flask import Flask, render_template, request, redirect, url_for, make_response
import jwt

app = Flask(__name__)
app.secret_key = "vouvou"

USERNAME = "Naruto"
PASSWORD = "Chakra"

SECRET_KEY = "minato"

@app.route("/", methods=["GET", "POST"])
def login():
    if request.method == "POST":
        username = request.form["username"]
        password = request.form["password"]
        if username == USERNAME and password == PASSWORD:
            jwt_token = jwt.encode({"username": username, "role": "shinobi"}, SECRET_KEY, algorithm="HS256")

            response = make_response(redirect(url_for("war")))
            response.set_cookie('jwt_token', jwt_token, httponly=True)
            return response
        else:
            return "Invalid credentials. Try Again."
    return render_template("login.html")

@app.route("/war")
def war():
    token = request.cookies.get('jwt_token')
    if token:
        try:
            return render_template("war.html")
        except jwt.ExpiredSignatureError:
            return "Token has expired. Please log in again."
        except jwt.InvalidTokenError:
            return "Invalid token. Please log in again."
    else:
        return redirect(url_for("login"))

@app.route("/secret_of_Naruto")
def fun():
    return render_template("fun.html")

@app.route("/secret_of_Kurama")
def flag():
    token = request.cookies.get('jwt_token')
    if token:
        try:
            decoded_token = jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
            if decoded_token.get("role") == "NineTails":
                return render_template("success.html")
            else:
                return render_template("wrong.html")
        except jwt.ExpiredSignatureError:
            return "Token has expired. Please log in again."
        except jwt.InvalidTokenError:
            return render_template("wrong.html")
    else:
        return redirect(url_for("login"))

if __name__ == "__main__":
    app.run(debug=False, host='0.0.0.0', port=4054)

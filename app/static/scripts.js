// connect to WebSocket
const ws = new WebSocket("ws://localhost:8080/ws")

ws.onopen = () => {
    console.log("Connected!");
    document.getElementById("messages").innerHTML += "<p>Connected to server!</p>";
};

ws.onmessage = e => {
    document.querySelector("#messages").innerHTML += "<p>Server says: " + e.data + "</p>";
};

function sendMessage() {
    const input = document.querySelector("#messageInput");
    ws.send(input.value);
    document.querySelector("#messages").innerHTML += "<p>You: " + input.value + "</p>";
    input.value = "";
}
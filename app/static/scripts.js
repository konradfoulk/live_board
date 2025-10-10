// connect to WebSocket
const ws = new WebSocket("ws://localhost:8080/ws")

ws.onopen = () => {
    console.log("Connected!");
    document.getElementById("messages").innerHTML += "Connected to server!\n";
};

ws.onmessage = e => {
    document.querySelector("#messages").innerHTML += "Serer says" + e.data + "\n";
};

function sendMessage() {
    const input = document.querySelector("#messageInput");
    ws.send(input.value);
    document.querySelector("#messages").innerHTML += "You:" + input.value + "\n";
    input.value = "";
}
// connect to websocket
let roomName = "general"
const ws = new WebSocket(`ws://localhost:8080/ws?room=${roomName}`)

ws.onopen = () => {
    console.log(`client connected to ${roomName}`)
    document.getElementById(`${roomName}`).innerHTML += `<p>Connected to ${roomName}</p>`
}

ws.onmessage = e => {
    document.getElementById(`${roomName}`).innerHTML += `<p>${e.data}</p>`
}

document.querySelector(".sendBtn").addEventListener("click", () => {
    const input = document.querySelector(".messageInput");
    ws.send(input.value);
    input.value = "";
})
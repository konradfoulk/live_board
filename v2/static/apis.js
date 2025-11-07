let username = ""
let roomName = "general"


function connectToChat(username, roomName) {
    const ws = new WebSocket(`ws://localhost:8080/ws?room=${roomName}&username=${username}`)

    ws.onopen = () => {
        console.log(`${username} connected to ${roomName}`)
        document.getElementById(`${roomName}`).innerHTML += `<p>${username} connected to ${roomName}</p>`
    }

    ws.onmessage = e => {
        document.getElementById(`${roomName}`).innerHTML += `<p>${e.data}</p>`
    }

    // need a websocket function to update the ui when a message comes saying the state has changed
    // could be a websocket message for the browser instructing it to do a http get

    document.querySelector(".sendBtn").addEventListener("click", () => {
        const input = document.querySelector(".messageInput");
        ws.send(input.value);
        input.value = "";
    })
}

document.querySelector("#joinBtn").addEventListener("click", () => {
    username = document.querySelector("#usernameInput").value

    if (username === "") {
        alert("Please enter a username")
        return;
    }

    // hide modal
    document.querySelector("#usernameModal").style.display = "none";

    // connect to WebSocket
    connectToChat(username, roomName)

    // could have add event listeners function here so that you can't edit anything until authenticated
})
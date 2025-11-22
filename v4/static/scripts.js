let ws = null
let currentRoom = ""

async function createRoom(roomName) {
    // call server API
    const response = await fetch("/api/rooms", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json"
        },
        body: JSON.stringify({ name: roomName })
    })
    const data = await response.json()

    if (response.ok) {
        console.log(`created room ${data.name}: ${response.status}`)
    } else {
        console.log(`could not create room ${data.name}: ${response.status}`)
    }
}

async function deleteRoom(event) {
    const roomName = event.target.dataset.room

    // call server API
    const response = await fetch(`/api/rooms/${roomName}`, { method: "POST" })
    const data = await response.json()

    if (response.ok) {
        console.log(`deleted room ${data.name}: ${response.status}`)
    } else {
        console.log(`could not delete room ${data.name}: ${response.status}`)
    }    
}

function joinRoom(event) {
    const roomName = event.target.dataset.room

    currentRoom = roomName
    console.log(currentRoom)

    msg = {
        type: "join_room",
        room: roomName
    }
    ws.send(JSON.stringify(msg))
}

// establishes websocket connection and recieving ports
function connectToChat(username) {
    ws = new WebSocket(`ws://localhost:8080/ws?username=${username}`)

    ws.onopen = () => {
        console.log(`${username} connected to server`)
    }

    ws.onmessage = e => {
        const msg = JSON.parse(e.data)

        switch (msg.type) {
            case "init_rooms":
                if (msg.rooms) {
                    msg.rooms.reverse().forEach(room => {
                        const newRoom = newRoomBtn(room)
                        document.querySelector("#roomBtns").prepend(newRoom)
                    })
                    // current room = msg.rooms[0]
                }
                break
                // build room buttons and click the first one (if there is one) [joining the "default" room on load]
            case "create_room":
                const newRoom = newRoomBtn(msg.room)
                document.querySelector("#roomBtns").appendChild(newRoom)
                break
            case "delete_room":
                // may need to change this if chats use query selector
                document.querySelectorAll(`[data-room="${msg.room}"`).forEach(element => {
                    element.remove()
                })
                break
            case "chat_message":
                if (msg.room != currentRoom) {
                    // don't send message
                }
        }
    }


    // build websocket events
    // receive a chat message (for what room?)
    // receive a create room update - done
    // receive a delete room update - done
    // receive intial room state - done
    // receive initial chat state
}
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
    if (roomName == currentRoom) {
        return
    }

    currentRoom = roomName
    console.log(currentRoom)

    // delete old room chat
    document.querySelector(".roomChat").remove()

    // make new room chat and append it to roomChats
    const roomChat = document.createElement("div")
    roomChat.className = "roomChat"
    roomChat.setAttribute("data-room",roomName)
    document.querySelector("#roomChats").append(roomChat)

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
            case "create_room":
                const newRoom = newRoomBtn(msg.room)
                document.querySelector("#roomBtns").append(newRoom)
                break
            case "delete_room":
                // may need to change this if chats use query selector
                if (msg.room === currentRoom) {
                    currentRoom = ""
                }
                document.querySelectorAll(`[data-room="${msg.room}"`).forEach(element => {
                    element.remove()
                })
                break
            case "init_rooms":
                if (msg.rooms) {
                    msg.rooms.reverse().forEach(room => {
                        const newRoom = newRoomBtn(room)
                        document.querySelector("#roomBtns").prepend(newRoom)
                    })
                    // click #roomBtns first child
                }
                break
                // build room buttons and click the first one (if there is one) [joining the "default" room on load]
            case "message":
                if (msg.room == currentRoom) {
                    switch (msg.messageType) {
                        case "join_message":
                            break
                        case "leave_message":
                            break
                        case "chat_message":
                            break
                        case "init_chat":
                            break
                    }
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
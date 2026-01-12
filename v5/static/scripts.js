let ws = null
let currentRoom = ""

function chatAutoScroll() {
    const chat = document.querySelector("#roomChats")
    chat.scrollTop = chat.scrollHeight
}

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

    // make new room chat and append it to roomChats
    const roomChat = document.createElement("div")
    roomChat.className = "roomChat"
    roomChat.setAttribute("data-room", roomName)
    document.querySelector("#roomChats").append(roomChat)

    if (currentRoom === "") {
        currentRoom = roomName

        // enable input bar
        document.querySelectorAll("#messageInput [disabled]").forEach(element => {
            element.disabled = false
        })
    } else {
        currentRoom = roomName

        // delete old roomChat
        document.querySelector(".roomChat").remove() 
    }

    msg = {
        type: "join_room",
        room: roomName
    }
    ws.send(JSON.stringify(msg))

    chatAutoScroll()
}

// establishes websocket connection and recieving ports
function connectToChat(username, password) {
    return new Promise((resolve, reject) => {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        ws = new WebSocket(`${protocol}//${window.location.host}/ws?username=${username}&password=${password}`)

        ws.onopen = () => {
            console.log(`${username} connected to server`)
            resolve()
        }

        ws.onerror = () => {
            console.log("connection failed")
            reject()
        }

        ws.onmessage = e => {
            const msg = JSON.parse(e.data)
            const roomBtns = document.querySelector("#roomBtns")

            switch (msg.type) {
                case "create_room":
                    const newRoom = newRoomBtn(msg.room)
                    roomBtns.append(newRoom)
                    break
                case "delete_room":
                    if (msg.room === currentRoom) {
                        document.querySelectorAll("#messageInput *").forEach(element => {
                            element.disabled = true
                        })
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
                            roomBtns.prepend(newRoom)
                        })
                    }
                    break
                case "message":
                    if (msg.room == currentRoom) {
                        switch (msg.messageType) {
                            case "join_message":
                                document.querySelector(`.roomChat[data-room="${msg.room}"]`).innerHTML += `<p>${msg.username} joined ${msg.room}<p>`
                                chatAutoScroll()
                                break
                            case "leave_message":
                                document.querySelector(`.roomChat[data-room="${msg.room}"]`).innerHTML += `<p>${msg.username} left ${msg.room}<p>`
                                chatAutoScroll()
                                break
                            case "chat_message":
                                document.querySelector(`.roomChat[data-room="${msg.room}"]`).innerHTML += `<p>${msg.username}: ${msg.content}<p>`
                                chatAutoScroll()
                                break
                            case "init_chat":
                                msg.messages.forEach(i => {
                                    const p = document.createElement('p')
                                    p.innerHTML = `${i.username}: ${i.content}`

                                    document.querySelector(`.roomChat[data-room="${msg.room}"]`).prepend(p)
                                    chatAutoScroll()
                                })
                                break
                        }
                    }
                    break
                case "user_count":
                    document.querySelector("#count").textContent = msg.userCount
                    break
            }
        }
    })
}
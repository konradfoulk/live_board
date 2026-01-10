const joinModal =  document.querySelector("#joinModal")
const newRoomModal = document.querySelector("#newRoomModal")

function newRoomBtn(roomName) {
    const roomBtnContainer = document.createElement("div")
    const roomBtn = document.createElement("button")
    const deleteBtn = document.createElement("button")

    roomBtnContainer.className = "roomBtnContainer"
    roomBtn.className = "roomBtn"
    deleteBtn.className = "deleteBtn"

    roomBtn.textContent = roomName
    deleteBtn.textContent = "delete"
    
    const elements = [roomBtnContainer, roomBtn, deleteBtn]
    elements.forEach(element => {
        element.setAttribute("data-room", roomName)
    })

    roomBtnContainer.append(roomBtn)
    roomBtnContainer.append(deleteBtn)

    roomBtn.addEventListener("click", joinRoom)
    deleteBtn.addEventListener("click", deleteRoom)

    return roomBtnContainer

    // create the element

    // add an event listener for when it is clicked
    // create a new chat element with the proper room attribute (give a place for incoming websocket message to go)
    // tell the ws endpoint to change rooms (get initial state from response and populate it in the new div)
    // delete old room div and make new room div visible (done immediately after resopnse is sent so while waiting to load, users waits at an empty screen and if they try to type, they will be in the new room so no weird races)
}

// deactivate new room form if clicked off
function clickOff(e) {
    if (!newRoomModal.contains(e.target))
        newRoomModal.style.display = "none"
        document.removeEventListener("click", clickOff)
}

function addEventListeners() {
    // inputBar
    document.querySelector("#messageInput").addEventListener("submit", e => {
        const message = e.target.elements.message

        msg = {
            type: "message",
            room: currentRoom,
            content: message.value
        }
        ws.send(JSON.stringify(msg))

        message.value = ""
    })

    // display new room form
    document.querySelector("#newRoom").addEventListener("click", e => {
        e.stopPropagation() // prevent click from bubbling to DOM and activating clickOff

        newRoomModal.style.display = ""
        document.addEventListener("click", clickOff) // form deactivates if clicked off
        newRoomModal.elements.roomName.focus()
    })

    // get room name and create room
    newRoomModal.addEventListener("submit", e => {
        const roomName = e.target.elements.roomName
        createRoom(roomName.value)

        roomName.value = ""
        newRoomModal.style.display = "none"
        document.removeEventListener("click", clickOff)
    })
}

document.querySelectorAll("form").forEach(element => {
    element.addEventListener("submit", e => {
        e.preventDefault() // stop page reload on form submissions
    })
})

// get username and start app
joinModal.addEventListener("submit", e => {
    const username = e.target.elements.username.value
    if (username === "") {
        return
    }
    connectToChat(username)
    addEventListeners()

    joinModal.style.display = "none"
})

joinModal.elements.username.focus()
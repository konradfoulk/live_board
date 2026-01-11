const joinModal =  document.querySelector("#joinModal")
const newRoomModal = document.querySelector("#newRoomModal")
const counter = document.querySelector("#counter")

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
joinModal.addEventListener("submit", async e => {
    const username = e.target.elements.username.value
    const password = e.target.elements.password
    if (username === "" || password.value === "") {
        return
    }

    try {
        await connectToChat(username, password.value)
        addEventListeners()
        joinModal.style.display = "none"
        counter.style.display = ""
    } catch {
        // auth failed or error
        password.value = ""
        return
    }
})

joinModal.elements.username.focus()
const joinModal =  document.querySelector("#joinModal")
const newRoom = document.querySelector("#newRoom")
const newRoomModal = document.querySelector("#newRoomModal")

// deactivate new room form if clicked off
function clickOff(e) {
    if (!newRoomModal.contains(e.target))
        newRoomModal.style.display = "none"
        document.removeEventListener("click", clickOff)
        
}

// get username and start app
joinModal.addEventListener("submit", e => {
    e.preventDefault() // stop page reload

    const username = e.target.elements.username.value
    if (username === "") {
        return
    }
    connectToChat(username)

    // enable app and remove modal
    document.querySelectorAll("[disabled]").forEach(element => {
        element.classList.add("active")
        element.disabled = false
    })
    joinModal.style.display = "none"
})

// display new room form
newRoom.addEventListener("click", e => {
    e.stopPropagation() // prevent click from bubbling to DOM and activating clickOff

    newRoomModal.style.display = ""
    document.addEventListener("click", clickOff) // form deactivates if clicked off
})

// get room name and create room
newRoomModal.addEventListener("submit", e => {
    e.preventDefault() // stop page reload

    const roomName = e.target.elements.roomName
    createRoom(roomName.value)

    roomName.value = ""
    newRoomModal.style.display = "none"
    document.removeEventListener("click", clickOff)
})
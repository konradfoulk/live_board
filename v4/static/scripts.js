function createRoom(roomName) {
    // send to backend
    // backend will broadcast update, that is when button will be built with event listener function to handle other logic
}

function newRoomBtn() {
    // create the element

    // add an event listener for when it is clicked
    // create a new chat element with the proper room attribute (give a place for incoming websocket message to go)
    // tell the ws endpoint to change rooms (get initial state from response and populate it in the new div)
    // delete old room div and make new room div visible (done immediately after resopnse is sent so while waiting to load, users waits at an empty screen and if they try to type, they will be in the new room so no weird races)
}

// establishes websocket connection
// builds websocket event listeners
function connectToChat(username) {
    const ws = new WebSocket(`ws://localhost:8080/ws?username=${username}`) // user automatically added to general  or first room on join
    // check this ^, adding them to general may not be necessary (could just map button, button may not even be necessary: could start standard create room procedure with the first user on the backend when they join in handlews)

    ws.onmessage = e => {
        const msg = JSON.parse(e.data)

        switch (msg.type) {
            case "create_room":

        }
    }


    // build websocket events
    // receive a chat message (for what room?)
    // receive a create room update
    // receive a delete room update
    // receive intial room state (and what room it's going to)
    // receive initial chat state
}
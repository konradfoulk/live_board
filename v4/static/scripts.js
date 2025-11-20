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
    const data = await response.json();

    if (response.ok) {
        console.log(`created room ${data.name}: ${response.status}`)
    } else {
        console.log(`could not create room ${data.name}: ${response.status}`)
    }
}

// establishes websocket connection
function connectToChat(username) {
    const ws = new WebSocket(`ws://localhost:8080/ws?username=${username}`)

    ws.onopen = () => {
        console.log(`${username} connected to server`)
    }

    ws.onmessage = e => {
        const msg = JSON.parse(e.data)

        switch (msg.type) {
            case "create_room":
                console.log(msg.room)
                break
            case "init_rooms":
                console.log(msg.rooms)
                break
                // build room buttons and click the first one (if there is one) [joining the "default" room on load]
        }
    }


    // build websocket events
    // receive a chat message (for what room?)
    // receive a create room update
    // receive a delete room update
    // receive intial room state (and what room it's going to)
    // receive initial chat state
}

function newRoomBtn() {
    // create the element

    // add an event listener for when it is clicked
    // create a new chat element with the proper room attribute (give a place for incoming websocket message to go)
    // tell the ws endpoint to change rooms (get initial state from response and populate it in the new div)
    // delete old room div and make new room div visible (done immediately after resopnse is sent so while waiting to load, users waits at an empty screen and if they try to type, they will be in the new room so no weird races)
}

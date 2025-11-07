const newRoomBtn = document.getElementById('newRoom');
const roomList = document.getElementById('roomList');

newRoomBtn.addEventListener('click', () => {
    newRoom();
});

function newRoom() {
    // create new button
    const newBtn = document.createElement("button")
    newBtn.textContent = "New Room"
    newBtn.contentEditable = true
    roomList.appendChild(newBtn)

    // focus and select all text
    newBtn.focus()
    selectAllText(newBtn)


    // stop editing and create room
    const finishEditing = () => {
        newBtn.contentEditable = false;
        const roomName = newBtn.textContent.trim();

        if (roomName === '' || roomName === 'New Room') {
            newBtn.remove();
            return;
        }

        createRoom(roomName);

        // add click event listener to switch the room divs
        newBtn.addEventListener("click", () => {
            document.querySelector(".active").classList.remove("active")
            document.querySelector(`#${roomName}`).classList.add("active")

            // then join room on back-end
        })
    }

    // handle when user finishes editing (press Enter or blur)
    newBtn.addEventListener('blur', finishEditing, { once: true})
    newBtn.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            newBtn.blur();
        }
    })
}

function selectAllText(element) {
    const range = document.createRange();
    range.selectNodeContents(element);
    const selection = window.getSelection();
    selection.removeAllRanges();
    selection.addRange(range);
}

async function createRoom(roomName) {
    // send to backend
    console.log(roomName);

    // create room on frontend
    const newRoom = document.createElement("div")
    newRoom.id = roomName
    newRoom.classList.add("room")
    document.querySelector(".rooms").appendChild(newRoom)
}

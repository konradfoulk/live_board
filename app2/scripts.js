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
        newBtn.contentEtible = false;
        const roomName = newBtn.textContent.trim();

        if (roomName === '' || roomName === 'New Room') {
            newBtn.remove();
            return;
        }

        createRoom(roomName);
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
    
}

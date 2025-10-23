document.addEventListener('DOMContentLoaded', () => {
            const newRoomBtn = document.getElementById('newRoom');
            const roomList = document.getElementById('roomList');

            newRoomBtn.addEventListener('click', () => {
                createNewRoom();
            });

            function createNewRoom() {
                // Create new button element
                const newButton = document.createElement('button');
                newButton.textContent = 'New Room';
                newButton.contentEditable = true;
                newButton.classList.add('room-button');
                
                // Add to room list
                roomList.appendChild(newButton);
                
                // Focus and select all text
                newButton.focus();
                selectAllText(newButton);
                
                // Handle when user finishes editing (press Enter or blur)
                let roomCreated = false;
                
                const finishEditing = () => {
                    if (roomCreated) return;
                    roomCreated = true;
                    
                    newButton.contentEditable = false;
                    const roomName = newButton.textContent.trim();
                    
                    if (roomName === '' || roomName === 'New Room') {
                        // If empty or unchanged, remove the button
                        newButton.remove();
                        return;
                    }
                    
                    // Send to backend
                    createRoomOnBackend(roomName);
                };
                
                // Press Enter to finish
                newButton.addEventListener('keydown', (e) => {
                    if (e.key === 'Enter') {
                        e.preventDefault();
                        newButton.blur();
                    }
                });
                
                // Blur (click away) to finish
                newButton.addEventListener('blur', finishEditing, { once: true });
            }

            function selectAllText(element) {
                const range = document.createRange();
                range.selectNodeContents(element);
                const selection = window.getSelection();
                selection.removeAllRanges();
                selection.addRange(range);
            }

            async function createRoomOnBackend(roomName) {
                try {
                    const response = await fetch('/api/rooms', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ name: roomName })
                    });
                    
                    if (!response.ok) {
                        throw new Error('Failed to create room');
                    }
                    
                    const data = await response.json();
                    console.log('Room created:', data);
                    
                    // You can handle the response here (e.g., update button with room ID)
                } catch (error) {
                    console.error('Error creating room:', error);
                    alert('Failed to create room. Please try again.');
                }
            }
        });
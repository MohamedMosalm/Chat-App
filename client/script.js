document.addEventListener("DOMContentLoaded", function () {
  const createChatRoomButton = document.getElementById("createChatRoomButton");
  const getChatRoomButton = document.getElementById("getChatRoomButton");
  const sendButton = document.getElementById("sendButton");
  const messageInput = document.getElementById("messageInput");
  const chatMessages = document.getElementById("chat-messages");
  const chatRoomNameInput = document.getElementById("chatRoomName");

  let myClientId = localStorage.getItem("client_id") || null; 
  let ws = null;

  createChatRoomButton.addEventListener("click", function () {
    const chatRoomName = chatRoomNameInput.value.trim();
    if (chatRoomName) {
      createChatRoom(chatRoomName);
    }
  });

  getChatRoomButton.addEventListener("click", function () {
    const chatRoomName = chatRoomNameInput.value.trim();
    if (chatRoomName) {
      getChatRoom(chatRoomName);
    }
  });

  function createChatRoom(roomName) {
    fetch("/create-room", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: `room_name=${roomName}`,
    })
      .then((response) => response.json())
      .then((data) => {
        const roomId = data.room_id;
        switchRoom(roomId);
      })
      .catch((error) => console.error("Error creating room:", error));
  }

  function getChatRoom(roomName) {
    fetch(`/get-room?room_name=${roomName}`)
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error("Room not found");
        }
      })
      .then((data) => {
        const roomId = data.room_id;
        switchRoom(roomId);
      })
      .catch((error) => console.error("Error getting room:", error));
  }

  function switchRoom(roomId) {
    chatMessages.innerHTML = "";

    if (ws) {
      ws.close();
    }

    connectToWebSocket(roomId);

    const baseUrl = window.location.origin;
    const newUrl = `${baseUrl}/chat/${roomId}`;
    window.history.pushState({ path: newUrl }, "", newUrl);
  }

  function connectToWebSocket(roomId) {
    ws = new WebSocket(
      `ws://localhost:4000/ws?room_id=${roomId}&client_id=${myClientId || ""}`
    );

    ws.onopen = function () {
      console.log("Connected to WebSocket server");
    };

    ws.onmessage = function (event) {
      const messageData = JSON.parse(event.data);

      if (!myClientId && messageData.sender_id) {
        myClientId = messageData.sender_id;
        localStorage.setItem("client_id", myClientId);
        console.log("Received client ID:", myClientId);
      } else {
        if (messageData.content === "Connected to the server") {
          return; 
        }
        const isOwnMessage = messageData.sender_id === myClientId;
        appendMessage(messageData, isOwnMessage);
      }
    };

    ws.onclose = function () {
      console.log("Disconnected from WebSocket server");
    };

    ws.onerror = function (error) {
      console.error("WebSocket error observed:", error);
    };
  }

  sendButton.addEventListener("click", function () {
    const message = messageInput.value.trim();
    if (message && myClientId && ws) {
      const messageData = {
        sender_id: myClientId,
        content: message,
        timestamp: new Date().toLocaleTimeString(),
      };
      ws.send(JSON.stringify(messageData));
      appendMessage(messageData, true);
      messageInput.value = "";
    }
  });

  messageInput.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
      event.preventDefault();
      sendButton.click();
    }
  });

  function appendMessage(message, isOwnMessage) {
    const messageElement = document.createElement("div");
    if (isOwnMessage) {
      messageElement.classList.add(
        "d-flex",
        "flex-row",
        "justify-content-end",
        "mb-4"
      );
      messageElement.innerHTML = `
          <div>
            <p class="small p-2 me-3 mb-1 text-white rounded-3 bg-primary">${message.content}</p>
            <p class="small me-3 mb-3 rounded-3 text-muted d-flex justify-content-end">${message.timestamp}</p>
          </div>
          <img src="https://mdbcdn.b-cdn.net/img/Photos/new-templates/bootstrap-chat/ava3-bg.webp"
              alt="avatar 1" style="width: 45px; height: 100%;">
        `;
    } else {
      messageElement.classList.add(
        "d-flex",
        "flex-row",
        "justify-content-start",
        "mb-4"
      );

      messageElement.innerHTML = `
          <img src="https://mdbcdn.b-cdn.net/img/Photos/new-templates/bootstrap-chat/ava4-bg.webp"
              alt="avatar 1" style="width: 45px; height: 100%;">
          <div>
            <p class="small p-2 ms-3 mb-1 rounded-3 bg-body-tertiary">${message.content}</p>
            <p class="small ms-3 mb-3 rounded-3 text-muted">${message.timestamp}</p>
          </div>
        `;
    }
    chatMessages.appendChild(messageElement);
    chatMessages.scrollTop = chatMessages.scrollHeight;
  }
});

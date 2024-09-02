document.addEventListener("DOMContentLoaded", function () {
  const sendButton = document.getElementById("sendButton");
  const messageInput = document.getElementById("messageInput");
  const chatMessages = document.getElementById("chat-messages");

  let myClientId = null;

  const ws = new WebSocket(`ws://localhost:4000/ws`);

  ws.onopen = function () {
      console.log("Connected to WebSocket server");
  };

  ws.onmessage = function (event) {
      const messageData = JSON.parse(event.data);

      if (!myClientId) {
          myClientId = messageData.content;
          console.log("Received client ID:", myClientId);
      } else {
          appendMessage(messageData, messageData.sender_id === myClientId);
      }
  };

  ws.onclose = function () {
      console.log("Disconnected from WebSocket server");
  };

  ws.onerror = function (error) {
      console.error("WebSocket error observed:", error);
  };

  sendButton.addEventListener("click", function () {
      const message = messageInput.value.trim();
      if (message && myClientId) {
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

document.addEventListener("DOMContentLoaded", function () {
  const sendButton = document.getElementById("sendButton");
  const messageInput = document.getElementById("messageInput");
  const chatMessages = document.getElementById("chat-messages");

  sendButton.addEventListener("click", function () {
    const message = messageInput.value.trim();

    if (message) {
      const messageElement = document.createElement("div");
      messageElement.classList.add(
        "d-flex",
        "flex-row",
        "justify-content-end",
        "mb-4"
      );

      messageElement.innerHTML = `
              <div>
                  <p class="small p-2 me-3 mb-1 text-white rounded-3 bg-primary">${message}</p>
                  <p class="small me-3 mb-3 rounded-3 text-muted d-flex justify-content-end">${new Date().toLocaleTimeString()}</p>
              </div>
              <img src="https://mdbcdn.b-cdn.net/img/Photos/new-templates/bootstrap-chat/ava3-bg.webp"
                  alt="avatar 1" style="width: 45px; height: 100%;">
          `;

      // <img src="https://mdbcdn.b-cdn.net/img/Photos/new-templates/bootstrap-chat/ava4-bg.webp"
      //         alt="avatar 1" style="width: 45px; height: 100%;">

      chatMessages.appendChild(messageElement);
      messageInput.value = "";
      chatMessages.scrollTop = chatMessages.scrollHeight;
    }
  });

  messageInput.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
      event.preventDefault();
      sendButton.click();
    }
  });
});

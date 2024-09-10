# Chat App

A real-time chat application built with **Go** and **Fiber** using **WebSockets** for client-server communication and **PostgreSQL** for message and room persistence. This application supports multiple chat rooms, message broadcasting, and message history loading upon room entry.

## Features

- Real-time messaging between clients in the same room using WebSockets.
- Multiple chat rooms with unique names.
- Simple and intuitive user interface with message input and display.
- UUID-based identification for users and rooms.

## Tech Stack

- **Backend**: Go, Fiber, Gorm, WebSockets
- **Database**: PostgreSQL
- **Frontend**: HTML, CSS, JavaScript

## Setup and Installation

### Prerequisites

- [Go](https://golang.org/doc/install) (v1.19 or above)
- [PostgreSQL](https://www.postgresql.org/download/) (with UUID extension)
- [Git](https://git-scm.com/)

### Clone the Repository


```bash
git clone https://github.com/MohamedMosalm/Chat-App.git
cd chat-app/server
```

### Install Dependencies

```bash
go mod tidy
```

## Configure Database

Create a PostgreSQL database and configure the connection details in .env file to match your PostgreSQL credentials:

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chat_app


Run the following SQL command to ensure the UUID extension is available:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### Run the Application

Start the application:

```bash
go run *.go
```
The app will run on http://localhost:4000.

Open a browser and navigate to http://localhost:4000/ to access the chat interface. You can create or join rooms and start chatting in real-time.

### Endpoints
- POST /create-room: Create a new chat room.
- GET /get-room?room_name={name}: Get an existing room by name.
- GET /ws?room_id={id}&client_id={id}: Connect to a WebSocket endpoint for real-time chat.

### Usage
- Enter a room name and click Create Room or Get Room to join or create a room.
- Type a message and press Send to chat with others in the room.

### Future Enhancements
- User authentication and profiles.
- Typing indicators and message status (delivered, seen).
- Support for private messaging and direct user-to-user chats.
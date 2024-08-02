# Simple TCP Server in Go

This project is a simple TCP server written in Go. The server listens for incoming TCP connections, reads messages from connected clients, and sends a response back to the clients. Received messages are logged to the console.

## Features

- Accepts multiple client connections concurrently.
- Reads and logs messages from clients.
- Sends a response back to the clients after receiving a message.
- Gracefully handles connection errors.

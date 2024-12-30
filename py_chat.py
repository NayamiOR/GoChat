# websocket_im.py
import asyncio
import websockets

clients = set()

async def server_handler(websocket, path):
    clients.add(websocket)
    try:
        async for message in websocket:
            # Broadcast the message to all connected clients
            for client in clients:
                if client != websocket:
                    await client.send(message)
    except websockets.ConnectionClosed:
        pass
    finally:
        clients.remove(websocket)

async def client():
    uri = "ws://localhost:8080"
    async with websockets.connect(uri) as websocket:
        asyncio.create_task(receive_messages(websocket))
        while True:
            message = input("Enter message (type 'exit' to quit): ")
            if message.lower() == "exit":
                break
            await websocket.send(message)

async def receive_messages(websocket):
    try:
        async for message in websocket:
            print("Received:", message)
    except websockets.ConnectionClosed:
        print("Connection closed.")

if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("Usage: python websocket_im.py [server|client]")
        sys.exit(1)

    if sys.argv[1] == "server":
        start_server = websockets.serve(server_handler, "localhost", 8080)
        asyncio.get_event_loop().run_until_complete(start_server)
        print("Server started on ws://localhost:8080")
        asyncio.get_event_loop().run_forever()
    elif sys.argv[1] == "client":
        asyncio.get_event_loop().run_until_complete(client())
    else:
        print("Invalid argument, use 'server' or 'client'")

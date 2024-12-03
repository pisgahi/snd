## Snd: Secure File Transfer

Snd is a file transfer program built using Golang, leveraging TCP for reliable and efficient data transmission. This application allows users to send and receive files securely over a network.

## Features

- Fast and reliable file transfer over TCP.
- Secure file transmission
- Simple command-line interface for ease of use.
- Lightweight and efficient, suitable for various use cases.

## Installation

To get started with the Snd application, follow these steps:

1. <b>Clone this repo</b>

```bash
git clone https://github.com/yourusername/snd.git
```

2. <b>Cd into the repo</b>

```bash
cd snd
```

3. <b>Manage Dependecies</b>

```bash
go mod tidy
```

4. <b>Configure Snd</b>

```bash
./snd -cert="cert/certs" -addr="127.0.0.1:4040" -dir="User/username/Desktop"
```

Before you can start using Snd you need to setup some config variables.

- <b>-cert:</b> the directory you would like to save your TLS certificates and keys
- <b>-addr</b>: the address your server starts on
- <b>-dir</b> the directory to store received files

5. <b>Build Snd</b>

```bash
go build -o snd
```

```bash
./snd -cert="cert/certs" -addr="127.0.0.1:4040" -dir="User/username/Desktop"
```

<b>Running the Server</b>

```bash
./snd -s
```

The server will listen for incoming file transfer requests, on the configured port.

## Sending a File

```bash
./snd -f="your_filepath" -to="127.0.0.1:4040"
```

## Receiving a File

The server will automatically receive the file sent by the client and store it in the specified directory.

## How It Works

- Server: Listens for incoming TCP connections and handles file reception.
- Client: Connects to the server and sends files over the established TCP connection.

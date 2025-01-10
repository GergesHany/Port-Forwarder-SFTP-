# Secure FTP-Like File Transfer Project in Go

This project demonstrates a **secure FTP-like file transfer system** using Go. It consists of three components:

1. **Sender Server**: Hosts a file and sends it to the forwarder using **SFTP**.
2. **Receiver Server**: Receives the file and saves it in a specified directory
3. **Forwarder**: Retrieves the file from the sender and forwards it to the receiver using **SFTP**.

The project uses **SSH keys** for secure authentication and encryption, ensuring secure file transfers over an unsecured network.

---

## Project Structure

```plaintext
Port-Forwarder/
├── sender/
│   └── sender.go          # SFTP sender code
├── receiver/
│   └── receiver.go        # SFTP receiver code
├── forward/
│   └── forward.go         # SFTP forwarder code
├── keys/
│   ├── id_rsa             # Private key for authentication
│   └── id_rsa.pub         # Public key for authentication
├── run_ftp.sh             # Bash script to automate running the project
└── README.md              # Project documentation
```

<hr>

## How to Use

#### 1. Clone the Repository

- Clone the project repository to your local machine:

  ```bash
     git clone git@github.com:GergesHany/Port-Forwarder-SFTP-.git
     cd Port-Forwarder-SFTP-

  ```

<hr>

#### 2. Generate SSH Keys (If Not Already Generated)

- The project requires SSH keys for secure authentication. If the keys/ directory doesn't exist, the run_ftp.sh script will generate them automatically.

- To manually generate SSH keys:

  ```bash
  mkdir -p keys
  ssh-keygen -t rsa -b 4096 -f keys/id_rsa -N "" # No passphrase
  ```

- This will create:
  - keys/id_rsa: Private key (used by the sender and forwarder).
  - keys/id_rsa.pub: Public key (used by the receiver).

<hr>

#### 3. Build and Run the Project

Option 1: Manual Execution

1. Start the Receiver Server:
   Navigate to the receiver directory and run the `receiver` with the upload directory as an argument:

   ```bash
    cd receiver
    go run receiver.go <upload_directory>
   ```

2. Start the Sender Server:
   Open a new terminal, navigate to the sender directory, and run the `sender` with the file path as an argument:

   ```bash
    cd sender
    go run sender.go <file_path>
   ```

3. Start the Forwarder:
   Open another terminal, navigate to the `forward` directory, and run the forwarder:
   ```bash
    cd forward
    go run forward.go <file_path> <directory>
   ```

Option 2: Automated Execution (Using run_ftp.sh)

1. Make the Script Executable:
   Run the following command to make the script executable:

   ```bash
   chmod +x run_ftp.sh
   ```

2. Run the Script:
   Execute the script with the file name and upload directory as arguments:

   ```bash
     ./run_ftp.sh <file_path> <upload_directory>
   ```

This will:

- Open a new terminal and start the receiver.
- Open another terminal and start the sender with the specified file.
- Open a third terminal and start the forwarder with the same file and upload directory.

<hr>

## Example

1. Create a Sample File:
   Create a file named example.txt in the sender directory:

   ```bash
   echo "Hello, World!" > sender/example.txt
   ```

2. Run the Project:
   Execute the run_ftp.sh script with the file name as an argument:

   ```bash
   ./run_ftp.sh example.txt upload
   ```

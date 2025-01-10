
# Check if the file path and upload directory are provided as arguments
if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Usage: ./run_ftp.sh <file_path> <upload_directory>"
  exit 1
fi

FILE_PATH=$1
UPLOAD_DIR=$2

# Get the current directory path
CURRENT_DIR=$(pwd)

# Open a new terminal and run files (sender, receiver, forwarder)
# and wait 2 seconds between each terminal to allow the previous one to start

gnome-terminal -- bash -c "cd $CURRENT_DIR/receiver; go run receiver.go $UPLOAD_DIR; exec bash"
sleep 2
gnome-terminal -- bash -c "cd $CURRENT_DIR/sender; go run sender.go $FILE_PATH; exec bash"
sleep 2
gnome-terminal -- bash -c "cd $CURRENT_DIR/forward; go run forward.go $FILE_PATH $UPLOAD_DIR; exec bash"

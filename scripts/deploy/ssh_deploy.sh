SERVICE_BIN=$1
SSH_HOST_LIST=$2

if [[ -z $SSH_HOST_LIST ]]; then
    echo "SSH_HOST_LIST variable is not set"
    exit 1
fi

if [[ -z $SERVICE_BIN ]]; then
    SERVICE_BIN="./bin/asvsoft"
fi

for SSH_HOST in $SSH_HOST_LIST; do
    scp -o ConnectTimeout=3 -q $SERVICE_BIN $SSH_HOST:/usr/local/bin/
    if [[ $? -ne 0 ]]; then 
        echo "failed deploy binary file to $SSH_HOST"
        exit 1
    fi
done

#!/bin/bash

starttime=$(date +%s)
export TIMEOUT=10
export DELAY=3

# Set environment 
export PATH=$PATH:${PWD}/bin
export FAB_BINS=${PWD}/bin
export FABRIC_CFG_PATH=${PWD}/config
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
export FTS_PATH=$GOPATH/src/github.com/hyperledger-labs/fabric-token-sdk
export OAPI_PATH=$GOPATH/src/github.com/deepmap/oapi-codegen
export $(./test-network/setOrgEnv.sh bank | xargs)
export $(./test-network/setOrgEnv.sh business | xargs)
export TEST_NETWORK_HOME=${PWD}/test-network;

IMAGE_NAME="hyperledger/fabric"
IMAGE_TAG="2.5.0"
if docker image ls -q ${IMAGE_NAME}:${IMAGE_TAG} > /dev/null 2>&1; then
    echo "[+] The Hyperledger Fabric 2.5.0 image exists locally. Start the network up.";
else
    bash install-fabric.sh --fabric-version 2.5.0 docker
    echo "[+] Install the Hyperledger Fabric 2.5.0 image locally. Start the network up.";
fi

# Start the network up
cd ./test-network && bash network.sh up;
# Create private channel
# bash network.sh createChannel -c fast-channel && bash network.sh createChannel -c stable-channel;
bash network.sh createChannel;

# Add Org3 join to channel
# bash ./addOrg3/addOrg3.sh up -c fast-channel;

# Install token sdk
sudo git clone https://github.com/hyperledger-labs/fabric-token-sdk.git $FTS_PATH;
sudo git clone https://github.com/deepmap/oapi-codegen.git $OAPI_PATH;
cd ../token-sdk;
go install github.com/hyperledger-labs/fabric-token-sdk/cmd/tokengen@v0.3.0;
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest;
bash ./scripts/up.sh;

# Install high-throughput chaincode
cd ../test-network;
bash network.sh deployCC -ccn bigdatacc -ccp ../high-throughput/chaincode-go/ -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')" -cci Init;

cat <<EOF

Total setup execution time : $(($(date +%s) - starttime)) seconds ...

EOF

# Start monitor
cd ./prometheus-grafana && docker-compose up -d;
cd .. && bash monitordocker.sh;

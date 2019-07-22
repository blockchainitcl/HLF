FROM hyperledger/fabric-peer:amd64-1.3.0
RUN mkdir fabric/
COPY channel1.tx /fabric/channel1.tx
COPY genesis.block /fabric/genesis.block
RUN mkdir crypto-config
COPY peerOrganizations fabric/crypto-config/peerorganizations
RUN mkdir -p chaincode/src/chaincode/
COPY chaincode/src/chaincode/Wallet/Wallet.go /fabric/chaincode/src/chaincode/Wallet/Wallet.go
COPY Org1MSPanchors.tx /fabric/Org1MSPanchors.tx
COPY configtx.yaml /fabric/configtx.yaml

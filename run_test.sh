#!/bin/bash

# Start Hardhat node in background
cd blockchain-hardhat/
npm install --save-dev hardhat
npm install --save-dev @nomicfoundation/hardhat-toolbox
npm install --save-dev @nomicfoundation/hardhat-ignition-ethers
npx hardhat node &

# Wait a bit for Hardhat to start
sleep 2

# Start Go app (foreground)
cd ../
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
go run zkrevoke --test

cd plots/
poetry install --no-root
poetry run python main.py --test
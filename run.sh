#!/bin/bash

# Start Hardhat node in background
cd blockchain-hardhat/
npx hardhat node &

# Wait a bit for Hardhat to start
sleep 2

# Start Go app (foreground)
cd ../
go run zkrevoke

cd plots/
poetry run python main.py
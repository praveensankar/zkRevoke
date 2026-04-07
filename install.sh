#!/bin/bash

cd blockchain-hardhat
npm install --save-dev hardhat
npm install --save-dev @nomicfoundation/hardhat-toolbox
npm install --save-dev @nomicfoundation/hardhat-ignition-ethers

cd ..
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

cd plots
poetry install --no-root
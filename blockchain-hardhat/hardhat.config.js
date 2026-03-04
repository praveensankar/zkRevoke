/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  networks: {
    hardhat: {
      accounts: [
        {
          privateKey: "df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e",
          balance: "1000000000000000000000000000000000000" // 1000 ETH in wei
        }
      ]
    }
  },
  solidity: "0.8.28",
};

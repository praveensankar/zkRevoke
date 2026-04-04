# Artifact Appendix

Paper title: **zkRevoke: Configurable Untraceability for Verifiable Credentials using ZKPs**

Requested Badge(s):
  - [✓] **Available**
  - [✓] **Functional**
  - [✓] **Reproduced**


## Description

Praveensankar Manimaran, Mayank Raikwar, Thiago Garrett, Arlindo F. da Conceição, Leander Jehl and Roman Vitenberg. "zkRevoke: Configurable Untraceability for Verifiable Credentials using ZKPs" Proceedings on Privacy Enhancing Technologies (2026).

This repository contains the implementation of zkRevoke, a protocol that limits the amount of time a verifier can track a VC’s revocation status from the information provided in a VP.


### Security/Privacy Issues and Ethical Concerns 

This implementation does have any security implications for testing. In addition, the implementation uses synthetic dataset for testing. Therefore, it does not have implications for privacy or ethical issues. 



## Basic Requirements 


### Hardware Requirements 

Can run on a laptop (No special hardware requirements).

The hardware choice only impacts computation times in the experiments, without affecting storage and bandwidth consumption. Regarding computation times, our main focus is on the holder. Since holders are likely to use commodity machines, we have performed the experiments on a desktop computer running a 12-Core Apple M4 Pro CPU with 24 GB RAM.  Using a more powerful server would affect the absolute values, but would be unlikely to significantly influence comparative evaluation relative to the baselines.


### Software Requirements 

We have performed the experiments on Mac OS 26.3. However, it is possible to run the experiments on any operating system.

Programming Languages: Golang, Python, Solidity.

#### The compilers and interpreters used for the evaluation in the paper
* go v1.23.0
* python >=3.13
* solidity v0.8.28
* nodejs >= 23.10.0 <= 24.14.1

#### Other dependencies
* hardhat v2.22.19
* Poetry v2.1.4

### Estimated Time and Storage Consumption

* Total time to run the artifact: 1 hour
* Required disk space: 1 GB


## Environment 

In the following, we describe how to access the artifact and all related and
necessary data and software components. Afterward, we describe how to set up
everything and how to verify that everything is set up correctly.

### Accessibility 

Link to the artifact: https://github.com/praveensankar/zkRevoke

### Set up the environment

Install the following languages: GoLang and Node.js.

Then install:

(a) Hardhat - Hardhat is used to run a local blockchain network.
```bash
npm install --save-dev hardhat
npm install --save-dev @nomicfoundation/hardhat-toolbox
npm install --save-dev @nomicfoundation/hardhat-ignition-ethers
```
Refer to the official hardhat tutorial if you encounter any errors in the installation of hardhat:
https://hardhat.org/tutorial/creating-a-new-hardhat-project

(b) Solc, abigen - Solc and abigen are used for deploying smart contracts to the hardhat network and interacting with the deployed smart contract.

Follow the instructions given in https://goethereumbook.org/en/smart-contract-compile/. The abigen tool can be installed using the following command:

```bash
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

### Testing the Environment (Required for Functional and Reproduced badges)

verify the installation of GoLand and Node.js.
```bash
go version
node -v
```



## Artifact Evaluation (Required for Functional and Reproduced badges)

This section includes all the steps required to evaluate the artifact's
functionality and validate the paper's key results and claims. Therefore, it
highlights the paper's main results and claims. 


### Test Experiments

These steps evaluate the artifact using small parameters to test the working of the complete flow.
This process takes a few minutes of computation time. 

####  (a) Run a local blockchain network: 
Open a terminal and go to the root directory of the project and execute the following:
```bash
cd blockchain-hardhat
npx hardhat node
```

#### (b) Run the golang program
Open a terminal and go to the root directory of the project and execute the following:
```bash
go run zkrevoke --test
```
This command will generate results and store them in the following directories:
```bash
benchmark/results/
irma/benchmark/results/
```

#### (c) Plot the results
Open a terminal and go to the root directory of the project and execute the following:
```bash
cd plots
poetry install --no-root 
poetry run python main.py --test
```
These commands will generate plots and tables used in the paper and store them in the following directory:
```bash
plots/graphs/
```
The filenames of the generated figures and tables include the caption numbers used in the paper.



### Main Experiments

####  (a) Run a local blockchain network: 
Open a terminal and go to the root directory of the project and execute the following:
```bash
cd blockchain-hardhat
npx hardhat node
```

#### (b) Run the golang program
Open a terminal and go to the root directory of the project and execute the following:
```bash
go run zkrevoke
```
This command will generate results and store them in the following directories:
```bash
benchmark/results/
irma/benchmark/results/
```

#### (c) Plot the results
Open a terminal and go to the root directory of the project and execute the following:
```bash
cd plots
poetry install --no-root 
poetry run python main.py
```
These commands will generate plots and tables used in the paper and store them in the following directory:
```bash
plots/graphs/
```
The filenames of the generated figures and tables include the caption numbers used in the paper.




### workload used for the results in the paper


+ total number of issued VCs (n): 1 million
+ expiration period of VCs: 365
+ epoch duration: 1 
+ revocation rate (R): 1% to 15%
+ revoke a fixed number ($r = (10^6*\mathcal{R})/(100*365)$)  of randomly choosen VCs in each epoch
+ verification period of a VP (m): 1 day to 60
+ we consider the values of k from $1$ to $2^5$


### Main Results and Claims



#### Main Result 1: Holder’s Bandwidth

Holders in zkRevoke require significantly less bandwidth compared to IRMA since constructing VPs in zkRevoke does not depend on the revocations of other VCs, unlike in IRMA.  In the experiment, we change the verification period (m) linearly, from 1 to 60, and we measure the holder's bandwidth. The results are stored in the file: plots/graphs/fig_1c_result_one_time_sharing_holder_bandwidth.png. We report these results in "Figure 1c" of our paper.

#### Main Result 2: Holder’s Computation

The computation for holders in zkRevoke does not depend on the revocation of other VCs, in contrast to IRMA. When the extension of proving multiple tokens is used, zkRevoke incurs computation overhead for holders that is comparable to or better than the one in IRMA. In the experiment, we change the verification period (m) linearly, from 1 to 60, and we measure the holder's computation requirements. The results are stored in the files: (a) plots/graphs/fig_1a_result_one_time_sharing_computation.png, and (b) plots/graphs/fig_1b_result_one_time_sharing_computation_irma_k_2.png. We report these results in "Figure 1a" and "Figure 1b" of our paper.

#### Main Result 3: Issuer's Bandwidth
zkRevoke avoids significant bandwidth consumption required in IRMA to broadcast factors, but requires more bandwidth than IRMA-registry since the blacklist is replaced every epoch, not just appended to. In the experiment, we change the current epoch (e) linearly, from 1 to 365, and we measure the issuer's bandwidth requirements. The results are stored in the files: (a) plots/graphs/fig_2a_result_revocation_issuer_bandwidth_without_repition_r1.png, and (b) plots/graphs/fig_2b_result_revocation_issuer_bandwidth_without_repition_r5.png. We report these results in "Figure 2a" and "Figure 2b" of our paper.

#### Main Result 4: Issuer's Computation
Issuers need to perform more computation in zkRevoke compared to IRMA, since the issuer
has to recompute tokens for all revoked VCs after each epoch. In the experiment, we change the current epoch (e) linearly, from 1 to 365, and we measure the issuer's computation requirements. The results are stored in the file: plots/graphs/fig_2c_result_revocation_computation_including_commitment_1.png. We report these results in "Figure 2c" of our paper.

#### Main Result 5: Verifier's computation for the verification of commitment 
The time to verify the commitment is on the order of microseconds even for higher revocation rates, which is insignificant compared to the computation requirements for verification of proofs in a VP. In the experiment, we change the current epoch (e) linearly, from 1 to 365, and we measure the verifier's computation requirements. The results are stored in the file: plots/graphs/fig_3_result_list_commitment_verification_time.png. We report these results in "Figure 3" of our paper.

#### Main Result 5: ZK Proof size
The size of a ZK proof in zkRevoke is 164 bytes, whereas the size of a non-revocation proof in IRMA is 1831 bytes. zkRevoke requires 11x less bandwidth for proofs, compared to IRMA.
In the experiment, we change the verification period (m) linearly, from 1 to 60, and we measure the proof size. The results are stored in the file: plots/graphs/fig_4a_result_zkp_proof_size.png. We report these results in "Figure 4a" of our paper.

#### Main Result 6: ZK Proof Generation Time.
The proof generation time in IRMA is 10x faster than the proof generation time in zkRevoke.
In the experiment, we change the verification period (m) linearly, from 1 to 60, and we measure the time to generation ZK proofs. The results are stored in the file: plots/graphs/fig_4a_result_zkp_proof_gen_time.png. We report these results in "Figure 4b" of our paper.

#### Main Result 7: ZK Proof Verification Time.
The proof verification time in zkRevoke is 6x faster than the proof generation time in
IRMA.
In the experiment, we change the verification period (m) linearly, from 1 to 60, and we measure the time to verify all ZK proofs for a VP. The results are stored in the file: plots/graphs/fig_4a_result_zkp_proof_ver_time.png. We report these results in "Figure 4c" of our paper.

#### Main Result 8: ZK Circuit
To generate and verify proofs, ZK circuits are one of the parameters in zkRevoke, whereas public keys are sufficient in IRMA. In IRMA, the ZKP scheme is designed as a
digital signature scheme and enables holders and verifiers to use
public keys to generate and verify proofs. IRMA outperforms zkRevoke since public keys are shorter in size and take a significantly smaller amount of time to generate
compared to ZK circuits. 
The size of a ZK proof in zkRevoke is 164 bytes, whereas the size of a non-revocation proof in IRMA is 1831 bytes. zkRevoke requires 11x less bandwidth for proofs, compared to IRMA.
In the experiment, we change the number of tokens per circuit (k) linearly, from 1 to 60 and compare the size and generation time.  The results are stored in the file: (a) plots/graphs/fig_5a_result_circuit_size.png and (b) plots/graphs/fig_5a_result_circuit_time.png. We report these results in "Figure 5a" and "Figure 5b" of our paper.

#### Main Result 9: ZK Circuit Breakdown
In the ZK circuit, the hash-related conditions require 661 constraints, whereas signature verification requires 7993 constraints. This indicates that signature verification is a dominant part (85.8%) of the circuit. The results are stored in the file: (a) plots/graphs/table_4:Circuit_Constraints_Breakdown.png. We report these results in "Table 4" of our paper.

#### Main Result 10: Groth16 Overhead
we evaluate the overhead of the Groth16 scheme by considering the performance of the “empty” circuit instantiated without any conditions. The witness and proof generation for the empty circuit add 5% and 1.5% overhead respectively compared to the complete circuit, whereas proof verification for the empty circuit requires almost the same amount of time as for the complete circuit. The results are stored in the file: (a) plots/graphs/table_5:Groth16_Overhead.png. We report these results in "Table 5" of our paper.


## Limitations 

All the results should be reproducible. 



## Notes on Reusability

The artifact can be deployed in any operating system. Since we have evaluated in a personal computer, the artifact can be adapted to different hardware requirements as well. 


############################################################################################################################

Please find documentation at: https://zk-revoke.mintlify.app/introduction


############################################################################################################################

License details: https://github.com/praveensankar/zkRevoke/blob/main/LICENSE

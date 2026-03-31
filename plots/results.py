
import json
from collections import defaultdict
from re import error

import math
import os

class ZKRevokePresentationResult:
    def __init__(self, vpValidity, numberOfTokensInCircuit, numberOfZKPProofs, zkpProofSize, totalZKPProofSize, zkpProofGenTime, totalZKPProofGenTime):
        self.VPValidity = vpValidity
        self.numberOfTokensInCircuit = numberOfTokensInCircuit
        self.numberOfZKPProofs = numberOfZKPProofs
        self.zkpProofSize = zkpProofSize
        self.totalZKPProofSize = totalZKPProofSize
        self.zkpProofGenTime = zkpProofGenTime
        self.totalZKPProofGenTime = totalZKPProofGenTime



    def __eq__(self, another):
        return self.VPValidity == another.VPValidity and self.numberOfTokensInCircuit == another.numberOfTokensInCircuit

    def __hash__(self):
        return hash(str(self.VPValidity) + str(self.numberOfTokensInCircuit))


def parse_zkrevoke_presentation_result_entry(file):
    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, "benchmark", "results", file)
    with open(file_path) as f:
        json_data = json.load(f)

    entries = []
    for entry in json_data:
        vpValidity = entry['vp_validity_period']
        numberOfTokensInCircuit = entry['number_of_tokens_in_circuit']
        numberOfZKPProofs = entry['number_of_zkp_proofs']
        zkpProofSize = entry['zkp_proof_size']
        totalZKPProofSize = entry['total_zkp_proof_size']
        zkpProofGenTime = entry['time_to_generate_one_zkp_proof']
        totalZKPProofGenTime = entry['time_to_generate_all_zkp_proofs']

        entry = ZKRevokePresentationResult(vpValidity=vpValidity,
                                           numberOfTokensInCircuit=numberOfTokensInCircuit,
                                           numberOfZKPProofs=numberOfZKPProofs,
                                           zkpProofSize=zkpProofSize,
                                           totalZKPProofSize=totalZKPProofSize,
                                           zkpProofGenTime=zkpProofGenTime,
                                           totalZKPProofGenTime=totalZKPProofGenTime)

        entries.append(entry)

    return entries



class ZKRevokeVerificationResult:
    def __init__(self, vpValidity, numberOfTokensInCircuit, zkpProofVerTime):
        self.vpValidity = vpValidity
        self.numberOfTokensInCircuit = numberOfTokensInCircuit
        self.zkpProofVerTime = zkpProofVerTime



    def __eq__(self, another):
        return self.vpValidity == another.vpValidity and self.numberOfTokensInCircuit == another.numberOfTokensInCircuit

    def __hash__(self):
        return hash(str(self.vpValidity) + str(self.numberOfTokensInCircuit))


def parse_zkrevoke_verification_result_entry(file):
    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, "benchmark", "results", file)
    with open(file_path) as f:
        json_data = json.load(f)

    entries = []
    for entry in json_data:
        vpValidity = entry['vp_validity_period']
        numberOfTokensInCircuit = entry['number_of_tokens_in_circuit']
        zkpProofVerTime = entry['zkp_proof_ver_time']


        entry = ZKRevokeVerificationResult(vpValidity=vpValidity,
                                           numberOfTokensInCircuit=numberOfTokensInCircuit,
                                           zkpProofVerTime=zkpProofVerTime)

        entries.append(entry)

    return entries



class IRMAPresentationVerificationResult:
    def __init__(self,  totalVCs, revocationRate, numberOfEpochs, totalRevokedVCsPerEpoch, witnessUpdateTime, proofGenTime,
                 proofSize, nonRevProofSize, proofVerTime, totalEpochs, currentEpoch, VPValidityPeriod, holderBandwidth):
        self.totalVCs = totalVCs
        self.revocationRate = revocationRate
        self.numberOfEpochs = numberOfEpochs
        self.totalRevokedVCsPerEpoch = totalRevokedVCsPerEpoch
        self.witnessUpdateTime = witnessUpdateTime
        self.proofGenTime = proofGenTime
        self.proofSize = proofSize
        self.nonRevProofSize = nonRevProofSize
        self.proofVerTime = proofVerTime
        self.totalEpochs = totalEpochs
        self.currentEpoch = currentEpoch
        self.VPValidityPeriod = VPValidityPeriod
        self.holderBandwidth = holderBandwidth


    def __eq__(self, another):
        return self.totalEpochs == another.totalEpochs and self.currentEpoch == another.currentEpoch and self.totalVCs == another.totalVCs and self.revocationRate == another.revocationRate


    def __hash__(self):
        return hash(str(self.totalEpochs) + str(self.currentEpoch) + str(self.totalVCs) + str(self.revocationRate))


def parse_irma_presentation_and_verification_result_entry(file):
    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, "irma", "benchmark", "results", file)

    with open(file_path) as f:
        json_data = json.load(f)

    entries = []
    for entry in json_data:
        totalVCs =entry['total_number_of_issued_vcs']
        revocationRate = entry['revocation_rate']
        numberOfEpochs = entry['total_number_of_epochs']
        totalRevokedVCsPerEpoch = entry['total_revoked_vcs_per_epoch']
        witnessUpdateTime = entry['time_to_update_witness']
        proofGenTime = entry['time_to_create_disclosure_proof_with_non_revocation']
        proofSize = entry['disclosure_proof_size']
        nonRevProofSize = entry['non_revocation_proof_size']
        proofVerTime = entry['proof_verification_time']
        totalEpochs = entry['total_number_of_epochs']
        currentEpoch = entry['current_epoch']
        holderBandwidth = entry['holder_bandwidth']
        VPValidityPeriod = entry['vp_validity_period']


        entry = IRMAPresentationVerificationResult(totalVCs=totalVCs,
                                                   totalEpochs=totalEpochs,
                                                   totalRevokedVCsPerEpoch=totalRevokedVCsPerEpoch,
                                                   revocationRate=revocationRate,
                                                   numberOfEpochs = numberOfEpochs,
                                                   witnessUpdateTime = witnessUpdateTime,
                                                   proofGenTime = proofGenTime,
                                                   proofSize = proofSize,
                                                   nonRevProofSize= nonRevProofSize,
                                                   proofVerTime = proofVerTime,
                                                   currentEpoch=currentEpoch,
                                                   VPValidityPeriod = VPValidityPeriod,
                                                   holderBandwidth=holderBandwidth)

        entries.append(entry)

    return entries




class IRMASetupResult:
    def __init__(self,  privateKeySize, publicKeySize, keyGenTime, accGenTime, accSize):
        self.privateKeySize = privateKeySize
        self.publicKeySize = publicKeySize
        self.keyGenTime = keyGenTime
        self.accGenTime = accGenTime
        self.accSize = accSize


def parse_irma_setup_result_entry(file):
    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, "irma", "benchmark", "results", file)

    with open(file_path) as f:
        json_data = json.load(f)

    entries = []
    for entry in json_data:
        privateKeySize =entry['private_key_size']
        publicKeySize = entry['public_key_size']
        keyGenTime = entry['key_gen_time']
        accGenTime = entry['accumulator_gen_time']
        accSize = entry['accumulator_size']



        entry = IRMASetupResult(privateKeySize=privateKeySize,
                                                   publicKeySize=publicKeySize,
                                                   keyGenTime=keyGenTime,
                                                   accGenTime=accGenTime,
                                                   accSize = accSize)

        entries.append(entry)

    return entries


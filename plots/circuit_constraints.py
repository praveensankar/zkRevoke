import json
import os

import numpy as np
from matplotlib import pyplot as plt


class CircuitConstraintResult:
    def __init__(self, circuitType, numberOfConstraints, witnessGenTime, proofGenTime, proofVerTime):
        self.circuitType = circuitType
        self.numberOfConstraints = numberOfConstraints
        self.witnessGenTime = witnessGenTime
        self.proofGenTime = proofGenTime
        self.proofVerTime = proofVerTime


    def __eq__(self, another):
        return self.circuitType == another.circuitType

    def __hash__(self):
        return hash(str(self.circuitType))


def parse_zkrevoke_circuit_constraint_result_entry(file):
    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, "benchmark", "results", file)

    with open(file_path) as f:
        json_data = json.load(f)

    entries = []
    for entry in json_data:
        circuitType = entry['circuit_type']
        numberOfConstraints = entry['number_of_constraints']
        witnessGenTime = entry['public_witness_gen_time']
        proofGenTime = entry['proof_gen_time']
        proofVerTime = entry['proof_verify_time']

        entry = CircuitConstraintResult(circuitType = circuitType,
                                        numberOfConstraints= numberOfConstraints,
                                        witnessGenTime = witnessGenTime,
                                        proofGenTime= proofGenTime,
                                        proofVerTime= proofVerTime)

        entries.append(entry)

    return entries

def plot_circuit_contstraint_result():
    circuit_constraint_result = parse_zkrevoke_circuit_constraint_result_entry("result_circuit_constraints_avg.json")

    witnessGenTimeEmptyCircuit = 0
    proofGenTimeEmptyCircuit = 0
    proofVerTimeEmptyCircuit = 0
    witnessGenTimeCompleteCircuit = 0
    proofGenTimeCompleteCircuit = 0
    proofVerTimeCompleteCircuit = 0

    signature_verification_constraints = 0
    token_verification_constraints = 0
    challenge_verification_constraints = 0
    complete_circuit_constraints = 0

    for entry in circuit_constraint_result:
        if entry.circuitType == "empty_circuit":
            witnessGenTimeEmptyCircuit = entry.witnessGenTime
            proofGenTimeEmptyCircuit = entry.proofGenTime
            proofVerTimeEmptyCircuit = entry.proofVerTime
        if entry.circuitType == "complete_circuit":
            witnessGenTimeCompleteCircuit = entry.witnessGenTime
            proofGenTimeCompleteCircuit = entry.proofGenTime
            proofVerTimeCompleteCircuit = entry.proofVerTime
            complete_circuit_constraints = entry.numberOfConstraints
        if entry.circuitType == "signature_verification":
            signature_verification_constraints = entry.numberOfConstraints
        if entry.circuitType == "challenge_verification":
            challenge_verification_constraints = entry.numberOfConstraints
        if entry.circuitType == "token_verification":
            token_verification_constraints = entry.numberOfConstraints


    col_labels = ["signature ver.", "token ver.", "challenge ver.", "complete circuit"]
    row_labels = ["No. of Constraints"]
    data = [[signature_verification_constraints, token_verification_constraints, challenge_verification_constraints, complete_circuit_constraints]]

    fig1, ax1 = plt.subplots()
    font = {'fontname': 'Times New Roman', 'weight': 'bold'}
    ax1.axis('off')
    ax1.axis('tight')
    table = ax1.table(cellText=data,
                     colLabels=col_labels,
                     rowLabels=row_labels,
                     loc='center',
                     cellLoc='center',
                     edges = 'closed')



    filename1 = "graphs/table_4:Circuit_Constraints_Breakdown"+".png"
    fig1.savefig(filename1, bbox_inches='tight')


    col_labels = ["Witness Gen. Time (µs)", "Proof Gen. Time (µs)", "Proof Ver. Time (µs)"]
    row_labels = ["Empty Circuit", "Complete Circuit"]
    data = [[witnessGenTimeEmptyCircuit, proofGenTimeEmptyCircuit, proofVerTimeEmptyCircuit],
            [witnessGenTimeCompleteCircuit, proofGenTimeCompleteCircuit, proofVerTimeCompleteCircuit]]

    fig, ax = plt.subplots()
    font = {'fontname': 'Times New Roman', 'weight': 'bold'}
    ax.axis('off')
    ax.axis('tight')
    table = ax.table(cellText=data,
                     colLabels=col_labels,
                     rowLabels=row_labels,
                     loc='center',
                     cellLoc='center',
                     edges = 'closed')



    filename = "graphs/table_5:Groth16_Overhead"+".png"
    fig.savefig(filename, bbox_inches='tight')



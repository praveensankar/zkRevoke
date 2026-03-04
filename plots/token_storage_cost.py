import math
import json
import os

from matplotlib.ticker import ScalarFormatter
from scipy.stats import sem
import numpy as np
from matplotlib import pyplot as plt

from revocation import parse_zkrevoke_refresh_result_entry

class TokenStorageResult:
    def __init__(self, number_of_tokens, size, gas):
        self.number_of_tokens = number_of_tokens
        self.size = size
        self.gas = gas



    def __eq__(self, another):
        return self.number_of_tokens == another.number_of_tokens

    def __hash__(self):
        return hash(str(self.number_of_tokens))


def parse_token_storage_result_entry(file):
    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, "benchmark", "results", file)

    with open(file_path) as f:
        json_data = json.load(f)

    entries = []
    for entry in json_data:
        number_of_tokens = int(entry['numberOfRevokedVCs'])
        size = int(entry['tokenSize']) * int(entry['numberOfRevokedVCs'])
        gas = int(entry['gas'])
        res = TokenStorageResult(
            number_of_tokens= number_of_tokens,
            size= size,
            gas=gas)

        entries.append(res)

    return entries

def plot_token_storage_cost_result():
    token_storage_cost_entry = parse_token_storage_result_entry("result_token_storage_avg.json")

    number_of_tokens = np.empty(6)
    size = np.empty(6)
    gas = np.empty(6)

    for entry in token_storage_cost_entry:
        if entry.number_of_tokens == 1:
            number_of_tokens[0]=1
            size[0] = entry.size
            gas[0] = entry.gas
        if entry.number_of_tokens == 10:
            number_of_tokens[1]=10
            size[1] = entry.size
            gas[1] = entry.gas
        if entry.number_of_tokens == 100:
            number_of_tokens[2]=100
            size[2] = entry.size
            gas[2] = entry.gas
        if entry.number_of_tokens == 1000:
            number_of_tokens[3]=1000
            size[3] = entry.size
            gas[3] = entry.gas
        if entry.number_of_tokens == 10000:
            number_of_tokens[4]=10000
            size[4] = entry.size
            gas[4] = entry.gas
        if entry.number_of_tokens == 100000:
            number_of_tokens[5]=100000
            size[5] = entry.size
            gas[5] = entry.gas


    col_labels = ["Number of Tokens", "Size (B)", "Gas"]

    num_rows = 3
    num_columns = 6

    data =  np.empty((6, 3))
    for i in range(6):
        data[i][0] = number_of_tokens[i]
        data[i][1] = size[i]
        data[i][2] = gas[i]


    fig1, ax1 = plt.subplots()
    font = {'fontname': 'Times New Roman', 'weight': 'bold'}
    ax1.axis('off')
    ax1.axis('tight')
    table = ax1.table(cellText=data,
                      colLabels=col_labels,
                      loc='center',
                      cellLoc='center',
                      edges = 'closed')



    filename1 = "graphs/table_6:The_Cost_of_Storing_Tokens_in_Smart_Contract"+".png"
    fig1.savefig(filename1, bbox_inches='tight')


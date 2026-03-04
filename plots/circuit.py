import json
from collections import defaultdict
from re import error

import math
import os

import numpy as np
from matplotlib import pyplot as plt
from matplotlib.ticker import ScalarFormatter
from scipy.stats import sem
from results import *
class CircuitResult:
    def __init__(self, numberOfTokensInCircuit, circuitSize, circuitGenTime, circuitGas, circuitCost, circuitConstraints):
        self.numberOfTokensInCircuit = numberOfTokensInCircuit
        self.circuitSize = circuitSize
        self.circuitGenTime = circuitGenTime
        self.circuitGas = circuitGas
        self.circuitCost = circuitCost
        self.circuitConstraints = circuitConstraints


    def __eq__(self, another):
        return self.numberOfTokensInCircuit == another.numberOfTokensInCircuit

    def __hash__(self):
        return hash(str(self.numberOfTokensInCircuit))


def parse_circuit_result_entry(file):
    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, "benchmark", "results", file)
    with open(file_path) as f:
        json_data = json.load(f)

    entries = []
    for results in json_data:
        circuitResults = results['zkpcircuitResult']

        for entry in circuitResults:
            numberOfTokensInCircuit = entry['numberOfTokensInCircuit']
            circuitSize = entry['grothCcsSize']
            circuitGenTime = entry['grothccsTime']
            circuitGas = entry['grothCcsGas']
            circuitCost = entry['GrothccsCost']
            circuitConstraints = entry['grothccsNumberOfConstraints']

            entry = CircuitResult(numberOfTokensInCircuit=numberOfTokensInCircuit,
                                  circuitSize=circuitSize,
                                  circuitGenTime=circuitGenTime,
                                  circuitGas=circuitGas,
                                  circuitCost=circuitCost,
                                  circuitConstraints=circuitConstraints)

            entries.append(entry)

    return entries




def plot_circuit_size(downsample_rate):
    circuit_entries = parse_circuit_result_entry("result_setup.json")
    irma_entries = parse_irma_setup_result_entry("result_setup.json")

    res_size = {}
    res_gen_time = {}
    res_gas = {}
    res_cost = {}
    error_size = {}
    error_gen_time = {}
    error_gas = {}
    error_cost = {}
    keys = set()


    for entry in circuit_entries:
        if entry.numberOfTokensInCircuit <=60:
            if entry.__hash__() in keys:
                size = np.asarray(entry.circuitSize)
                res_size[entry.numberOfTokensInCircuit] = np.append(res_size[entry.numberOfTokensInCircuit], size)
                circuit_generation_time = np.asarray(entry.circuitGenTime)
                res_gen_time[entry.numberOfTokensInCircuit] = np.append(res_gen_time[entry.numberOfTokensInCircuit], circuit_generation_time)
                gas = np.asarray(entry.circuitGas)
                res_gas[entry.numberOfTokensInCircuit] = np.append(res_gas[entry.numberOfTokensInCircuit], gas)
                cost = np.asarray(entry.circuitCost)
                res_cost[entry.numberOfTokensInCircuit] = np.append(res_cost[entry.numberOfTokensInCircuit], cost)
            else:
                res_size[entry.numberOfTokensInCircuit]= np.asarray(entry.circuitSize)
                res_gen_time[entry.numberOfTokensInCircuit] = np.asarray(entry.circuitGenTime)
                res_gas[entry.numberOfTokensInCircuit] = np.asarray(entry.circuitGas)
                res_cost[entry.numberOfTokensInCircuit] = np.asarray(entry.circuitCost)
                keys.add(entry.__hash__())


    for key, value in res_size.items():
        # print("zkRevoke size: number of tokens in a circuit: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        res_size[key] = int(np.mean(value))
        error_size[key] = np.std(value)

    res_size = dict(sorted(res_size.items()))
    xpoints = np.array(list(res_size.keys()))
    y1points = np.array(list(res_size.values()))
    y1points = np.ceil(y1points/1024)
    error_size = dict(sorted(error_size.items()))
    ey1points = np.array(list(error_size.values()))
    ey1points = np.ceil(ey1points/1024)

    for key, value in res_gen_time.items():
        # print("zkRevoke time: number of tokens in a circuit: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        res_gen_time[key] = int(np.mean(value)/1000)
        error_gen_time[key] = np.std(value)/1000


    res_gen_time = dict(sorted(res_gen_time.items()))
    y2points = np.array(list(res_gen_time.values()))
    # y2points = np.ceil(y2points/1000)
    error_gen_time = dict(sorted(error_gen_time.items()))
    ey2points = np.array(list(error_gen_time.values()))
    # ey2points = np.ceil(ey2points/1000)


    for key, value in res_gas.items():
        # print("zkRevoke gas: number of tokens in a circuit: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        res_gas[key] = int(np.mean(value))
        error_gas[key] = np.std(value)

    res_gas = dict(sorted(res_gas.items()))
    y3points = np.array(list(res_gas.values()))
    y3points = np.ceil(y3points/1000000)
    error_gas = dict(sorted(error_gas.items()))
    ey3points = np.array(list(error_gas.values()))
    ey2points = np.ceil(ey3points/1000000)

    for key, value in res_cost.items():
        res_cost[key] = int(np.mean(value))
        res_cost[key] = np.ceil(((np.mean(value))/1000000000000000000)*1880)
        error_cost[key] = np.std(value)
        error_cost[key] = np.ceil(((np.std(value))/1000000000000000000)*1880)
        # print("zkRevoke cost: number of tokens in a circuit: ", key, "\t cost: ", res_cost[key],"\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))


    res_cost = dict(sorted(res_cost.items()))
    y4points = np.array(list(res_cost.values()))
    # y3points = np.ceil(y2points/1000)
    error_cost = dict(sorted(error_cost.items()))
    ey4points = np.array(list(error_cost.values()))
    # ey2points = np.ceil(ey2points/1000)




    resIRMAPublicKeySize = np.empty(0)
    resIRMAKeyGenTime = np.empty(0)

    for entry in irma_entries:
        resIRMAPublicKeySize = np.append(resIRMAPublicKeySize, np.asarray(entry.publicKeySize))
        resIRMAKeyGenTime = np.append(resIRMAKeyGenTime, np.asarray(entry.keyGenTime))


    resIRMAPublicKeySize = resIRMAPublicKeySize/1024
    resIRMAPublicKeySize = np.mean(resIRMAPublicKeySize)
    # print("IRMA public key size: ", resIRMAPublicKeySize)
    resIRMAKeyGenTime = round(np.mean(resIRMAKeyGenTime)/1000,3)

    errorIRMAPublicKeySize = np.std(resIRMAPublicKeySize)
    errorIRMAKeyGenTime = round(np.std(resIRMAKeyGenTime)/1000,3)


    yIRMAPublicKeySizepoints = np.array([resIRMAPublicKeySize for i in range(len(xpoints))])
    yIRMAKeyGenTimepoints = np.array([resIRMAKeyGenTime for i in range(len(xpoints))])
    eyIRMAPublicKeySizepoints = np.array([errorIRMAPublicKeySize for i in range(len(xpoints))])
    eyIRMAKeyGenTimepoints = np.array([errorIRMAKeyGenTime for i in range(len(xpoints))])



    xpoints = xpoints[::downsample_rate]
    y1points = y1points[::downsample_rate]
    ey1points = ey1points[::downsample_rate]
    y2points = y2points[::downsample_rate]
    ey2points = ey2points[::downsample_rate]
    y3points = y3points[::downsample_rate]
    ey3points = ey3points[::downsample_rate]
    y4points = y4points[::downsample_rate]
    ey4points = ey4points[::downsample_rate]

    yIRMAPublicKeySizepoints = yIRMAPublicKeySizepoints[::downsample_rate]
    eyIRMAPublicKeySizepoints = eyIRMAPublicKeySizepoints[::downsample_rate]
    yIRMAKeyGenTimepoints = yIRMAKeyGenTimepoints[::downsample_rate]
    eyIRMAKeyGenTimepoints = eyIRMAKeyGenTimepoints[::downsample_rate]


    font = {'fontname': 'Times New Roman',  'weight': 'bold'}




    # y1labels = [str(i) for i in y1points]
    #
    # for i, label in enumerate(y1labels):
    #     # bbox_props = dict(boxstyle='square,pad=0.2', alpha=0.5)
    #     # ax.text(x1points[i], y1points[i], label)
    #     ax.annotate(label, (xpoints[i], y1points[i]), xytext=(-5, 8), textcoords='offset points',
    #                 arrowprops=dict(arrowstyle='<-'))
    fig1, ax1 = plt.subplots(layout='constrained')
    plt.tight_layout()
    for i in  range(len(xpoints)):
        # bbox_props = dict(boxstyle='square,pad=0.2', alpha=0.5)
        # ax.text(x1points[i], y1points[i], label)
        if i==0:
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(y1points[i]))), (xpoints[i], y1points[i]), xytext=(-5, -15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
            ax1.annotate((str(np.ceil(xpoints[i])), str(round(yIRMAPublicKeySizepoints[i], 2))), (xpoints[i], yIRMAPublicKeySizepoints[i]), xytext=(-5, 10), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
        if i==4:
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(y1points[i]))), (xpoints[i], y1points[i]), xytext=(-15, -15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
            ax1.annotate((str(np.ceil(xpoints[i])), str(round(yIRMAPublicKeySizepoints[i], 2))), (xpoints[i], yIRMAPublicKeySizepoints[i]), xytext=(-25, 10), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
        if i==9:
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(y1points[i]))), (xpoints[i], y1points[i]), xytext=(-35, -25), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
            ax1.annotate((str(np.ceil(xpoints[i])), str(round(yIRMAPublicKeySizepoints[i], 2))), (xpoints[i], yIRMAPublicKeySizepoints[i]), xytext=(-40, 10), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))

    ax1.errorbar(xpoints, y1points, color="#226363",  marker='d',  label=r'zkRevoke circuit size',  yerr=ey1points)
    ax1.errorbar(xpoints, yIRMAPublicKeySizepoints,  linestyle=(0, (5,1)), marker='x', color="red", label=r'IRMA public key size', yerr=eyIRMAPublicKeySizepoints)

    ax1.set_xlabel(r'The number of tokens per circuit: k', font, fontsize=14)
    ax1.set_ylabel(r'size (in KB)', font, fontsize=14)
    ax1.legend(fontsize="11", framealpha=0.3)
    fig1.set_size_inches(5, 3)
    fig1.savefig("graphs/fig_5a_result_circuit_size.png", bbox_inches='tight')

    fig2, ax2 = plt.subplots(layout='constrained')
    plt.tight_layout()
    for i in  range(len(xpoints)):
        # bbox_props = dict(boxstyle='square,pad=0.2', alpha=0.5)
        # ax.text(x1points[i], y1points[i], label)
        if i==0:
            ax2.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(y2points[i])))), (xpoints[i], y2points[i]), xytext=(-2, -15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
            ax2.annotate((str(np.ceil(xpoints[i])), str(round(yIRMAKeyGenTimepoints[i],3))), (xpoints[i], yIRMAKeyGenTimepoints[i]), xytext=(-2, 5), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
        if i==4:
            ax2.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(y2points[i])))), (xpoints[i], y2points[i]), xytext=(-5, -15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
            ax2.annotate((str(np.ceil(xpoints[i])), str(round(yIRMAKeyGenTimepoints[i],3))), (xpoints[i], yIRMAKeyGenTimepoints[i]), xytext=(-25, 5), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
        if i==9:
            ax2.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(y2points[i])))), (xpoints[i], y2points[i]), xytext=(-50, -45), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
            ax2.annotate((str(np.ceil(xpoints[i])), str(round(yIRMAKeyGenTimepoints[i],3))), (xpoints[i], yIRMAKeyGenTimepoints[i]), xytext=(-45, 5), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
    ax2.errorbar(xpoints, y2points, color="#226363",  marker='d',  label=r'zkRevoke circuit gen. time',  yerr=ey2points)
    ax2.errorbar(xpoints, yIRMAKeyGenTimepoints,  linestyle=(0, (5,1)), marker='x', color="red", label=r'IRMA key gen. time', yerr=eyIRMAKeyGenTimepoints)

    ax2.set_xlabel(r'The number of tokens per circuit: k', font, fontsize=14)
    ylabel_string = f"time (in ms)"
    ax2.set_ylabel(ylabel_string, font, fontsize=14)
    ax2.legend(fontsize="11", framealpha=0.3)
    fig2.set_size_inches(5, 3)
    fig2.savefig("graphs/fig_5b_result_circuit_time.png", bbox_inches='tight')
    #
    # fig3, ax3 = plt.subplots(layout='constrained')
    # for i in  range(len(xpoints)):
    #     # bbox_props = dict(boxstyle='square,pad=0.2', alpha=0.5)
    #     # ax.text(x1points[i], y1points[i], label)
    #     if i%10==0:
    #         ax3.annotate((str(np.ceil(xpoints[i])), str(np.ceil(y3points[i]))), (xpoints[i], y3points[i]), xytext=(-15, 2), textcoords='offset points',
    #                      arrowprops=dict(arrowstyle='->'))
    #
    # y3label = [str(i)+r'* $10^6$' for i in [100, 200, 300, 400, 500, 600, 700, 800, 900, 1000]]
    #
    # ax3.errorbar(xpoints, y3points, color="#226363",  yerr=ey3points)
    # ax3.set_xlabel(r'$k$: number of tokens in a circuit', font)
    # ax3.set_ylabel('Total gas consumption', font)
    # ax3.set_yticks([100, 200, 300, 400, 500, 600, 700, 800, 900, 1000], y3label)
    # # ax3.legend(fontsize="10")
    # fig3.savefig("graphs/result_circuit_gas.png")
    #
    #
    # fig4, ax4 = plt.subplots(layout='constrained')
    # for i in  range(len(xpoints)):
    #     # bbox_props = dict(boxstyle='square,pad=0.2', alpha=0.5)
    #     # ax.text(x1points[i], y1points[i], label)
    #     if i%10==0:
    #         ax4.annotate((str(np.ceil(xpoints[i])), str(np.ceil(y4points[i]))), (xpoints[i], y4points[i]), xytext=(-15, 2), textcoords='offset points',
    #                 arrowprops=dict(arrowstyle='->'))
    #
    # ax4.errorbar(xpoints, y4points, color="#226363",  yerr=ey4points)
    # ax4.set_xlabel(r'$k$: number of tokens in a circuit', font)
    # ax4.set_ylabel('Total cost (in USD)', font)
    # # ax4.legend(fontsize="10")
    # fig4.savefig("graphs/result_circuit_cost.png")
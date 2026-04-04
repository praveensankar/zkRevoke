
import json
from collections import defaultdict
from re import error

import math
import os

import numpy as np
from matplotlib import pyplot as plt
from matplotlib.ticker import ScalarFormatter
from scipy.stats import sem


class IRMARevocationResult:
    def __init__(self, totalEpochs, currentEpoch, totalRevokedVCsPerEpoch, timeToRevokeVCs, witnessUpdateSize, totalValidVCs, revocationRate, issuerBandwidth, totalHoldersWithValidVCs):
        self.totalEpochs = totalEpochs
        self.currentEpoch = currentEpoch
        self.totalRevokedVCsPerEpoch = totalRevokedVCsPerEpoch
        self.timeToRevokeVCs = timeToRevokeVCs
        self.witnessUpdateSize = witnessUpdateSize
        self.totalValidVCs = totalValidVCs
        self.revocationRate = revocationRate
        self.issuerBandwidth = issuerBandwidth
        self.totalHoldersWithValidVCs = totalHoldersWithValidVCs


    def __eq__(self, another):
        return self.totalEpochs == another.totalEpochs and self.currentEpoch == another.currentEpoch and self.totalValidVCs == another.totalValidVCs and self.revocationRate == another.revocationRate

    def __hash__(self):
        return hash(str(self.totalEpochs) + str(self.currentEpoch) + str(self.totalValidVCs) + str(self.revocationRate))


def parse_irma_revocation_result_entry(file):
    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, "irma", "benchmark", "results", file)

    with open(file_path) as f:
        json_data = json.load(f)

    entries = []
    for entry in json_data:
        totalEpochs = entry['total_epochs']
        currentEpoch = entry['current_epoch']
        totalRevokedVCsPerEpoch = entry['total_revoked_vcs_per_epoch']
        timeToRevokeVCs = entry['time_to_revoke_vcs']
        witnessUpdateSize = entry['witness_update_size']
        totalValidVCs = entry['total_valid_vcs']
        issuerBandwidth = entry['issuer_bandwidth']
        revocationRate = entry['revocation_rate']
        totalHoldersWithValidVCs = entry['total_holders_with_valid']
        entry = IRMARevocationResult(totalEpochs = totalEpochs,
                                     currentEpoch= currentEpoch,
                                     totalRevokedVCsPerEpoch = totalRevokedVCsPerEpoch,
                                     timeToRevokeVCs= timeToRevokeVCs,
                                     witnessUpdateSize= witnessUpdateSize,
                                     totalValidVCs= totalValidVCs,
                                     issuerBandwidth= issuerBandwidth,
                                     revocationRate= revocationRate,
                                     totalHoldersWithValidVCs= totalHoldersWithValidVCs)

        entries.append(entry)

    return entries


class ZKRevokeRevocationResult:
    def __init__(self, totalEpochs, currentEpoch, totalRevokedVCsPerEpoch, timeToRevokeVCs, totalValidVCs, revocationRate, issuerBandwidth):
        self.totalEpochs = totalEpochs
        self.currentEpoch = currentEpoch
        self.totalRevokedVCsPerEpoch = totalRevokedVCsPerEpoch
        self.timeToRevokeVCs = timeToRevokeVCs
        self.totalValidVCs = totalValidVCs
        self.revocationRate = revocationRate
        self.issuerBandwidth = issuerBandwidth



    def __eq__(self, another):
        return self.totalEpochs == another.totalEpochs and self.currentEpoch == another.currentEpoch and self.totalValidVCs == another.totalValidVCs and self.revocationRate == another.revocationRate

    def __hash__(self):
        return hash(str(self.totalEpochs) + str(self.currentEpoch) + str(self.totalValidVCs) + str(self.revocationRate))


def parse_zkrevoke_revocation_result_entry(file):
    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, "benchmark", "results", file)

    with open(file_path) as f:
        json_data = json.load(f)

    entries = []
    for entry in json_data:
        totalEpochs = entry['total_epochs']
        currentEpoch = entry['current_epoch']
        totalRevokedVCsPerEpoch = entry['total_revoked_vcs_per_epoch']
        timeToRevokeVCs = entry['time_to_revoke_vcs']
        totalValidVCs = entry['total_valid_vcs']
        issuerBandwidth = entry['issuer_bandwidth']
        revocationRate = entry['revocation_rate']
        entry = ZKRevokeRevocationResult(totalEpochs = totalEpochs,
                                     currentEpoch= currentEpoch,
                                     totalRevokedVCsPerEpoch = totalRevokedVCsPerEpoch,
                                     timeToRevokeVCs= timeToRevokeVCs,
                                     totalValidVCs= totalValidVCs,
                                     issuerBandwidth= issuerBandwidth,
                                     revocationRate= revocationRate)

        entries.append(entry)

    return entries

class ZKRevokeRefreshResult:
    def __init__(self, timeToRefresh, revocationRate, currentEpoch):
        self.timeToRefresh = timeToRefresh
        self.revocationRate = revocationRate
        self.currentEpoch = currentEpoch



    def __eq__(self, another):
        return self.revocationRate == another.revocationRate  and self.currentEpoch == another.currentEpoch

    def __hash__(self):
        return hash(str(self.revocationRate) + str(self.currentEpoch))


def parse_zkrevoke_refresh_result_entry(file):
    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, "benchmark", "results", file)

    with open(file_path) as f:
        json_data = json.load(f)

    entries = []
    for entry in json_data:
        timeToRefresh = entry['time']
        revocationRate = entry['revocation_rate']
        currentEpoch = entry['current_epoch']
        entry = ZKRevokeRefreshResult(
            timeToRefresh= timeToRefresh,
            revocationRate= revocationRate,
        currentEpoch=currentEpoch)

        entries.append(entry)

    return entries


# all the results should have the same total number of epochs
# revocationRate should be in percentage. e.g. 1%
def plot_issuer_bandwidth_revocation(totalVCs, revocationRate, downsample_rate):
    zkrevoke_entries = parse_zkrevoke_revocation_result_entry("result_revocation.json")
    irma_entries = parse_irma_revocation_result_entry("result_revocation.json")

    irma_entries_without_repetition = parse_irma_revocation_result_entry("result_revocation_witness_update_without_repetition.json")

    total_epochs = 0
    for entry in zkrevoke_entries:
        total_epochs = entry.totalEpochs
        break


    resZKRevoke = np.empty((100, total_epochs), dtype=object)
    keysZKRevoke = set()

    index = revocationRate - 1
    for entry in zkrevoke_entries:
        if entry.totalValidVCs == totalVCs:
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            bandwidth = np.asarray(entry.issuerBandwidth)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysZKRevoke:
                resZKRevoke[revocation_rate-1][current_epoch-1] = np.append(resZKRevoke[revocation_rate-1][current_epoch-1], bandwidth)
            else:
                resZKRevoke[revocation_rate-1][current_epoch-1] = bandwidth
                keysZKRevoke.add(entry.__hash__())


    num_rows = resZKRevoke.shape[0]
    num_columns = resZKRevoke.shape[1]

    ypointsZKRevoke =  np.empty((resZKRevoke.shape[0], resZKRevoke.shape[1]))
    errorZKRevoke =  np.empty((resZKRevoke.shape[0], resZKRevoke.shape[1]))
    max_revocation_rate = 0
    xpoints = np.array([i+1 for i in range(num_columns)])
    for i in range(num_rows):
        if  resZKRevoke[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(num_columns):
                ypointsZKRevoke[i][j] = int(np.mean(resZKRevoke[i][j])/1024)
                errorZKRevoke[i][j] = int(np.std(resZKRevoke[i][j])/1024)
                if i==1:
                    print("revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsZKRevoke[i][j], "\t std: ", errorZKRevoke[i][j], "\t bandwidth: ", resZKRevoke[i][j])



    resIRMA = np.empty((100, total_epochs), dtype=object)
    keysIRMA = set()


    for entry in irma_entries:
        if entry.totalValidVCs == totalVCs:
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            bandwidth = np.asarray(entry.issuerBandwidth)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysIRMA:
                resIRMA[revocation_rate-1][current_epoch-1] = np.append(resIRMA[revocation_rate-1][current_epoch-1], bandwidth)
            else:
                resIRMA[revocation_rate-1][current_epoch-1] = bandwidth
                keysIRMA.add(entry.__hash__())


    num_rows = resZKRevoke.shape[0]
    num_columns = resZKRevoke.shape[1]

    ypointsIRMA =  np.empty((resIRMA.shape[0], resIRMA.shape[1]))
    errorIRMA =  np.empty((resIRMA.shape[0], resIRMA.shape[1]))
    max_revocation_rate = 0
    xpoints = np.array([i+1 for i in range(num_columns)])
    for i in range(num_rows):
        if  resIRMA[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(num_columns):
                ypointsIRMA[i][j] = int(np.mean(resIRMA[i][j])/1024)
                errorIRMA[i][j] = int(np.std(resIRMA[i][j])/1024)
                if i==1:
                    print("revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsIRMA[i][j], "\t std: ", errorIRMA[i][j], "\t bandwidth: ", resIRMA[i][j])


    resIRMAWithoutRepition = np.empty((100, total_epochs), dtype=object)
    keysIRMAWithoutRepition = set()


    for entry in irma_entries_without_repetition:
        if entry.totalValidVCs == totalVCs:
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            bandwidth = np.asarray(entry.issuerBandwidth)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysIRMAWithoutRepition:
                resIRMAWithoutRepition[revocation_rate-1][current_epoch-1] = np.append(resIRMAWithoutRepition[revocation_rate-1][current_epoch-1], bandwidth)
            else:
                resIRMAWithoutRepition[revocation_rate-1][current_epoch-1] = bandwidth
                keysIRMAWithoutRepition.add(entry.__hash__())


    num_rows = resZKRevoke.shape[0]
    num_columns = resZKRevoke.shape[1]

    ypointsIRMAWithoutRepition =  np.empty((resIRMAWithoutRepition.shape[0], resIRMAWithoutRepition.shape[1]))
    errorIRMAWithoutRepition =  np.empty((resIRMAWithoutRepition.shape[0], resIRMAWithoutRepition.shape[1]))
    max_revocation_rate = 0
    xpoints = np.array([i+1 for i in range(num_columns)])
    for i in range(num_rows):
        if  resIRMAWithoutRepition[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(num_columns):
                ypointsIRMAWithoutRepition[i][j] = int(np.mean(resIRMAWithoutRepition[i][j])/1024)
                errorIRMAWithoutRepition[i][j] = int(np.std(resIRMAWithoutRepition[i][j])/1024)
                if i==1:
                    print("revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsIRMAWithoutRepition[i][j], "\t std: ", errorIRMAWithoutRepition[i][j], "\t bandwidth: ", resIRMAWithoutRepition[i][j])




    resIRMAWithoutRepitionWithSingleRegistry = np.empty((100, total_epochs), dtype=object)
    keysIRMAWithoutRepitionWithSingleRegistry = set()


    for entry in irma_entries_without_repetition:
        if entry.totalValidVCs == totalVCs:
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            bandwidth = np.asarray(entry.witnessUpdateSize)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysIRMAWithoutRepitionWithSingleRegistry:
                resIRMAWithoutRepitionWithSingleRegistry[revocation_rate-1][current_epoch-1] = np.append(resIRMAWithoutRepitionWithSingleRegistry[revocation_rate-1][current_epoch-1], bandwidth)
            else:
                resIRMAWithoutRepitionWithSingleRegistry[revocation_rate-1][current_epoch-1] = bandwidth
                keysIRMAWithoutRepitionWithSingleRegistry.add(entry.__hash__())


    num_rows = resZKRevoke.shape[0]
    num_columns = resZKRevoke.shape[1]

    ypointsIRMAWithoutRepitionWithSingleRegistry =  np.empty((resIRMAWithoutRepitionWithSingleRegistry.shape[0], resIRMAWithoutRepitionWithSingleRegistry.shape[1]))
    errorIRMAWithoutRepitionnWithSingleRegistry =  np.empty((resIRMAWithoutRepitionWithSingleRegistry.shape[0], resIRMAWithoutRepitionWithSingleRegistry.shape[1]))
    max_revocation_rate = 0
    xpoints = np.array([i+1 for i in range(num_columns)])
    for i in range(num_rows):
        if  resIRMAWithoutRepitionWithSingleRegistry[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(num_columns):
                ypointsIRMAWithoutRepitionWithSingleRegistry[i][j] = int(np.mean(resIRMAWithoutRepitionWithSingleRegistry[i][j])/1024)
                errorIRMAWithoutRepitionnWithSingleRegistry[i][j] = int(np.std(resIRMAWithoutRepitionWithSingleRegistry[i][j])/1024)
                if i==1:
                    print("revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsIRMAWithoutRepitionWithSingleRegistry[i][j], "\t std: ", errorIRMAWithoutRepitionnWithSingleRegistry[i][j], "\t bandwidth: ", resIRMAWithoutRepitionWithSingleRegistry[i][j])



    xpoints = xpoints[::downsample_rate]
    ypointsZKRevoke = ypointsZKRevoke[:, ::downsample_rate]
    errorZKRevoke = errorZKRevoke[:, ::downsample_rate]
    ypointsIRMA = ypointsIRMA[:, ::downsample_rate]
    errorIRMA = errorIRMA[:, ::downsample_rate]
    ypointsIRMAWithoutRepition = ypointsIRMAWithoutRepition[:, ::downsample_rate]
    errorIRMAWithoutRepition = errorIRMAWithoutRepition[:, ::downsample_rate]
    ypointsIRMAWithoutRepitionWithSingleRegistry = ypointsIRMAWithoutRepitionWithSingleRegistry[:, ::downsample_rate]
    errorIRMAWithoutRepitionnWithSingleRegistry = errorIRMAWithoutRepitionnWithSingleRegistry[:, ::downsample_rate]


    font = {'fontname': 'Times New Roman', 'weight': 'bold'}



    fig3, ax3 = plt.subplots(layout='constrained')
    plt.tight_layout()
    ax3.set_yscale('symlog', linthresh=1)


    for i in  range(len(ypointsZKRevoke[index])):
        # if revocationRate==1 and i==0:
        #     ax3.annotate((str(np.ceil(xpoints[i])), str(np.round(ypointsZKRevoke[index][i], 2))+"KB"), (xpoints[i], ypointsZKRevoke[index][i]), xytext=(5, -8), textcoords='offset points',
        #                  arrowprops=dict(arrowstyle='->'), fontsize=8)
        if i==5 or i==10:
            ax3.annotate((str(np.ceil(xpoints[i])), str(np.round(ypointsZKRevoke[index][i], 2))+"KB"), (xpoints[i], ypointsZKRevoke[index][i]), xytext=(-32, -15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=8)
    for i in  range(len(ypointsIRMA[index])):
        if   i==5 or i==10:
            ax3.annotate((str(np.ceil(xpoints[i])), str(np.round(ypointsIRMA[index][i]/(1024*1024), 2))+"GB"), (xpoints[i], ypointsIRMA[index][i]), xytext=(-40, -15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=8)
    for i in  range(len(ypointsIRMAWithoutRepition[index])):
        if   i==5 or i==10:
            ax3.annotate((str(np.ceil(xpoints[i])), str(np.round(ypointsIRMAWithoutRepition[index][i]/(1024*1024),2))+"GB"), (xpoints[i], ypointsIRMAWithoutRepition[index][i]), xytext=(-20, 6), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'),  fontsize=8)
    for i in  range(len(ypointsIRMAWithoutRepitionWithSingleRegistry[index])):
        if i==5 or i==10:
            ax3.annotate((str(np.ceil(xpoints[i])), str(np.round(ypointsIRMAWithoutRepitionWithSingleRegistry[index][i], 2))+"KB"), (xpoints[i], ypointsIRMAWithoutRepitionWithSingleRegistry[index][i]), xytext=(-20, -10), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=8)

    zkRevoke_string = f"zkRevoke:  $\mathcal{{R}}$= {revocationRate}%"
    irma_string = f"IRMA:  $\mathcal{{R}}$= {revocationRate}%"
    irma_without_repetition_string = f"IRMA-no repetition: $\mathcal{{R}}$= {revocationRate}%"
    irma_with_single_registry_string = f"IRMA-registry: $\mathcal{{R}}$= {revocationRate}%"
    ax3.errorbar(xpoints, ypointsZKRevoke[index], color="#226363", marker='+', label=zkRevoke_string, yerr=errorZKRevoke[index])
    ax3.errorbar(xpoints, ypointsIRMA[index],  linestyle=(0, (5,3)), marker='x', color="red", label=irma_string, yerr=errorIRMA[index])
    ax3.errorbar(xpoints, ypointsIRMAWithoutRepition[index],  linestyle=(0, (5,10)), marker='o', color="blue", label=irma_without_repetition_string, yerr=errorIRMAWithoutRepition[index])
    ax3.errorbar(xpoints, ypointsIRMAWithoutRepitionWithSingleRegistry[index],  linestyle=(0, (5,1)), marker='v', color="#ffb000", label=irma_with_single_registry_string, yerr=errorIRMAWithoutRepitionnWithSingleRegistry[index])
    marker_indices = [0]
    for i in range(len(ypointsZKRevoke[index])):
        if i%10==0:
            marker_indices.append(i)
    for idx in marker_indices:
        ax3.errorbar(xpoints[idx], ypointsZKRevoke[index][idx], yerr=errorZKRevoke[index][idx], fmt='+', color='#226363')
        ax3.errorbar(xpoints[idx], ypointsIRMA[index][idx], yerr=errorIRMA[index][idx], fmt='x',  color="red")
        ax3.errorbar(xpoints[idx], ypointsIRMAWithoutRepition[index][idx], yerr=errorIRMAWithoutRepition[index][idx], fmt='o',  color="blue")
        ax3.errorbar(xpoints[idx], ypointsIRMAWithoutRepitionWithSingleRegistry[index][idx], yerr=errorIRMAWithoutRepitionnWithSingleRegistry[index][idx], fmt='v',  color="#ffb000")


    ax3.set_xlabel(r'Epochs', font, fontsize=14)
    ax3.set_ylabel('Issuer: bandwidth (in KB)', font, fontsize=14)
    fig3.set_size_inches(3.8, 3.5)
    ax3.legend(fontsize="8.5", framealpha=0.3)
    if revocationRate==1:
        figNumber = "a"
    if revocationRate==5:
        figNumber = "b"

    if totalVCs!=1000000:
        for child in ax3.get_children():
            if isinstance(child, plt.Annotation):
                child.set_visible(False)

    filename = "graphs/fig_2"+figNumber+"_result_revocation_issuer_bandwidth_without_repition_r"+str(revocationRate)+".png"
    fig3.savefig(filename, bbox_inches='tight')





# all the results should have the same total number of epochs
# revocationRate should be in percentage. e.g. 1%
def plot_issuer_computation_revocation(totalVCs, revocationRate, downsample_rate):
    zkrevoke_entries = parse_zkrevoke_revocation_result_entry("result_revocation.json")
    zkrevoke_refresh_entries = parse_zkrevoke_refresh_result_entry("result_refresh.json")
    irma_entries = parse_irma_revocation_result_entry("result_revocation.json")

    irma_entries_without_repetition = parse_irma_revocation_result_entry("result_revocation_witness_update_without_repetition.json")

    total_epochs = 0
    for entry in zkrevoke_entries:
        total_epochs = entry.totalEpochs
        break

    index = revocationRate - 1

    resZKRevoke = np.empty((100, total_epochs), dtype=object)
    keysZKRevoke = set()


    for entry in zkrevoke_entries:
        if entry.totalValidVCs == totalVCs:
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            time = np.asarray(entry.timeToRevokeVCs)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysZKRevoke:
                resZKRevoke[revocation_rate-1][current_epoch-1] = np.append(resZKRevoke[revocation_rate-1][current_epoch-1], time)
            else:
                resZKRevoke[revocation_rate-1][current_epoch-1] = time
                keysZKRevoke.add(entry.__hash__())


    num_rows = resZKRevoke.shape[0]
    num_columns = resZKRevoke.shape[1]

    ypointsZKRevoke =  np.empty((resZKRevoke.shape[0], resZKRevoke.shape[1]))
    errorZKRevoke =  np.empty((resZKRevoke.shape[0], resZKRevoke.shape[1]))
    max_revocation_rate = 0
    xpoints = np.array([i+1 for i in range(num_columns)])
    for i in range(num_rows):
        if  resZKRevoke[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(num_columns):
                ypointsZKRevoke[i][j] = int(np.mean(resZKRevoke[i][j]))
                errorZKRevoke[i][j] = int(np.std(resZKRevoke[i][j]))
                if i==1:
                    print("revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsZKRevoke[i][j], "\t std: ", errorZKRevoke[i][j], "\t revocation time (in micro seconds): ", resZKRevoke[i][j])



    resZKRevokeRefresh = np.empty((100, total_epochs), dtype=object)
    keysZKRevokeRefresh = set()


    for entry in zkrevoke_refresh_entries:
        revocation_rate = entry.revocationRate
        if revocation_rate != 0 and revocation_rate==revocationRate:
            current_epoch = entry.currentEpoch
            time = np.asarray(entry.timeToRefresh)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysZKRevokeRefresh:
                resZKRevokeRefresh[revocation_rate-1][current_epoch-1] = np.append(resZKRevokeRefresh[revocation_rate-1][current_epoch-1], time)
            else:
                resZKRevokeRefresh[revocation_rate-1][current_epoch-1] = time
                keysZKRevokeRefresh.add(entry.__hash__())


    num_rows = resZKRevokeRefresh.shape[0]
    num_columns = resZKRevokeRefresh.shape[1]

    ypointsZKRevokeRefresh =  np.empty((resZKRevokeRefresh.shape[0], resZKRevokeRefresh.shape[1]))
    errorZKRevokeRefresh =  np.empty((resZKRevokeRefresh.shape[0], resZKRevokeRefresh.shape[1]))
    max_revocation_rate = 0
    xpoints = np.array([i+1 for i in range(num_columns)])
    for i in range(num_rows):
        if  resZKRevokeRefresh[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(num_columns):
                ypointsZKRevokeRefresh[i][j] = int(np.mean(resZKRevokeRefresh[i][j]))
                errorZKRevokeRefresh[i][j] = int(np.std(resZKRevokeRefresh[i][j]))
                if i==1:
                    print("revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsZKRevokeRefresh[i][j], "\t std: ", errorZKRevokeRefresh[i][j], "\t time (in micro seconds): ", resZKRevokeRefresh[i][j])




    resIRMA = np.empty((100, total_epochs), dtype=object)
    keysIRMA = set()


    for entry in irma_entries:
        if entry.totalValidVCs == totalVCs:
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            time = np.asarray(entry.timeToRevokeVCs)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysIRMA:
                resIRMA[revocation_rate-1][current_epoch-1] = np.append(resIRMA[revocation_rate-1][current_epoch-1], time)
            else:
                resIRMA[revocation_rate-1][current_epoch-1] = time
                keysIRMA.add(entry.__hash__())


    num_rows = resZKRevoke.shape[0]
    num_columns = resZKRevoke.shape[1]

    ypointsIRMA =  np.empty((resIRMA.shape[0], resIRMA.shape[1]))
    errorIRMA =  np.empty((resIRMA.shape[0], resIRMA.shape[1]))
    max_revocation_rate = 0
    xpoints = np.array([i+1 for i in range(num_columns)])
    for i in range(num_rows):
        if  resIRMA[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(num_columns):
                ypointsIRMA[i][j] = int(np.mean(resIRMA[i][j]))
                errorIRMA[i][j] = int(np.std(resIRMA[i][j]))
                if i==1:
                    print("revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsIRMA[i][j], "\t std: ", errorIRMA[i][j], "\t time (in micro seconds): ", resIRMA[i][j])





    font = {'fontname': 'Times New Roman', 'weight': 'bold'}




    title = "Total VCs: "
    if totalVCs==1000:
        title = title + "1K"
    if totalVCs == 10000:
        title = title+"10K"
    if totalVCs == 100000:
        title = title + "100K"
    if totalVCs== 1000000:
        title = title + "1M"



    xpoints = xpoints[::downsample_rate]
    ypointsZKRevoke = ypointsZKRevoke[:, ::downsample_rate]
    errorZKRevoke = errorZKRevoke[:, ::downsample_rate]
    ypointsIRMA = ypointsIRMA[:, ::downsample_rate]
    errorIRMA = errorIRMA[:, ::downsample_rate]
    ypointsZKRevokeRefresh = ypointsZKRevokeRefresh[:, ::downsample_rate]
    errorZKRevokeRefresh = errorZKRevokeRefresh[:, ::downsample_rate]


    fig3, ax3 = plt.subplots(layout='constrained')


    plt.tight_layout()
    # for i in  range(len(ypointsZKRevoke[index])):
    #     if  i==0 or i==5:
    #         ax3.annotate((str(np.ceil(xpoints[i])), str(np.round(ypointsZKRevoke[index][i]/1000, 2))), (xpoints[i], ypointsZKRevoke[index][i]/1000), xytext=(-15, 2), textcoords='offset points',
    #                      arrowprops=dict(arrowstyle='->'),  fontsize=9)
    #     if  i==10:
    #         ax3.annotate((str(np.ceil(xpoints[i])), str(np.round(ypointsZKRevoke[index][i]/1000, 2))), (xpoints[i], ypointsZKRevoke[index][i]/1000), xytext=(-25, 2), textcoords='offset points',
    #                      arrowprops=dict(arrowstyle='->'),  fontsize=9)
    for i in  range(len(ypointsIRMA[index])):
        if  i==0 or i==5:
            ax3.annotate((str(np.ceil(xpoints[i])), str(int(ypointsIRMA[index][i]/1000))+"$*10^3$"), (xpoints[i], ypointsIRMA[index][i]), xytext=(-5, 10), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'),  fontsize=9)
        if  i==10:
            ax3.annotate((str(np.ceil(xpoints[i])), str(int(ypointsIRMA[index][i]/1000))+"$*10^3$"), (xpoints[i], ypointsIRMA[index][i]), xytext=(-20, 10), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'),  fontsize=9)


    for i in range(num_rows):
        if  resZKRevoke[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(len(ypointsZKRevoke[i])):
                # print("i: ", i, "\t j: ", j, "\t bandwidth: ", ypointsIRMA[i][j])
                ypointsZKRevoke[i][j] = ypointsZKRevoke[i][j]/(1000)
                errorZKRevoke[i][j] = ypointsZKRevoke[i][j]/(1000)

    # for i in range(len(ypointsZKRevokeRefresh[index])):
    #     ypointsZKRevokeRefresh[index][i] = ypointsZKRevokeRefresh[index][i] + ypointsZKRevoke[index][i]
    #     errorZKRevokeRefresh[index][i] = errorZKRevokeRefresh[index][i] + errorZKRevoke[index][i]
    #     if errorZKRevokeRefresh[index][i] > 5000:
    #         print(i, ypointsZKRevokeRefresh[index][i], errorZKRevokeRefresh[index][i] )


    for i in  range(len(ypointsZKRevokeRefresh[index])):
        if i==0:
            ax3.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(ypointsZKRevokeRefresh[index][i]/1000)))), (xpoints[i], ypointsZKRevokeRefresh[index][i]), xytext=(15,2), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=9)
        if i==5:
            ax3.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(ypointsZKRevokeRefresh[index][i]/1000)))), (xpoints[i], ypointsZKRevokeRefresh[index][i]), xytext=(-50, 10), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=9)
        if i==10:
            ax3.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(ypointsZKRevokeRefresh[index][i]/1000)))), (xpoints[i], ypointsZKRevokeRefresh[index][i]), xytext=(-90, -10), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=9)




    zkRevoke_string = f"zkRevoke: $\mathcal{{R}}$= {revocationRate}%"
    zkRevoke_refresh_string = f"zkRevoke Revocation + Refresh: $\mathcal{{R}}$= {revocationRate}%"
    irma_string = f"IRMA: $\mathcal{{R}}$= {revocationRate}%"
    # ax3.errorbar(xpoints, ypointsZKRevoke[index], color="#226363", marker='+', label=zkRevoke_string, yerr=errorZKRevoke[index])
    ax3.errorbar(xpoints, ypointsIRMA[index],  linestyle=(0, (5,1)),  marker='x', color="red", label=irma_string, yerr=errorIRMA[index])
    ax3.errorbar(xpoints, ypointsZKRevokeRefresh[index],  marker='o', linestyle=(0, (5,5)), color="blue", label=zkRevoke_string, yerr=errorZKRevokeRefresh[index])
    # marker_indices = [0]
    # for i in range(len(ypointsZKRevoke[index])):
    #     if i%10==0:
    #         marker_indices.append(i)
    # for idx in marker_indices:
    #     ax3.errorbar(xpoints[idx], ypointsZKRevoke[index][idx], yerr=errorZKRevoke[index][idx], fmt='+', color='#226363')
    #     ax3.errorbar(xpoints[idx], ypointsIRMA[index][idx], yerr=errorIRMA[index][idx], fmt='*',  color="red")
    #     ax3.errorbar(xpoints[idx], ypointsZKRevokeRefresh[index][idx], yerr=errorZKRevokeRefresh[index][idx], fmt=',',  color="blue")



    end = ypointsZKRevokeRefresh[index][len(ypointsZKRevokeRefresh[0])-1]
    yticks_label = []

    yticks_locations = np.linspace(0, end, 6)
    for i in yticks_locations:
        yticks_label.append(str(int(np.ceil(np.ceil(i)/1000))))

    print(yticks_locations)
    print(yticks_label)

    ax3.set_yticks(yticks_locations, yticks_label)
    plt.tight_layout()
    ax3.set_xlabel(r'Epochs', font, fontsize=14)
    ax3.set_ylabel('Issuer: computation (in ms)', font, fontsize=14)

    if totalVCs!=1000000:
        ann = ax3.annotate('Text', xy=(1, 1))
    ann.remove()

    # ax3.set_title(title)
    ax3.legend(fontsize="9", framealpha=0.3)
    filename = "graphs/result_revocation_time_"+str(revocationRate)+".png"
    fig3.set_size_inches(4, 3)
    fig3.savefig(filename,  bbox_inches='tight')












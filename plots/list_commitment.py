
import json
from collections import defaultdict
from re import error

import math
import os

import numpy as np
from matplotlib import pyplot as plt
from matplotlib.ticker import ScalarFormatter
from scipy.stats import sem

from revocation import *


class ZKRevokeListCommitmentResult:
    def __init__(self, totalEpochs, currentEpoch, totalRevokedVCsPerEpoch, totalValidVCs, revocationRate, timeToCreateListCommitment, timeToVerifyListCommitment, sizeOfTheListAtTheCurrentEpoch):
        self.totalEpochs = totalEpochs
        self.currentEpoch = currentEpoch
        self.totalRevokedVCsPerEpoch = totalRevokedVCsPerEpoch
        self.totalValidVCs = totalValidVCs
        self.revocationRate = revocationRate
        self.timeToCreateListCommitment = timeToCreateListCommitment
        self.timeToVerifyListCommitment = timeToVerifyListCommitment
        self.sizeOfTheListAtTheCurrentEpoch = sizeOfTheListAtTheCurrentEpoch



    def __eq__(self, another):
        return self.totalEpochs == another.totalEpochs and self.currentEpoch == another.currentEpoch and self.totalValidVCs == another.totalValidVCs and self.revocationRate == another.revocationRate

    def __hash__(self):
        return hash(str(self.totalEpochs) + str(self.currentEpoch) + str(self.totalValidVCs) + str(self.revocationRate))


def parse_zkrevoke_list_commitment_result_entry(file):
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

        totalValidVCs = entry['total_valid_vcs']

        revocationRate = entry['revocation_rate']
        timeToCreateListCommitment = entry['time_to_create_commitment']
        timeToVerifyListCommitment = entry['time_to_verify_commitment']
        sizeOfTheListAtTheCurrentEpoch = entry['size_of_the_list_at_the_current_epoch']
        entry = ZKRevokeListCommitmentResult(totalEpochs = totalEpochs,
                                         currentEpoch= currentEpoch,
                                         totalRevokedVCsPerEpoch = totalRevokedVCsPerEpoch,
                                         totalValidVCs= totalValidVCs,
                                         revocationRate= revocationRate,
                                             timeToCreateListCommitment=timeToCreateListCommitment,
                                             timeToVerifyListCommitment=timeToVerifyListCommitment,
                                             sizeOfTheListAtTheCurrentEpoch=sizeOfTheListAtTheCurrentEpoch)

        entries.append(entry)

    return entries




# all the results should have the same total number of epochs
# revocationRate should be in percentage. e.g. 1%
def plot_list_commitment_verification_time(totalVCs, revocationRate, downsample_rate):
    zkrevoke_entries_commitment = parse_zkrevoke_list_commitment_result_entry("result_list_commitment.json")

    total_epochs = 0
    for entry in zkrevoke_entries_commitment:
        total_epochs = entry.totalEpochs
        break


    resZKRevokeCommitment = np.empty((100, total_epochs), dtype=object)
    keysZKRevokeCommitment = set()

    index = revocationRate - 1
    for entry in zkrevoke_entries_commitment:
        if entry.totalValidVCs == totalVCs:
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            timeToVerifyListCommitment = np.asarray(entry.timeToVerifyListCommitment)
            # print(current_epoch, revocation_rate, timeToVerifyListCommitment)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysZKRevokeCommitment:
                resZKRevokeCommitment[revocation_rate-1][current_epoch-1] = np.append(resZKRevokeCommitment[revocation_rate-1][current_epoch-1], timeToVerifyListCommitment)
            else:
                resZKRevokeCommitment[revocation_rate-1][current_epoch-1] = timeToVerifyListCommitment
                keysZKRevokeCommitment.add(entry.__hash__())


    num_rows = resZKRevokeCommitment.shape[0]
    num_columns = resZKRevokeCommitment.shape[1]

    ypointsZKRevokeCommitment =  np.empty((resZKRevokeCommitment.shape[0], resZKRevokeCommitment.shape[1]))
    errorZKRevokeCommitment =  np.empty((resZKRevokeCommitment.shape[0], resZKRevokeCommitment.shape[1]))
    max_revocation_rate = 0
    xpoints = np.array([i+1 for i in range(num_columns)])
    for i in range(num_rows):
        if  resZKRevokeCommitment[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(num_columns):
                ypointsZKRevokeCommitment[i][j] = int(np.mean(resZKRevokeCommitment[i][j]))
                errorZKRevokeCommitment[i][j] = int(np.std(resZKRevokeCommitment[i][j]))
                # if i==1:
                #     print("revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsZKRevokeCommitment[i][j], "\t std: ", errorZKRevokeCommitment[i][j], "\t timeToVerifyListCommitment (micro seconds): ", resZKRevokeCommitment[i][j])



    xpoints = xpoints[::downsample_rate]
    ypointsZKRevokeCommitment = ypointsZKRevokeCommitment[:, ::downsample_rate]
    errorZKRevokeCommitment = errorZKRevokeCommitment[:, ::downsample_rate]

    # print(ypointsZKRevokeCommitment)


    font = {'fontname': 'Times New Roman', 'weight': 'bold'}


    fig1, ax1 = plt.subplots(layout='constrained')
    plt.tight_layout()
    # ax1.set_yscale('log')
    for i in  range(len(ypointsZKRevokeCommitment[0])):
        if i==100 or i==300:
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsZKRevokeCommitment[index][i]))), (xpoints[i], ypointsZKRevokeCommitment[index][i]), xytext=(-40, -15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='<-'), fontsize=9)


    title = "Total VCs: "
    if totalVCs==1000:
        title = title + "1K"
    if totalVCs == 10000:
        title = title+"10K"
    if totalVCs == 100000:
        title = title + "100K"
    if totalVCs== 1000000:
        title = title + "1M"
    ax1.errorbar(xpoints, ypointsZKRevokeCommitment[0], color="#226363", label=r'zkRevoke: $\mathcal{R}$=1%', yerr=errorZKRevokeCommitment[0])
    ax1.errorbar(xpoints, ypointsZKRevokeCommitment[4], linestyle=(0, (1, 1)),  color="#581845",  label=r'zkRevoke: $\mathcal{R}$=5%', yerr=errorZKRevokeCommitment[4])
    ax1.errorbar(xpoints, ypointsZKRevokeCommitment[8], linestyle=(0, (3, 1, 1, 1, 1, 1)),  color="#1f4be1",  label=r'zkRevoke: $\mathcal{R}$=9%', yerr=errorZKRevokeCommitment[8])
    ax1.errorbar(xpoints, ypointsZKRevokeCommitment[12], linestyle=(0, (3, 5, 1, 5)),  color="#a93883",  label=r'zkRevoke: $\mathcal{R}$=13%', yerr=errorZKRevokeCommitment[12])

    marker_indices = [0]
    for i in range(len(ypointsZKRevokeCommitment[0])):
        if i%10==0:
            marker_indices.append(i)
    for idx in marker_indices:
        ax1.errorbar(xpoints[idx], ypointsZKRevokeCommitment[0][idx], yerr=errorZKRevokeCommitment[0][idx], fmt='+', color='#226363')
        ax1.errorbar(xpoints[idx], ypointsZKRevokeCommitment[4][idx], yerr=errorZKRevokeCommitment[4][idx], fmt='o', color='#581845')
        ax1.errorbar(xpoints[idx], ypointsZKRevokeCommitment[8][idx], yerr=errorZKRevokeCommitment[8][idx], fmt='s', color='#1f4be1')
        ax1.errorbar(xpoints[idx], ypointsZKRevokeCommitment[12][idx], yerr=errorZKRevokeCommitment[12][idx], fmt='p', color='#a93883')





    ax1.set_xlabel(r'Epochs', font, fontsize=14)
    ax1.set_ylabel('Verification Time (in µs)', font, fontsize=14)
    # ax1.set_title(title)
    fig1.set_size_inches(3.5, 3)
    ax1.legend(fontsize="9", framealpha=0.3, loc='upper left')
    filename = "graphs/fig_3_result_list_commitment_verification_time"+".png"
    fig1.savefig(filename, bbox_inches='tight')




# all the results should have the same total number of epochs
# revocationRate should be in percentage. e.g. 1%
def plot_issuer_computation_revocation_including_commitment(totalVCs, revocationRate, downsample_rate, include_refresh_in_commitment):
    zkrevoke_entries = parse_zkrevoke_revocation_result_entry("result_revocation.json")
    zkrevoke_refresh_entries = parse_zkrevoke_refresh_result_entry("result_refresh.json")
    irma_entries = parse_irma_revocation_result_entry("result_revocation.json")
    zkrevoke_entries_commitment = parse_zkrevoke_list_commitment_result_entry("result_list_commitment.json")
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
    # print("number of rows: ", num_rows, "\t number of columns: ", num_columns)

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
                # if i==0 and j%50==0:
                #     print("zkrevoke: revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsZKRevoke[i][j], "\t std: ", errorZKRevoke[i][j], "\t revocation time (in micro seconds): ", resZKRevoke[i][j])



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
                # if i==0 and j%25==0:
                #     print("zkRevoke-refresh: revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsZKRevokeRefresh[i][j], "\t std: ", errorZKRevokeRefresh[i][j], "\t time (in micro seconds): ", resZKRevokeRefresh[i][j])




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
                # if i==0 and j%50==0:
                #     print("IRMA: revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsIRMA[i][j], "\t std: ", errorIRMA[i][j], "\t time (in micro seconds): ", resIRMA[i][j])


    resZKRevokeCommitment = np.empty((100, total_epochs), dtype=object)
    keysZKRevokeCommitment = set()

    index = revocationRate - 1
    for entry in zkrevoke_entries_commitment:
        if entry.totalValidVCs == totalVCs:
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            timeToCreateListCommitment = np.asarray(entry.timeToCreateListCommitment)
            # print(current_epoch, revocation_rate, timeToCreateListCommitment)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysZKRevokeCommitment:
                resZKRevokeCommitment[revocation_rate-1][current_epoch-1] = np.append(resZKRevokeCommitment[revocation_rate-1][current_epoch-1], timeToCreateListCommitment)
            else:
                resZKRevokeCommitment[revocation_rate-1][current_epoch-1] = timeToCreateListCommitment
                keysZKRevokeCommitment.add(entry.__hash__())


    num_rows = resZKRevokeCommitment.shape[0]
    # print("number of rows: ", num_rows)
    num_columns = resZKRevokeCommitment.shape[1]

    ypointsZKRevokeCommitment =  np.empty((resZKRevokeCommitment.shape[0], resZKRevokeCommitment.shape[1]))
    errorZKRevokeCommitment =  np.empty((resZKRevokeCommitment.shape[0], resZKRevokeCommitment.shape[1]))
    max_revocation_rate = 0
    xpoints = np.array([i+1 for i in range(num_columns)])
    for i in range(num_rows):
        if  resZKRevokeCommitment[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(num_columns):
                if include_refresh_in_commitment == True:
                    ypointsZKRevokeCommitment[i][j] = int(np.mean(resZKRevokeCommitment[i][j])) + ypointsZKRevokeRefresh[i][j]
                    errorZKRevokeCommitment[i][j] = int(np.std(resZKRevokeCommitment[i][j])) + errorZKRevokeRefresh[i][j]
                else:
                    ypointsZKRevokeCommitment[i][j] = int(np.mean(resZKRevokeCommitment[i][j]))
                    errorZKRevokeCommitment[i][j] = int(np.std(resZKRevokeCommitment[i][j]))

                # if i==0 and j%50==0:
                #     print("zkRevoke: list-commitment: revocation rate: ", i, "\t current epoch: ", j, "\t comm: ",np.mean(resZKRevokeCommitment[i][j]), "\t refresh: ",ypointsZKRevokeRefresh[i][j], "\t mean: ",  ypointsZKRevokeCommitment[i][j], "\t std: ", errorZKRevokeCommitment[i][j], "\t timeToCreateListCommitment + refresh time(micro seconds): ", resZKRevokeCommitment[i][j])







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
    ypointsZKRevokeCommitment = ypointsZKRevokeCommitment[:, ::downsample_rate]
    errorZKRevokeCommitment = errorZKRevokeCommitment[:, ::downsample_rate]


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
        if  i==0:
            ax3.annotate((str(np.ceil(xpoints[i])), str(int(ypointsIRMA[index][i]/1000))), (xpoints[i], ypointsIRMA[index][i]), xytext=(-8, 15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'),  fontsize=9)
        if  i==5:
            ax3.annotate((str(np.ceil(xpoints[i])), str(int(ypointsIRMA[index][i]/1000))), (xpoints[i], ypointsIRMA[index][i]), xytext=(-20, 15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'),  fontsize=9)
        if  i==10:
            ax3.annotate((str(np.ceil(xpoints[i])), str(int(ypointsIRMA[index][i]/1000))), (xpoints[i], ypointsIRMA[index][i]), xytext=(-30, 15), textcoords='offset points',
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
            ax3.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(ypointsZKRevokeRefresh[index][i]/1000)))), (xpoints[i], ypointsZKRevokeRefresh[index][i]), xytext=(-50, 20), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=9)
        if include_refresh_in_commitment == True:
            if i==10:
                ax3.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(ypointsZKRevokeRefresh[index][i]/1000)))), (xpoints[i], ypointsZKRevokeRefresh[index][i]), xytext=(-25, -40), textcoords='offset points',
                             arrowprops=dict(arrowstyle='->'), fontsize=9)


    if include_refresh_in_commitment == False:
        for i in  range(len(ypointsZKRevokeCommitment[index])):
            if i==0:
                ax3.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(ypointsZKRevokeCommitment[index][i])))), (xpoints[i], ypointsZKRevokeCommitment[index][i]), xytext=(1,-10), textcoords='offset points',
                             arrowprops=dict(arrowstyle='->'), fontsize=9)
            if i==5:
                ax3.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(ypointsZKRevokeCommitment[index][i])))), (xpoints[i], ypointsZKRevokeCommitment[index][i]), xytext=(-50, -10), textcoords='offset points',
                             arrowprops=dict(arrowstyle='->'), fontsize=9)
            if i==10:
                ax3.annotate((str(np.ceil(xpoints[i])), str(int(np.ceil(ypointsZKRevokeCommitment[index][i])))), (xpoints[i], ypointsZKRevokeCommitment[index][i]), xytext=(-50, -10), textcoords='offset points',
                             arrowprops=dict(arrowstyle='->'), fontsize=9)




    zkRevoke_string = f"zkRevoke: $\mathcal{{R}}$= {revocationRate}%"
    zkRevoke_refresh_string = f"zkRevoke Revocation + Refresh: $\mathcal{{R}}$= {revocationRate}%"
    irma_string = f"IRMA: $\mathcal{{R}}$= {revocationRate}%"
    zkrevoke_commitment_string = f"zkRevoke-commitment: $\mathcal{{R}}$= {revocationRate}%"
    # ax3.errorbar(xpoints, ypointsZKRevoke[index], color="#226363", marker='+', label=zkRevoke_string, yerr=errorZKRevoke[index])
    ax3.errorbar(xpoints, ypointsIRMA[index],  linestyle=(0, (5,1)),  marker='x', color="red", label=irma_string, yerr=errorIRMA[index])
    ax3.errorbar(xpoints, ypointsZKRevokeCommitment[index],  marker='d', linestyle=(0, (5,1)),  color='#f3d00e', label=zkrevoke_commitment_string, yerr=errorZKRevokeCommitment[index])
    ax3.errorbar(xpoints, ypointsZKRevokeRefresh[index],  marker='o', linestyle=(0, (5,10)), color="blue", label=zkRevoke_string, yerr=errorZKRevokeRefresh[index])


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
        if i ==0:
            yticks_label.append(str(0))
        else:
            yticks_label.append(str(int(np.ceil(np.ceil(i)/1000))))

    # print(yticks_locations)
    # print(yticks_label)

    ax3.set_yticks(yticks_locations, yticks_label)
    plt.tight_layout()
    ax3.set_xlabel(r'Epochs', font, fontsize=14)
    ax3.set_ylabel('Issuer: computation (in ms)', font, fontsize=14)
    # ax3.set_title(title)
    ax3.legend(fontsize="9", framealpha=0.3)
    filename = "graphs/fig_2c_result_revocation_computation_including_commitment_"+str(revocationRate)+".png"
    fig3.set_size_inches(4, 3)
    fig3.savefig(filename,  bbox_inches='tight')










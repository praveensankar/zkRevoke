

import numpy as np
from matplotlib import pyplot as plt, ticker
from matplotlib.ticker import ScalarFormatter, LogFormatterMathtext
from scipy.stats import sem
from results import *



def plot_one_time_sharing_computation(totalVCs, onlyZKRevoke, downsample_rate):
    zkrevoke_entries = parse_zkrevoke_presentation_result_entry("result_presentation.json")
    irma_entries = parse_irma_presentation_and_verification_result_entry("result_presentation_verification.json")
    irma_entries_without_repetition = parse_irma_presentation_and_verification_result_entry("result_presentation_verification_witness_update_without_repetition.json")

    resK1 = {}
    resK2 = {}
    resK3 = {}
    resK4 = {}
    resK5 = {}
    resK6 = {}
    resK7 = {}
    resK8 = {}
    resK9 = {}
    resK10 = {}


    error = {}



    max_vp_validity = 0
    for entry in irma_entries:
        total_epochs = entry.totalEpochs
        max_vp_validity = entry.VPValidityPeriod
        break


    resIRMA = np.empty((100, max_vp_validity), dtype=object)
    keysIRMA = set()


    for entry in irma_entries:
        if entry.totalVCs == totalVCs:
            # current_epoch denotes the validity of VP from epoch 1 till current epoch
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            temp = entry.witnessUpdateTime + entry.proofGenTime
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysIRMA:
                resIRMA[revocation_rate-1][current_epoch-1] = np.append(resIRMA[revocation_rate-1][current_epoch-1], temp)
            else:
                resIRMA[revocation_rate-1][current_epoch-1] = temp
                keysIRMA.add(entry.__hash__())


    num_rows = resIRMA.shape[0]
    num_columns = resIRMA.shape[1]

    ypointsIRMA =  np.empty((resIRMA.shape[0], resIRMA.shape[1]))
    errorIRMA =  np.empty((resIRMA.shape[0], resIRMA.shape[1]))
    max_revocation_rate = 0
    xpoints = np.array([i+1 for i in range(num_columns)])
    for i in range(num_rows):
        if  resIRMA[i][1] is not None:
            max_revocation_rate = max_revocation_rate  + 1
            for j in range(num_columns):
                ypointsIRMA[i][j] = int(np.mean(resIRMA[i][j])/1000)
                errorIRMA[i][j] = np.std(resIRMA[i][j])/1000
                # if i==1:
                #     print("IRMA: revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsIRMA[i][j], "\t std: ", errorIRMA[i][j], "\t time: ", resIRMA[i][j])


    resIMRAOnlyProof = {}
    errorIMRAOnlyProof = {}
    irmaKeysOnlyProof = set()
    for entry in irma_entries:
        if entry.totalVCs == totalVCs:
            match entry.revocationRate:
                case 1:
                    temp =  entry.proofGenTime
                    if entry.__hash__() in irmaKeysOnlyProof:
                        resIMRAOnlyProof[entry.currentEpoch] = np.append(resIMRAOnlyProof[entry.currentEpoch], np.asarray(temp))
                    else:
                        resIMRAOnlyProof[entry.currentEpoch] = np.asarray(temp)
                        irmaKeysOnlyProof.add(entry.__hash__())

    for key, value in resIMRAOnlyProof.items():
        # print("IRMA: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resIMRAOnlyProof[key] = int(np.mean(value))
        errorIMRAOnlyProof[key] = np.std(value)


    resIMRAWithoutRepetition = {}
    errorIMRAWithoutRepetition = {}
    irmaKeysWithoutRepetition = set()
    for entry in irma_entries_without_repetition:
        if entry.totalVCs == totalVCs:
            match entry.revocationRate:
                case 1:
                    temp = entry.witnessUpdateTime + entry.proofGenTime
                    if entry.__hash__() in irmaKeysWithoutRepetition:
                        resIMRAWithoutRepetition[entry.currentEpoch] = np.append(resIMRAWithoutRepetition[entry.currentEpoch], np.asarray(temp))
                    else:
                        resIMRAWithoutRepetition[entry.currentEpoch] = np.asarray(temp)
                        irmaKeysWithoutRepetition.add(entry.__hash__())

    for key, value in resIMRAWithoutRepetition.items():
        resIMRAWithoutRepetition[key] = int(np.mean(value))
        errorIMRAWithoutRepetition[key] = np.std(value)
        # print("IRMA: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))


    resK1Keys = set()
    resK2Keys = set()
    resK3Keys = set()
    resK4Keys = set()
    resK5Keys = set()
    resK6Keys = set()
    resK7Keys = set()
    resK8Keys = set()
    resK9Keys = set()
    resK10Keys = set()



    for entry in zkrevoke_entries:
        if entry.VPValidity == 0:
            continue

        match entry.numberOfTokensInCircuit:
            case 1:
                if entry.__hash__() in resK1Keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK1[entry.VPValidity] = np.append(resK1[entry.VPValidity], totalProofGenTime)
                else:
                    resK1[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    resK1Keys.add(entry.__hash__())
            case 2:
                if entry.__hash__() in resK2Keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK2[entry.VPValidity] = np.append(resK2[entry.VPValidity], totalProofGenTime)
                else:
                    resK2[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    resK2Keys.add(entry.__hash__())

            case 4:
                if entry.__hash__() in resK3Keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK3[entry.VPValidity] = np.append(resK3[entry.VPValidity], totalProofGenTime)
                else:
                    resK3[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    resK3Keys.add(entry.__hash__())

            case 8:
                if entry.__hash__() in resK4Keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK4[entry.VPValidity] = np.append(resK4[entry.VPValidity], totalProofGenTime)
                else:
                    resK4[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    resK4Keys.add(entry.__hash__())

            case 16:

                if entry.__hash__() in resK5Keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK5[entry.VPValidity] = np.append(resK5[entry.VPValidity], totalProofGenTime)
                else:
                    resK5[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    resK5Keys.add(entry.__hash__())

            case 32:
                if entry.__hash__() in resK6Keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK6[entry.VPValidity] = np.append(resK6[entry.VPValidity], totalProofGenTime)
                else:
                    resK6[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    resK6Keys.add(entry.__hash__())

            case 64:
                if entry.__hash__() in resK7Keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK7[entry.VPValidity] = np.append(resK7[entry.VPValidity], totalProofGenTime)
                else:
                    resK7[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    resK7Keys.add(entry.__hash__())

            case 128:
                if entry.__hash__() in resK8Keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK8[entry.VPValidity] = np.append(resK8[entry.VPValidity], totalProofGenTime)
                else:
                    resK8[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    resK8Keys.add(entry.__hash__())

            case 256:
                if entry.__hash__() in resK9Keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK9[entry.VPValidity] = np.append(resK9[entry.VPValidity], totalProofGenTime)
                else:
                    resK9[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    resK9Keys.add(entry.__hash__())

            case 512:
                if entry.__hash__() in resK10Keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK10[entry.VPValidity] = np.append(resK10[entry.VPValidity], totalProofGenTime)
                else:
                    resK10[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    resK10Keys.add(entry.__hash__())

    errorK1 = {}
    for key, value in resK1.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK1[key] = int(np.mean(value))
        errorK1[key] = np.std(value)

    errorK2 = {}
    # print("\n\n**************** number of tokens: 2 ********************")
    for key, value in resK2.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK2[key] = int(np.mean(value))
        errorK2[key] = np.std(value)

    errorK3 = {}
    for key, value in resK3.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK3[key] = int(np.mean(value))
        errorK3[key] = np.std(value)

    errorK4 = {}
    for key, value in resK4.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK4[key] = int(np.mean(value))
        errorK4[key] = np.std(value)

    errorK5 = {}
    for key, value in resK5.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK5[key] = int(np.mean(value))
        errorK5[key] = np.std(value)

    errorK6 = {}
    for key, value in resK6.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK6[key] = int(np.mean(value))
        errorK6[key] = np.std(value)

    errorK7 = {}
    for key, value in resK7.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK7[key] = int(np.mean(value))
        errorK7[key] = np.std(value)

    errorK8 = {}
    for key, value in resK8.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK8[key] = int(np.mean(value))
        errorK8[key] = np.std(value)

    errorK9 = {}
    for key, value in resK9.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK9[key] = int(np.mean(value))
        errorK9[key] = np.std(value)

    errorK10 = {}
    for key, value in resK10.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK10[key] = int(np.mean(value))
        errorK10[key] = np.std(value)

    resK1 = dict(sorted(resK1.items()))
    resK2 = dict(sorted(resK2.items()))
    resK3 = dict(sorted(resK3.items()))
    resK4 = dict(sorted(resK4.items()))
    resK5 = dict(sorted(resK5.items()))
    resK6 = dict(sorted(resK6.items()))
    resK7 = dict(sorted(resK7.items()))
    resK8 = dict(sorted(resK8.items()))
    resK9 = dict(sorted(resK9.items()))
    resK10 = dict(sorted(resK10.items()))

    xpoints = np.array(list(resK1.keys()))
    y1points = np.ceil(np.array(list(resK1.values()))/1000)
    y2points = np.ceil(np.array(list(resK2.values()))/1000)
    y3points = np.ceil(np.array(list(resK3.values()))/1000)
    y4points = np.ceil(np.array(list(resK4.values()))/1000)
    y5points = np.ceil(np.array(list(resK5.values()))/1000)
    y6points = np.ceil(np.array(list(resK6.values()))/1000)
    y7points = np.ceil(np.array(list(resK7.values()))/1000)
    y8points = np.ceil(np.array(list(resK8.values()))/1000)
    y9points = np.ceil(np.array(list(resK9.values()))/1000)
    y10points = np.ceil(np.array(list(resK10.values()))/1000)

    yIRMApointsWithoutRepetition = np.ceil(np.array(list(resIMRAWithoutRepetition.values()))/1000)
    yIRMApointsOnlyProof = np.ceil(np.array(list(resIMRAOnlyProof.values()))/1000)

    errorK1 = dict(sorted(errorK1.items()))
    ey1points = np.array(list(errorK1.values()))
    ey1points = np.ceil(ey1points/1000)
    errorK2 = dict(sorted(errorK2.items()))
    ey2points = np.array(list(errorK2.values()))
    ey2points = np.ceil(ey2points/1000)
    errorK3 = dict(sorted(errorK3.items()))
    ey3points = np.array(list(errorK3.values()))
    ey3points = np.ceil(ey3points/1000)
    errorK4 = dict(sorted(errorK4.items()))
    ey4points = np.array(list(errorK4.values()))
    ey4points = np.ceil(ey4points/1000)
    errorK5 = dict(sorted(errorK5.items()))
    ey5points = np.array(list(errorK5.values()))
    ey5points = np.ceil(ey5points/1000)
    errorK6 = dict(sorted(errorK6.items()))
    ey6points = np.array(list(errorK6.values()))
    ey6points = np.ceil(ey6points/1000)
    errorK7 = dict(sorted(errorK7.items()))
    ey7points = np.array(list(errorK7.values()))
    ey7points = np.ceil(ey7points/1000)
    errorK8 = dict(sorted(errorK8.items()))
    ey8points = np.array(list(errorK8.values()))
    ey8points = np.ceil(ey8points/1000)
    errorK9 = dict(sorted(errorK9.items()))
    ey9points = np.array(list(errorK9.values()))
    ey9points = np.ceil(ey9points/1000)
    errorK10 = dict(sorted(errorK10.items()))
    ey10points = np.array(list(errorK10.values()))
    ey10points = np.ceil(ey10points/1000)




    errorIMRAWithoutRepetition = dict(sorted(errorIMRAWithoutRepetition.items()))
    eyIMRAWithoutRepetitionpoints = np.array(list(errorIMRAWithoutRepetition.values()))
    eyIMRAWithoutRepetitionpoints = np.ceil(eyIMRAWithoutRepetitionpoints/1000)


    errorIMRAOnlyProof = dict(sorted(errorIMRAOnlyProof.items()))
    eyIMRAOnlyProofpoints = np.array(list(errorIMRAOnlyProof.values()))
    eyIMRAOnlyProofpoints = np.ceil(eyIMRAOnlyProofpoints/1000)

    font = {'fontname': 'Times New Roman', 'weight': 'bold'}


    xpoints = xpoints[::downsample_rate]
    y1points = y1points[::downsample_rate]
    ey1points = ey1points[::downsample_rate]
    y2points = y2points[::downsample_rate]
    ey2points = ey2points[::downsample_rate]
    y3points = y3points[::downsample_rate]
    ey3points = ey3points[::downsample_rate]
    y4points = y4points[::downsample_rate]
    ey4points = ey4points[::downsample_rate]
    y4points = y4points[::downsample_rate]
    ey4points = ey4points[::downsample_rate]
    y5points = y5points[::downsample_rate]
    ey5points = ey5points[::downsample_rate]
    y6points = y6points[::downsample_rate]
    ey6points = ey6points[::downsample_rate]


    yIRMApointsWithoutRepetition = yIRMApointsWithoutRepetition[::downsample_rate]
    eyIMRAWithoutRepetitionpoints = eyIMRAWithoutRepetitionpoints[::downsample_rate]

    yIRMApointsOnlyProof = yIRMApointsOnlyProof[::downsample_rate]
    eyIMRAOnlyProofpoints = eyIMRAOnlyProofpoints[::downsample_rate]

    ypointsIRMA = ypointsIRMA[:, ::downsample_rate]
    errorIRMA = errorIRMA[:, ::downsample_rate]




    fig1, ax1 = plt.subplots(layout='constrained')
    plt.tight_layout()


    if onlyZKRevoke==False:
        for i in  range(len(ypointsIRMA[0])):
            if  i==20 or i==40:
                ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(y1points[i]))), (xpoints[i], y1points[i]), xytext=(-15, -15), textcoords='offset points',
                             arrowprops=dict(arrowstyle='->'))
                ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsIRMA[0][i]))), (xpoints[i], ypointsIRMA[0][i]), xytext=(-40, 5), textcoords='offset points',
                             arrowprops=dict(arrowstyle='->'))
                ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(yIRMApointsOnlyProof[i]))), (xpoints[i], yIRMApointsOnlyProof[i]), xytext=(-10, -10), textcoords='offset points',
                             arrowprops=dict(arrowstyle='->'))
        # ax.set_xscale('log')
    # ax.set_yscale('log')
    x = np.arange(9)
    # xlabel = [str('0'), str(r'$2^{1}$'), str(r'$2^{2}$'), str(r'$2^{3}$'), str(r'$2^{4}$'), str(r'$2^{5}$'), str(r'$2^{6}$'), str(r'$2^{7}$'), str(r'$2^{8}$'), str(r'$2^{9}$'), str(r'$2^{10}$')]
    # ax.set_xticks([0, 2, 4, 8, 16, 32,64,128,256,512, 1024], xlabel)


    title = "Total VCs: "
    if totalVCs==1000:
        title = title + "1K"
    if totalVCs == 10000:
        title = title+"10K"
    if totalVCs == 100000:
        title = title + "100K"
    if totalVCs== 1000000:
        title = title + "1M"



    if onlyZKRevoke==True:
        # print(xpoints)
        # print(y3points)
        ax1.errorbar(xpoints, y3points, yerr= ey3points, color="#1fd8e1", marker='v', label=r'zkRevoke: $\it{k}=2^{2}$')
        ax1.errorbar(xpoints, y4points,  yerr= ey4points, color="#1f4be1",  marker='o', label=r'zkRevoke: $\it{k}=2^{3}$')
        ax1.errorbar(xpoints, y5points,yerr= ey5points,  color="#e1a61f",  marker='s', label=r'zkRevoke: $\it{k}=2^{4}$')
        ax1.errorbar(xpoints, y6points,yerr= ey6points, color='#921fe1',marker='p', label=r'zkRevoke: $\it{k}=2^{5}$')
        ax1.errorbar(xpoints, ypointsIRMA[0], yerr=errorIRMA[0], linestyle=(0, (5,10)), marker='*', color="red", label=r'IRMA: $\mathcal{R}$=1%')
        ax1.errorbar(xpoints, yIRMApointsWithoutRepetition, yerr=eyIMRAWithoutRepetitionpoints, linestyle=(0, (5,1)), marker=',', color="blue", label=r'IRMA-no repetition: $\mathcal{R}$=1%')
    else:
        # print(y1points)
        # print(ypointsIRMA[0])
        ax1.errorbar(xpoints, y1points, yerr=ey1points, color="#226363",  marker='d',  label=r'zkRevoke')
        ax1.errorbar(xpoints, ypointsIRMA[0], yerr=errorIRMA[0], linestyle=(0, (5,2)), marker='*', color="red", label=r'IRMA: $\mathcal{R}$=1%')
        ax1.errorbar(xpoints, yIRMApointsWithoutRepetition, yerr=eyIMRAWithoutRepetitionpoints, linestyle=(0, (5,1)), marker='v', color="blue", label=r'IRMA-no repetition: $\mathcal{R}$=1%')
        ax1.errorbar(xpoints, yIRMApointsOnlyProof, yerr=eyIMRAOnlyProofpoints, linestyle=(0, (5,5)), marker='p', color="#cc79a7", label=r'IRMA only proof: $\mathcal{R}$=1%')

    # plt.errorbar(xpoints, yIRMApoints, color="red", marker='X', label=r'IRMA', yerr=eyIRMApoints)


    # ax1.set_title(title)
    ax1.set_xlabel(r'VP verification period:  m (in epochs)', font, fontsize=14)
    ax1.set_ylabel('Holder: computation (in ms) ', font, fontsize=14)
    ax1.legend(fontsize="9", framealpha=0.3)
    fig1.set_size_inches(3.5, 3)

    if onlyZKRevoke==True:
        fig1.savefig("graphs/result_one_time_sharing_computation_zkRevoke.png", bbox_inches='tight')
    else:
        fig1.savefig("graphs/fig_1a_result_one_time_sharing_computation.png", bbox_inches='tight')


    fig2, ax2 = plt.subplots(layout='constrained')
    plt.tight_layout()

    ax2.errorbar(xpoints, y1points, yerr=ey1points, color="#226363",  marker='d',  label=r'zkRevoke')
    ax2.errorbar(xpoints, ypointsIRMA[0], yerr=errorIRMA[0], linestyle=(0, (5,2)), marker='*', color="red", label=r'IRMA: $\mathcal{R}$=1%')
    ax2.errorbar(xpoints, ypointsIRMA[4], yerr=errorIRMA[4], linestyle=(0, (5,1)), marker='o', color="blue", label=r'IRMA: $\mathcal{R}$=5%')
    ax2.errorbar(xpoints, ypointsIRMA[8], yerr=errorIRMA[8], linestyle=(0, (5,5)), marker='p', color="#dc267f", label=r'IRMA: $\mathcal{R}$=9%')
    ax2.errorbar(xpoints, ypointsIRMA[12], yerr=errorIRMA[12],  color="#ffb000", marker='v', label=r'IRMA: $\mathcal{R}$=13%')
    ax2.set_xlabel(r'VP verification period:  m (in epochs)', font, fontsize=14)
    ax2.set_ylabel('Holder: computation (in ms) ', font, fontsize=14)
    ax2.legend(fontsize="9", framealpha=0.3)
    fig2.set_size_inches(3.5, 3)
    # fig2.savefig("graphs/result_one_time_sharing_computation_irma.png", bbox_inches='tight')

    fig3, ax3 = plt.subplots(layout='constrained')
    plt.tight_layout()

    ax3.errorbar(xpoints, y2points, yerr=ey2points, color="#226363",  marker='d',  label=r'zkRevoke: $\it{k}=2$')
    ax3.errorbar(xpoints, y3points, yerr=ey3points, color="#226363",  marker='o',  label=r'zkRevoke: $\it{k}=2^2$')
    ax3.errorbar(xpoints, ypointsIRMA[0], yerr=errorIRMA[0], linestyle=(0, (5,2)), marker='*', color="red", label=r'IRMA: $\mathcal{R}$=1%')
    ax3.errorbar(xpoints, ypointsIRMA[4], yerr=errorIRMA[4], linestyle=(0, (5,1)), marker=',', color="blue", label=r'IRMA: $\mathcal{R}$=5%')
    ax3.errorbar(xpoints, ypointsIRMA[8], yerr=errorIRMA[8], linestyle=(0, (5,5)), marker='p', color="#dc267f", label=r'IRMA: $\mathcal{R}$=9%')
    ax3.errorbar(xpoints, ypointsIRMA[12], yerr=errorIRMA[12], color="#ffb000", marker='v', label=r'IRMA: $\mathcal{R}$=13%')
    ax3.set_xlabel(r'VP verification period:  m (in epochs)', font, fontsize=14)
    ax3.set_ylabel('Holder: computation (in ms) ', font, fontsize=14)
    ax3.legend(fontsize="9", framealpha=0.3)
    fig3.set_size_inches(3.5, 3)
    fig3.savefig("graphs/fig_1b_result_one_time_sharing_computation_irma_k_2.png", bbox_inches='tight')





# all the results should have the same total number of epochs
def plot_one_time_sharing_holder_bandwidth(totalVCs, downsample_rate):
    irma_entries = parse_irma_presentation_and_verification_result_entry("result_presentation_verification.json")
    irma_entries_without_repetition = parse_irma_presentation_and_verification_result_entry("result_presentation_verification_witness_update_without_repetition.json")

    total_epochs = 0
    max_vp_validity = 0
    for entry in irma_entries:
        total_epochs = entry.totalEpochs
        max_vp_validity = entry.VPValidityPeriod
        break


    ypointsZKRevoke =  np.array([((164+32)*(i+1))/1024 for i in range(max_vp_validity)])


    resIRMA = np.empty((100, max_vp_validity), dtype=object)
    keysIRMA = set()


    for entry in irma_entries:
        if entry.totalVCs == totalVCs:
            # current_epoch denotes the validity of VP from epoch 1 till current epoch
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            bandwidth = np.asarray(entry.holderBandwidth)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysIRMA:
                resIRMA[revocation_rate-1][current_epoch-1] = np.append(resIRMA[revocation_rate-1][current_epoch-1], bandwidth)
            else:
                resIRMA[revocation_rate-1][current_epoch-1] = bandwidth
                keysIRMA.add(entry.__hash__())


    num_rows = resIRMA.shape[0]
    num_columns = resIRMA.shape[1]

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
                # if i==1:
                #     print("IRMA: revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsIRMA[i][j], "\t std: ", errorIRMA[i][j], "\t bandwidth: ", resIRMA[i][j])


    resIRMAWithoutRepition  = np.empty((100, max_vp_validity), dtype=object)
    keysIRMAWithoutRepition  = set()


    for entry in irma_entries_without_repetition:
        if entry.totalVCs == totalVCs:
            # current_epoch denotes the validity of VP from epoch 1 till current epoch
            current_epoch = entry.currentEpoch
            revocation_rate = entry.revocationRate
            bandwidth = np.asarray(entry.holderBandwidth)
            # print("current_epoch: ", current_epoch, "\t revocation_rate: ", revocation_rate, "\t bandwidth: ", bandwidth)
            if entry.__hash__() in keysIRMAWithoutRepition:
                resIRMAWithoutRepition[revocation_rate-1][current_epoch-1] = np.append(resIRMAWithoutRepition[revocation_rate-1][current_epoch-1], bandwidth)
            else:
                resIRMAWithoutRepition[revocation_rate-1][current_epoch-1] = bandwidth
                keysIRMAWithoutRepition.add(entry.__hash__())


    num_rows = resIRMAWithoutRepition.shape[0]
    num_columns = resIRMAWithoutRepition.shape[1]

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
                # if i==1:
                    # print("IRMA: revocation rate: ", i, "\t current epoch: ", j, "\t mean: ",  ypointsIRMAWithoutRepition[i][j], "\t std: ", errorIRMAWithoutRepition[i][j], "\t bandwidth: ", resIRMAWithoutRepition[i][j])


    font = {'fontname': 'Times New Roman',  'weight': 'bold'}


    xpoints = xpoints[::downsample_rate]
    ypointsZKRevoke = ypointsZKRevoke[::downsample_rate]

    ypointsIRMA = ypointsIRMA[:, ::downsample_rate]
    errorIRMA = errorIRMA[:, ::downsample_rate]
    ypointsIRMAWithoutRepition = ypointsIRMAWithoutRepition[:, ::downsample_rate]
    errorIRMAWithoutRepition = errorIRMAWithoutRepition[:, ::downsample_rate]

    fig1, ax1 = plt.subplots(layout='constrained')
    plt.tight_layout()
    ax1.set_yscale('linear')
    # formatter = ScalarFormatter(useMathText=True)
    # formatter.set_scientific(True)
    # formatter.set_powerlimits((0, 0))
    #
    # ax1.yaxis.set_major_formatter(formatter)


    # for i in  range(len(ypointsIRMA[0])):
    #     if i==0:
    #         ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsZKRevoke[i]))), (xpoints[i], ypointsZKRevoke[i]), xytext=(-5, -10), textcoords='offset points',
    #                      arrowprops=dict(arrowstyle='<-'))
    #     if i==4:
    #         ax1.annotate((str(np.ceil(xpoints[i])), str(round(np.float32(ypointsZKRevoke[i]/1024),2))+"KB"), (xpoints[i], ypointsZKRevoke[i]), xytext=(-35, -10), textcoords='offset points',
    #                      arrowprops=dict(arrowstyle='<-'))
    #     if i==8:
    #         ax1.annotate((str(np.ceil(xpoints[i])), str(round(np.float32(ypointsZKRevoke[i]/1024),2))+"KB"), (xpoints[i], ypointsZKRevoke[i]), xytext=(-35, -10), textcoords='offset points',
    #                      arrowprops=dict(arrowstyle='<-'))
    # for i in  range(len(ypointsIRMAWithoutRepition[0])):
    #     if i==4:
    #         ax1.annotate((str(np.ceil(xpoints[i])), str(round(np.float32(ypointsIRMAWithoutRepition[0][i]/1024),2))+"KB"), (xpoints[i], ypointsIRMAWithoutRepition[0][i]), xytext=(-35, 12), textcoords='offset points',
    #                      arrowprops=dict(arrowstyle='<-'))
    #     if i==8:
    #         ax1.annotate((str(np.ceil(xpoints[i])), str(round(np.float32(ypointsIRMAWithoutRepition[0][i]/1024),2))+"KB"), (xpoints[i], ypointsIRMAWithoutRepition[0][i]), xytext=(-42, 12), textcoords='offset points',
    #                      arrowprops=dict(arrowstyle='<-'))
    # for i in  range(len(ypointsIRMA[0])):
    #     if i==4:
    #         ax1.annotate((str(np.ceil(xpoints[i])), str(round(np.float32(ypointsIRMA[0][i]/1024),2))+"KB"), (xpoints[i], ypointsIRMA[0][i]), xytext=(-35, 12), textcoords='offset points',
    #                      arrowprops=dict(arrowstyle='->'))
    #     if i==8:
    #         ax1.annotate((str(np.ceil(xpoints[i])), str(round(np.float32(ypointsIRMA[0][i]/1024),2))+"KB"), (xpoints[i], ypointsIRMA[0][i]), xytext=(-47, -15), textcoords='offset points',
    #                      arrowprops=dict(arrowstyle='->'))

    title = "Total VCs: "
    if totalVCs==1000:
        title = title + "1K"
    if totalVCs == 10000:
        title = title+"10K"
    if totalVCs == 100000:
        title = title + "100K"
    if totalVCs== 1000000:
        title = title + "1M"


    ax1.errorbar(xpoints, ypointsZKRevoke, color="#226363",  marker='d',  label=r'zkRevoke')
    ax1.errorbar(xpoints, ypointsIRMAWithoutRepition[0],  linestyle=(0, (5,1)), color="blue", marker='v', label=r'IRMA-no repetition: $\mathcal{R}$=1%', yerr=errorIRMAWithoutRepition[0])
    ax1.errorbar(xpoints, ypointsIRMA[0],  linestyle=(0, (5,10)), marker='*', color="red", label=r'IRMA: $\mathcal{R}$=1%', yerr=errorIRMA[0])
    ax1.set_xlabel(r'VP verification period: m (in epochs)',  font, fontsize=14)
    ax1.set_ylabel('Holder: bandwidth (in KB)', font, fontsize=14)
    # ax1.set_title(title)
    fig1.set_size_inches(3.5, 3)
    ax1.legend(fontsize="9", framealpha=0.3)
    filename = "graphs/fig_1c_result_one_time_sharing_holder_bandwidth"+".png"
    fig1.savefig(filename, bbox_inches='tight')

    #
    #
    # for i in range(num_rows):
    #     if  resIRMA[i][1] is not None:
    #         max_revocation_rate = max_revocation_rate  + 1
    #         for j in range(len(ypointsIRMA[i])):
    #             print("i: ", i, "\t j: ", j, "\t bandwidth: ", ypointsIRMA[i][j])
    #             ypointsIRMA[i][j] = ypointsIRMA[i][j]
    #             errorIRMA[i][j] = errorIRMA[i][j]
    #
    # fig2, ax2 = plt.subplots(layout='constrained')
    #
    #
    #
    #
    # ax2.errorbar(xpoints, ypointsZKRevoke, color="#226363",  marker='d',  label=r'zkRevoke')
    # ax2.errorbar(xpoints, ypointsIRMA[0],  color="red", marker="*", linestyle=(0, (5,5)),  label=r'IRMA: $\mathcal{R}$=1%', yerr=errorIRMA[0])
    # ax2.errorbar(xpoints, ypointsIRMA[4], color='#581845', marker="o", linestyle=(0, (5, 10)),   label=r'IRMA: $\mathcal{R}$=5%', yerr=errorIRMA[4])
    # ax2.errorbar(xpoints, ypointsIRMA[8], color='#f3d00e', marker="s", linestyle=(0, (5, 1)), label=r'IRMA: $\mathcal{R}$=9%', yerr=errorIRMA[8])
    # ax2.errorbar(xpoints, ypointsIRMA[12],  color='#0e38f3',marker="p",  linestyle = ":", label=r'IRMA: $\mathcal{R}$=13%', yerr=errorIRMA[12])
    #
    # ax2.set_xlabel(r'VP verification period: $\it{m}$ (in epochs)', font)
    # ax2.set_ylabel('Holder bandwidth consumption (in B)', font)
    # ax2.legend(fontsize="10")
    #
    # ax2.set_title(title)
    # filename = "graphs/result_one_time_sharing_holder_bandwidth_irma"+".png"
    # fig2.savefig(filename)
    #
    #
    # for i in range(num_rows):
    #     if  resIRMAWithoutRepition[i][1] is not None:
    #         max_revocation_rate = max_revocation_rate  + 1
    #         for j in range(len(ypointsIRMAWithoutRepition[i])):
    #             print("i: ", i, "\t j: ", j, "\t bandwidth: ", ypointsIRMAWithoutRepition[i][j])
    #             ypointsIRMAWithoutRepition[i][j] = ypointsIRMAWithoutRepition[i][j]
    #             errorIRMAWithoutRepition[i][j] = errorIRMAWithoutRepition[i][j]
    #
    # fig3, ax3 = plt.subplots(layout='constrained')
    #
    #
    #
    #
    # ax3.errorbar(xpoints, ypointsZKRevoke, color="#226363",  marker='d',  label=r'zkRevoke')
    # ax3.errorbar(xpoints, ypointsIRMAWithoutRepition[0],  color="red", marker="*", linestyle=(0, (5,5)),  label=r'IRMA-no repetition: $\mathcal{R}$=1%', yerr=errorIRMAWithoutRepition[0])
    # ax3.errorbar(xpoints, ypointsIRMAWithoutRepition[4], color='#581845', marker="o", linestyle=(0, (5, 10)),   label=r'IRMA-no repetition: $\mathcal{R}$=5%', yerr=errorIRMAWithoutRepition[4])
    # ax3.errorbar(xpoints, ypointsIRMAWithoutRepition[8], color='#f3d00e', marker="s", linestyle=(0, (5, 1)), label=r'IRMA-no repetition: $\mathcal{R}$=9%', yerr=errorIRMAWithoutRepition[8])
    # ax3.errorbar(xpoints, ypointsIRMAWithoutRepition[12],  color='#0e38f3',marker="p",  linestyle = ":", label=r'IRMA-no repetition: $\mathcal{R}$=13%', yerr=errorIRMAWithoutRepition[12])
    #
    # ax3.set_xlabel(r'VP verification period: $\it{m}$ (in epochs)', font)
    # ax3.set_ylabel('Holder bandwidth consumption (in B)', font)
    # ax3.legend(fontsize="10")
    #
    # ax3.set_title(title)
    # filename = "graphs/result_one_time_sharing_holder_bandwidth_irma_without_repetition"+".png"
    # fig3.savefig(filename)

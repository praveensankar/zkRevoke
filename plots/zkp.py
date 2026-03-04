
import numpy as np
from matplotlib import pyplot as plt
from matplotlib.ticker import ScalarFormatter
from scipy.stats import sem

from results import *



def plot_zkp_proof_size_multi():
    zkrevoke_entries = parse_zkrevoke_presentation_result_entry("result_presentation_avg.json")

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
    resIMRA = {}

    error = {}
    for entry in zkrevoke_entries:
        print(entry.numberOfEpochs, entry.totalZKPProofSize, entry.numberOfTokensInCircuit)
        if entry.numberOfEpochs==1024:
            continue
        match entry.numberOfTokensInCircuit:
            case 1:
                resK1[entry.numberOfEpochs] = entry.totalZKPProofSize
                print(resK1)
            case 2:
                resK2[entry.numberOfEpochs] = entry.totalZKPProofSize

            case 4:
                resK3[entry.numberOfEpochs] = entry.totalZKPProofSize

            case 8:
                resK4[entry.numberOfEpochs] = entry.totalZKPProofSize

            case 16:
                resK5[entry.numberOfEpochs] = entry.totalZKPProofSize

            case 32:
                resK6[entry.numberOfEpochs] = entry.totalZKPProofSize

            case 64:
                resK7[entry.numberOfEpochs] = entry.totalZKPProofSize

            case 128:
                resK8[entry.numberOfEpochs] = entry.totalZKPProofSize

            case 256:
                resK9[entry.numberOfEpochs] = entry.totalZKPProofSize

            case 512:
                resK10[entry.numberOfEpochs] = entry.totalZKPProofSize





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
    y1points = np.array(list(resK1.values()))
    y2points = np.array(list(resK2.values()))
    y3points = np.array(list(resK3.values()))
    y4points = np.array(list(resK4.values()))
    y5points = np.array(list(resK5.values()))
    y6points = np.array(list(resK6.values()))
    y7points = np.array(list(resK7.values()))
    y8points = np.array(list(resK8.values()))
    y9points = np.array(list(resK9.values()))
    y10points = np.array(list(resK10.values()))
    print(y1points)
    print("y7 points: " , y7points)
    y1error = sem(y1points)
    y2error = sem(y2points)
    y3error = sem(y3points)
    y4error = sem(y4points)
    y5error = sem(y5points)
    y6error = sem(y6points)
    y7error = sem(y7points)
    y8error = sem(y8points)
    y9error = sem(y9points)
    y10error = sem(y10points)
    font = {'fontname': 'Times New Roman', 'size': 15, 'weight': 'bold'}


    fig, ax = plt.subplots(layout='constrained')

    ax.set_xscale('log')
    ax.set_yscale('log')
    x = np.arange(9)
    xlabel = [str('0'), str(r'$2^{1}$'), str(r'$2^{2}$'), str(r'$2^{3}$'), str(r'$2^{4}$'), str(r'$2^{5}$'), str(r'$2^{6}$'), str(r'$2^{7}$'), str(r'$2^{8}$'), str(r'$2^{9}$'), str(r'$2^{10}$')]
    ax.set_xticks([0, 2, 4, 8, 16, 32,64,128,256,512, 1024], xlabel)


    plt.plot(xpoints, y1points, color="#226363",  marker='d', label=r'$\it{k}=1$')
    plt.plot(xpoints, y2points, color="#581845",marker='d', label=r'$\it{k}=2^{1}$')
    plt.plot(xpoints, y3points, color="#1fd8e1", marker='v', label=r'$\it{k}=2^{2}$')
    plt.plot(xpoints, y4points,  color="#1f4be1",  marker='o', label=r'$\it{k}=2^{3}$')
    plt.plot(xpoints, y5points, color="#e1a61f",  marker='s', label=r'$\it{k}=2^{4}$')
    plt.plot(xpoints, y6points, color="#921fe1", marker='p', label=r'$\it{k}=2^{5}$')
    plt.plot(xpoints, y7points,  color="#e3c4ee", marker='h', label=r'$\it{k}=2^{6}$')



    plt.xlabel(r'verification period ($m$: number of epochs)', font)
    plt.ylabel('non-revocation proof size (in bytes)', font)
    plt.legend(fontsize="10")
    plt.savefig("graphs/result_one_time_sharing_log.png")






def plot_zkp_verification_time():
    zkrevoke_entries = parse_zkrevoke_verification_result_entry("result_presentation.json")
    irma_entries = parse_irma_presentation_and_verification_result_entry("result_presentation_verification.json")

    resK1 = {}
    errorK1 = {}
    keys = set()


    for entry in zkrevoke_entries:
        if entry.numberOfEpochs==1024:
            continue
        match entry.numberOfTokensInCircuit:
            case 1:
                if entry.__hash__() in keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK1[entry.numberOfEpochs] = np.append(resK1[entry.numberOfEpochs], totalProofGenTime)
                else:
                    resK1[entry.numberOfEpochs]= np.asarray(entry.totalZKPProofGenTime)
                    keys.add(entry.__hash__())


    for key, value in resK1.items():
        print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK1[key] = int(np.mean(value))
        errorK1[key] = np.std(value)



    resIMRA = {}
    errorIMRA = {}
    irmaKeys = set()
    for entry in irma_entries:
        if entry.numberOfEpochs==1024:
            continue
        temp = entry.witnessUpdateTime + entry.proofGenTime
        if entry.__hash__() in irmaKeys:
            resIMRA[entry.numberOfEpochs] = np.append(resIMRA[entry.numberOfEpochs], np.asarray(temp))
        else:
            resIMRA[entry.numberOfEpochs] = np.asarray(temp)
            irmaKeys.add(entry.__hash__())

    for key, value in resIMRA.items():
        print("IRMA: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resIMRA[key] = int(np.mean(value))
        errorIMRA[key] = np.std(value)




    resK1 = dict(sorted(resK1.items()))
    resIMRA = dict(sorted(resIMRA.items()))
    xpoints = np.array(list(resK1.keys()))
    y1points = np.array(list(resK1.values()))
    y1points = np.ceil(y1points/1000)
    yIRMApoints = np.array(list(resIMRA.values()))
    yIRMApoints = np.ceil(yIRMApoints / 1000)
    errorK1 = dict(sorted(errorK1.items()))
    ey1points = np.array(list(errorK1.values()))
    ey1points = np.ceil(ey1points/1000)
    errorIMRA = dict(sorted(errorIMRA.items()))
    eyIRMApoints = np.array(list(errorIMRA.values()))
    eyIRMApoints = np.ceil(eyIRMApoints/1000)

    font = {'fontname': 'Times New Roman',  'weight': 'bold'}

    fig, ax = plt.subplots(layout='constrained')
    plt.tight_layout()
    plt.figure(figsize=(3.5,3))
    ax.set_xscale('log')
    # ax.set_yscale('log')
    x = np.arange(9)
    xlabel = [str('0'), str(r'$2^{1}$'), str(r'$2^{2}$'), str(r'$2^{3}$'), str(r'$2^{4}$'), str(r'$2^{5}$'), str(r'$2^{6}$'), str(r'$2^{7}$'), str(r'$2^{8}$'), str(r'$2^{9}$'), str(r'$2^{10}$')]
    ax.set_xticks([0, 2, 4, 8, 16, 32,64,128,256,512, 1024], xlabel)

    y1labels = [str(i) for i in y1points]
    yIRMAlabels = [str(i) for i in yIRMApoints]

    for i, label in enumerate(y1labels):
        # bbox_props = dict(boxstyle='square,pad=0.2', alpha=0.5)
        # ax.text(x1points[i], y1points[i], label)
        ax.annotate(label, (xpoints[i], y1points[i]), xytext=(-5, 8), textcoords='offset points',
                    arrowprops=dict(arrowstyle='<-'))

    for i, label in enumerate(yIRMAlabels):
        # bbox_props = dict(boxstyle='square,pad=0.2', alpha=0.5)
        # ax.text(x1points[i], y1points[i], label)
        ax.annotate(label, (xpoints[i], yIRMApoints[i]), xytext=(3, -7), textcoords='offset points',
                    arrowprops=dict(arrowstyle='->'))

    plt.errorbar(xpoints, y1points, color="#226363",  marker='d', label=r'$\it{k}=1$', yerr=ey1points)
    plt.errorbar(xpoints, yIRMApoints, color="red", marker='X', label=r'IRMA', yerr=eyIRMApoints)



    plt.xlabel(r'VP verification period:  m (in epochs)', font, fontsize=14)
    plt.ylabel('time (in ms) ', font)
    plt.legend(fontsize="10")

    plt.savefig("graphs/result_zkp_proof_gen_time.png", bbox_inches='tight')




def plot_zkp_proof_gen_time(totalVCs, downsample_rate):
    zkrevoke_entries = parse_zkrevoke_presentation_result_entry("result_presentation.json")
    irma_entries = parse_irma_presentation_and_verification_result_entry("result_presentation_verification.json")

    resK1 = {}
    errorK1 = {}
    keys = set()


    for entry in zkrevoke_entries:
        if entry.VPValidity == 0:
            continue

        match entry.numberOfTokensInCircuit:
            case 1:
                if entry.__hash__() in keys:
                    totalProofGenTime = np.asarray(entry.totalZKPProofGenTime)
                    resK1[entry.VPValidity] = np.append(resK1[entry.VPValidity], totalProofGenTime)
                else:
                    resK1[entry.VPValidity]= np.asarray(entry.totalZKPProofGenTime)
                    keys.add(entry.__hash__())


    for key, value in resK1.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK1[key] = int(np.mean(value))
        errorK1[key] = np.std(value)



    resIMRA = {}
    errorIMRA = {}
    irmaKeys = set()
    for entry in irma_entries:
        if entry.totalVCs == totalVCs:
            match entry.revocationRate:
                case 1:
                    temp = entry.proofGenTime
                    if entry.__hash__() in irmaKeys:
                        resIMRA[entry.currentEpoch] = np.append(resIMRA[entry.currentEpoch], np.asarray(temp))
                    else:
                        resIMRA[entry.currentEpoch] = np.asarray(temp)
                        irmaKeys.add(entry.__hash__())

    for key, value in resIMRA.items():
        # print("IRMA: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resIMRA[key] = int(np.mean(value))
        errorIMRA[key] = np.std(value)




    resK1 = dict(sorted(resK1.items()))
    resIMRA = dict(sorted(resIMRA.items()))
    xpoints = np.array(list(resK1.keys()))
    y1points = np.array(list(resK1.values()))
    ypointsZKRevoke = np.ceil(y1points/1000)
    yIRMApoints = np.array(list(resIMRA.values()))
    yIRMApoints = np.ceil(yIRMApoints / 1000)
    errorK1 = dict(sorted(errorK1.items()))
    ey1points = np.array(list(errorK1.values()))
    ey1points = np.ceil(ey1points/1000)
    errorIMRA = dict(sorted(errorIMRA.items()))
    eyIRMApoints = np.array(list(errorIMRA.values()))
    eyIRMApoints = np.ceil(eyIRMApoints/1000)


    xpoints = xpoints[::downsample_rate]
    ypointsZKRevoke = ypointsZKRevoke[::downsample_rate]
    ey1points = ey1points[::downsample_rate]
    yIRMApoints = yIRMApoints[::downsample_rate]
    eyIRMApoints = eyIRMApoints[::downsample_rate]


    font = {'fontname': 'Times New Roman', 'weight': 'bold'}


    fig1, ax1 = plt.subplots(layout='constrained')

    for i in  range(len(yIRMApoints)):
        # if  i==0:
        #     ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsZKRevoke[i]))), (xpoints[i], ypointsZKRevoke[i]), xytext=(0, 12), textcoords='offset points',
        #                      arrowprops=dict(arrowstyle='->'), fontsize=9)
        #     ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(yIRMApoints[i]))), (xpoints[i], yIRMApoints[i]), xytext=(0, -7), textcoords='offset points',
        #              arrowprops=dict(arrowstyle='->'), fontsize=9)
        if  i==4:
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsZKRevoke[i]))), (xpoints[i], ypointsZKRevoke[i]), xytext=(-35, 20), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=9)
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(yIRMApoints[i]))), (xpoints[i], yIRMApoints[i]), xytext=(-15, 12), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=9)
        if  i==9:
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsZKRevoke[i]))), (xpoints[i], ypointsZKRevoke[i]), xytext=(-55, 15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=9)
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(yIRMApoints[i]))), (xpoints[i], yIRMApoints[i]), xytext=(-35, 14), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=9)


    title = "Total VCs: "
    if totalVCs==1000:
        title = title + "1K"
    if totalVCs == 10000:
        title = title+"10K"
    if totalVCs == 100000:
        title = title + "100K"
    if totalVCs== 1000000:
        title = title + "1M"

    # end = yIRMApoints[len(yIRMApoints)-1]
    # yticks_label = []
    #
    # yticks_locations = np.linspace(0, end, 6)
    # for i in yticks_locations:
    #     yticks_label.append(str(int(np.ceil(np.ceil(i)/10000)))+"$ * 10^4$")
    #
    #
    # ax1.set_yticks(yticks_locations, yticks_label)
    plt.tight_layout()




    ax1.errorbar(xpoints, ypointsZKRevoke, color="#226363",  marker='d',  label=r'zkRevoke', yerr=ey1points)
    ax1.errorbar(xpoints, yIRMApoints,  linestyle=(0, (5,1)), marker='x', color="red", label=r'IRMA: $\mathcal{R}$=1%', yerr=eyIRMApoints)

    ax1.set_ylabel('ZK proof gen. time (in ms)', font, fontsize=14)
    ax1.set_xlabel(r'VP verification period:  m (in epochs)', font, fontsize=14)
    ax1.legend(fontsize="10", framealpha=0.3)
    fig1.set_size_inches(3.5, 3)
    filename = "graphs/fig_4b_result_zkp_proof_gen_time"+".png"
    fig1.savefig(filename, bbox_inches='tight')






def plot_zkp_proof_size(totalVCs, downsample_rate):
    zkrevoke_entries = parse_zkrevoke_presentation_result_entry("result_presentation.json")
    irma_entries = parse_irma_presentation_and_verification_result_entry("result_presentation_verification.json")

    resK1 = {}
    errorK1 = {}
    keys = set()


    for entry in zkrevoke_entries:
        if entry.VPValidity == 0:
            continue

        match entry.numberOfTokensInCircuit:
            case 1:
                if entry.__hash__() in keys:
                    totalZKPProofSize = np.asarray(entry.totalZKPProofSize)
                    resK1[entry.VPValidity] = np.append(resK1[entry.VPValidity], totalZKPProofSize)
                else:
                    resK1[entry.VPValidity]= np.asarray(entry.totalZKPProofSize)
                    keys.add(entry.__hash__())


    for key, value in resK1.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK1[key] = np.mean(value)/1024
        errorK1[key] = np.std(value)/1024



    resIMRA = {}
    errorIMRA = {}
    irmaKeys = set()
    for entry in irma_entries:
        if entry.totalVCs == totalVCs:
            match entry.revocationRate:
                case 1:
                    temp = entry.nonRevProofSize
                    if entry.__hash__() in irmaKeys:
                        resIMRA[entry.currentEpoch] = np.append(resIMRA[entry.currentEpoch], np.asarray(temp))
                    else:
                        resIMRA[entry.currentEpoch] = np.asarray(temp)
                        irmaKeys.add(entry.__hash__())

    for key, value in resIMRA.items():
        # print("IRMA: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resIMRA[key] = np.mean(value)/1024
        errorIMRA[key] = np.std(value)/1024




    resK1 = dict(sorted(resK1.items()))
    resIMRA = dict(sorted(resIMRA.items()))
    xpoints = np.array(list(resK1.keys()))
    ypointsZKRevoke = np.array(list(resK1.values()))
    # ypointsZKRevoke = np.ceil(ypointsZKRevoke/1000)
    yIRMApoints = np.array(list(resIMRA.values()))
    # yIRMApoints = np.ceil(yIRMApoints / 1000)
    errorK1 = dict(sorted(errorK1.items()))
    ey1points = np.array(list(errorK1.values()))
    # ey1points = np.ceil(ey1points/1000)
    errorIMRA = dict(sorted(errorIMRA.items()))
    eyIRMApoints = np.array(list(errorIMRA.values()))
    # eyIRMApoints = np.ceil(eyIRMApoints/1000)


    font = {'fontname': 'Times New Roman', 'weight': 'bold'}


    xpoints = xpoints[::downsample_rate]
    ypointsZKRevoke = ypointsZKRevoke[::downsample_rate]
    ey1points = ey1points[::downsample_rate]
    yIRMApoints = yIRMApoints[::downsample_rate]
    eyIRMApoints = eyIRMApoints[::downsample_rate]

    fig1, ax1 = plt.subplots(layout='constrained')

    for i in  range(len(yIRMApoints)):
        # if i==0:
        #     ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsZKRevoke[i]))), (xpoints[i], ypointsZKRevoke[i]), xytext=(0,-5), textcoords='offset points',
        #          arrowprops=dict(arrowstyle='->'), fontsize=9)
        #     ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(yIRMApoints[i]))), (xpoints[i], yIRMApoints[i]), xytext=(-5, 10), textcoords='offset points',
        #          arrowprops=dict(arrowstyle='->'), fontsize=9)
        if  i==4 :
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsZKRevoke[i]))), (xpoints[i], ypointsZKRevoke[i]), xytext=(-15, 15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=9)
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(yIRMApoints[i]))), (xpoints[i], yIRMApoints[i]), xytext=(-80, 3), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'), fontsize=9)
        if  i==9:
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsZKRevoke[i]))), (xpoints[i], ypointsZKRevoke[i]), xytext=(-35, 15), textcoords='offset points',
                     arrowprops=dict(arrowstyle='->'), fontsize=9)
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(yIRMApoints[i]))), (xpoints[i], yIRMApoints[i]),  xytext=(-70, 3), textcoords='offset points',
                     arrowprops=dict(arrowstyle='->'), fontsize=9)


    title = "Total VCs: "
    if totalVCs==1000:
        title = title + "1K"
    if totalVCs == 10000:
        title = title+"10K"
    if totalVCs == 100000:
        title = title + "100K"
    if totalVCs== 1000000:
        title = title + "1M"
    ax1.errorbar(xpoints, ypointsZKRevoke, color="#226363",   marker='d',  label=r'zkRevoke', yerr=ey1points)
    ax1.errorbar(xpoints, yIRMApoints,  linestyle=(0, (5,10)), marker='x', color="red", label=r'IRMA: $\mathcal{R}$=1%', yerr=eyIRMApoints)


    end = yIRMApoints[len(yIRMApoints)-1]
    yticks_label = []

    # yticks_locations = np.linspace(0, end, 6)
    # for i in yticks_locations:
    #     yticks_label.append(str(int(np.ceil(np.ceil(i)/10000)))+"$ * 10^4$")
    #
    #
    # ax1.set_yticks(yticks_locations, yticks_label)
    plt.tight_layout()

    ax1.set_ylabel('ZK proof size (in KB)', font, fontsize=14)
    ax1.set_xlabel(r'VP verification period:  m (in epochs)', font, fontsize=14)

    ax1.legend(fontsize="10", framealpha=0.3)
    fig1.set_size_inches(3.5, 3)
    filename = "graphs/fig_4a_result_zkp_proof_size"+".png"
    fig1.savefig(filename, bbox_inches='tight')




def plot_zkp_proof_ver_time(totalVCs, downsample_rate):
    zkrevoke_entries = parse_zkrevoke_verification_result_entry("result_verification.json")
    irma_entries = parse_irma_presentation_and_verification_result_entry("result_presentation_verification.json")


    resK1 = {}
    errorK1 = {}
    keys = set()


    for entry in zkrevoke_entries:
        if entry.vpValidity == 0:
            continue

        match entry.numberOfTokensInCircuit:
            case 1:
                if entry.__hash__() in keys:
                    zkpProofVerTime = np.asarray(entry.zkpProofVerTime)
                    resK1[entry.vpValidity] = np.append(resK1[entry.vpValidity], zkpProofVerTime)
                else:
                    resK1[entry.vpValidity]= np.asarray(entry.zkpProofVerTime)
                    keys.add(entry.__hash__())


    for key, value in resK1.items():
        # print("zkRevoke: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resK1[key] = int(np.mean(value))
        errorK1[key] = np.std(value)



    resIMRA = {}
    errorIMRA = {}
    irmaKeys = set()
    for entry in irma_entries:
        if entry.totalVCs == totalVCs:
            match entry.revocationRate:
                case 1:
                    temp = entry.proofVerTime
                    if entry.__hash__() in irmaKeys:
                        resIMRA[entry.currentEpoch] = np.append(resIMRA[entry.currentEpoch], np.asarray(temp))
                    else:
                        resIMRA[entry.currentEpoch] = np.asarray(temp)
                        irmaKeys.add(entry.__hash__())

    for key, value in resIMRA.items():
        # print("IRMA: Key: ", key, "\t values: ", value, "\t mean: ", np.mean(value), "\t std: ", np.std(value))
        resIMRA[key] = int(np.mean(value))
        errorIMRA[key] = np.std(value)




    resK1 = dict(sorted(resK1.items()))
    resIMRA = dict(sorted(resIMRA.items()))
    xpoints = np.array(list(resK1.keys()))
    ypointsZKRevoke = np.array(list(resK1.values()))
    ypointsZKRevoke = np.ceil(ypointsZKRevoke/1000)
    yIRMApoints = np.array(list(resIMRA.values()))
    yIRMApoints = np.ceil(yIRMApoints / 1000)
    errorK1 = dict(sorted(errorK1.items()))
    ey1points = np.array(list(errorK1.values()))
    ey1points = np.ceil(ey1points/1000)
    errorIMRA = dict(sorted(errorIMRA.items()))
    eyIRMApoints = np.array(list(errorIMRA.values()))
    eyIRMApoints = np.ceil(eyIRMApoints/1000)
    xpoints = xpoints[::downsample_rate]
    ypointsZKRevoke = ypointsZKRevoke[::downsample_rate]
    ey1points = ey1points[::downsample_rate]
    yIRMApoints = yIRMApoints[::downsample_rate]
    eyIRMApoints = eyIRMApoints[::downsample_rate]

    font = {'fontname': 'Times New Roman',  'weight': 'bold'}


    fig1, ax1 = plt.subplots(layout='constrained')
    plt.tight_layout()

    for i in  range(len(yIRMApoints)):
        if  i==4 or i==0:
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsZKRevoke[i]))), (xpoints[i], ypointsZKRevoke[i]), xytext=(-12, 15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(yIRMApoints[i]))), (xpoints[i], yIRMApoints[i]), xytext=(-12, -10), textcoords='offset points',
                         arrowprops=dict(arrowstyle='<-'))
        if  i==9:
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(ypointsZKRevoke[i]))), (xpoints[i], ypointsZKRevoke[i]), xytext=(-35, 15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='->'))
            ax1.annotate((str(np.ceil(xpoints[i])), str(np.ceil(yIRMApoints[i]))), (xpoints[i], yIRMApoints[i]), xytext=(-35, -15), textcoords='offset points',
                         arrowprops=dict(arrowstyle='<-'))


    title = "Total VCs: "
    if totalVCs==1000:
        title = title + "1K"
    if totalVCs == 10000:
        title = title+"10K"
    if totalVCs == 100000:
        title = title + "100K"
    if totalVCs== 1000000:
        title = title + "1M"
    ax1.errorbar(xpoints, ypointsZKRevoke, color="#226363",   marker='d',  label=r'zkRevoke', yerr=ey1points)
    ax1.errorbar(xpoints, yIRMApoints,  linestyle=(0, (5,1)), marker='x', color="red", label=r'IRMA: $\mathcal{R}$=1%', yerr=eyIRMApoints)
    ax1.set_xlabel(r'VP verification period:  m (in epochs)', font, fontsize=14)
    ax1.set_ylabel('time (in ms)', font, fontsize=14)
    ax1.legend(fontsize="11", framealpha=0.3)
    fig1.set_size_inches(3.5, 3)
    filename = "graphs/fig_4c_result_zkp_proof_ver_time"+".png"
    fig1.savefig(filename, bbox_inches='tight')






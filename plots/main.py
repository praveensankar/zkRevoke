import json
import math

from one_time_sharing import *
from circuit import *
from list_commitment import *
from circuit_constraints import *
from zkp import *
from token_storage_cost import *
from revocation import *

def main():

    totalVCs = 1000000
    downsample_rate = 5
    plot_one_time_sharing_holder_bandwidth(totalVCs, downsample_rate)
    plot_one_time_sharing_computation(totalVCs, False, downsample_rate)

    downsample_rate = 30
    plot_issuer_bandwidth_revocation(totalVCs, 1, downsample_rate)
    plot_issuer_bandwidth_revocation(totalVCs, 5, downsample_rate)
    downsample_rate = 30
    include_refresh_in_commitment = True
    plot_issuer_computation_revocation_including_commitment(totalVCs, 1, downsample_rate, include_refresh_in_commitment)
    downsample_rate = 5
    plot_list_commitment_verification_time(totalVCs, 1,downsample_rate)
    plot_zkp_proof_size(totalVCs, downsample_rate)
    plot_zkp_proof_gen_time(totalVCs, downsample_rate)
    plot_zkp_proof_ver_time(totalVCs, downsample_rate)
    plot_circuit_size(downsample_rate)
    plot_circuit_contstraint_result()
    plot_token_storage_cost_result()


    # plot_token_storage_cost()





if __name__=="__main__":
    main()
import json
import math
import numpy as np
from one_time_sharing import *
from circuit import *
from list_commitment import *
from circuit_constraints import *
from zkp import *
from token_storage_cost import *
from revocation import *
import argparse
import os


def get_total_vcs_from_config():

    parser = argparse.ArgumentParser()
    parser.add_argument("--test", action="store_true", help="Enable test mode")
    args = parser.parse_args()

    config_file = ""
    # Use flags
    if args.test:
        config_file = "test_config.json"
    else:
        config_file = "config.json"


    path = os.path.realpath(__file__)
    current_dir = os.path.dirname(path)
    dir = os.path.dirname(current_dir)
    file_path = os.path.join(dir, config_file)

    with open(file_path) as file:
        data = json.load(file)
        total_vcs = data["params"]["total_vcs"]
        return int(total_vcs)




def main():

    totalVCs = get_total_vcs_from_config()
    print(totalVCs)

    downsample_rate = math.ceil(0.000005 * totalVCs)

    plot_one_time_sharing_holder_bandwidth(totalVCs, downsample_rate)
    plot_one_time_sharing_computation(totalVCs, False, downsample_rate)


    downsample_rate = math.ceil(0.00003 * totalVCs)
    plot_issuer_bandwidth_revocation(totalVCs, 1, downsample_rate)
    plot_issuer_bandwidth_revocation(totalVCs, 5, downsample_rate)

    downsample_rate = math.ceil(0.00003 * totalVCs)
    include_refresh_in_commitment = True
    plot_issuer_computation_revocation_including_commitment(totalVCs, 1, downsample_rate, include_refresh_in_commitment)
    downsample_rate = math.ceil(0.000005 * totalVCs)
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
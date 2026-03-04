//SPDX-License-Identifier: MIT

pragma solidity ^0.8.7;
//import "hardhat/console.sol";

contract RevocationList{

    struct List{
        uint256 epoch;
        bytes32[] tokens;
    }

    // stores epoch number and the revoked tokens in that epoch
    List public list;


    // entities is the owner of the contract
    address issuer;

    // stores the hash of ccs used in the groth16 zkp scheme
    bytes public ccsHash;

    // stores the ccs used in the groth16 zkp scheme
    bytes public ccs;


    // stores the hash of ccs used in the groth16 zkp scheme
    bytes public bbsPublicKey;

    // groth16 verifying key
    bytes public ZKPVerifyingKey;

    // eddsa public key used inside the zkp circuit
    bytes public EDDSAPublicKey;

    bytes public InitialTimeStamp;
    uint public EpochDuration;
    // sets the entities - contract creator is the entities
    constructor(){
        issuer = msg.sender;
    }

    /*
    This function is used to register new issuers.
    Register did of issuers and public keys. (maybe in the form of DID Docs).
    input: did doc
    */
    function registerIssuers() public{

    }

    /*
    This function stores the hash of ccs used in  groth16 zkp scheme

    input: public keys
    */
    function publishCCSHash(bytes memory _ccsHash) public{
        //only entities can perform the revocation
        require(msg.sender==issuer);

        for (uint i = 0; i < _ccsHash.length; i++) {
            ccsHash.push(_ccsHash[i]);
        }
    }

    function RetrieveCCSHash() public view returns (bytes memory){
        return ccsHash;
    }

    /*
    This function stores the ccs used in  groth16 zkp scheme

    input: public keys
    */
    function publishCCS(bytes memory _ccs) public{
        //only entities can perform the revocation
        require(msg.sender==issuer);
        for (uint i=0;i<_ccs.length;i++){
            ccs.push(_ccs[i]);
        }
    }

    function RetrieveCCS() public view returns (bytes memory){
        return ccs;
    }


    /*
    This function stores the bbs public key used by the issuer
    input: bbs public key
    */
    function PublishBBSPublicKey(bytes memory _bbsPublicKey) public{
        //only entities can perform the revocation
        require(msg.sender==issuer);
        bbsPublicKey = _bbsPublicKey;
    }
    /*
    This function retrieves the bbs public key used by the issuer
    */
    function RetrieveBBSPublicKey() public view returns (bytes memory){
        return bbsPublicKey;
    }

    /*
    This function stores the zkp verifying key
    input: zkp verifying key (size 460 bytes)
    */
    function PublishZKPVerifyingKey(bytes memory _verifyingKey) public{
        //only entities can perform the revocation
        require(msg.sender==issuer);
        ZKPVerifyingKey = _verifyingKey;
    }
    /*
    This function retrieves the zkp verifying key
    */
    function RetrieveZKPVerifyingKey() public view returns (bytes memory){
        return ZKPVerifyingKey;
    }

    /*
This function stores the eddsa public key
input: public key (size 32 bytes)
*/
    function PublishEDDSAPublicKey(bytes memory _publicKey) public{
        //only entities can perform the revocation
        require(msg.sender==issuer);
        EDDSAPublicKey = _publicKey;
    }
    /*
    This function retrieves the eddsa public key
    */
    function RetrieveEDDSAPublicKey() public view returns (bytes memory){
        return EDDSAPublicKey;
    }


    /*
    This function stores the zkp verifying key
    input: zkp verifying key (size 460 bytes)
    */
    function PublishEpochConfigurations(uint epoch_duration, bytes memory intial_timestamp) public{
        //only entities can perform the revocation
        require(msg.sender==issuer);
        EpochDuration = epoch_duration;
        InitialTimeStamp = intial_timestamp;
    }

    /*
    This function retrieves the epoch duration
    */
    function RetrieveEpochDuration() public view returns (uint){
        return EpochDuration;
    }

    /*
    This function retrieves the epoch duration
    */
    function RetrieveInitialTimeStamp() public view returns (bytes memory){
        return InitialTimeStamp;
    }
    /*
    This function is used to set proofs at merkle tree accumulator when one or more VC is issued.
    The merkle tree stores hash of VCs in leaves. Arrays using Level order structure is used to store the merkle tree.
    Every time new VC is issued, update the array.


    Note: The logic for mapping VCs to level order indexes should be done at the issuers side.
    */
    function RefreshRevokedTokens(uint256 _epoch, bytes32[] memory _tokens) public{
        //only entities can perform the revocation
        require(msg.sender==issuer);

        if (list.epoch==_epoch){
            for (uint i = 0; i < _tokens.length; i++) {
                list.tokens.push(_tokens[i]);
            }
        } else{
            delete list;
            list = List(_epoch, _tokens);
        }
    }




    /*
    Returns the latest set of revoked tokens.
    */
    function GetTokens() public view returns(List memory){
        return list;
    }

    function test() public{

        bytes32[] memory tokens = new bytes32[](2);
        tokens[0] = 0x68656c6c6f000000000000000000000000000000000000000000000000000000;
        tokens[1] = 0x68656c6c6f000000000000000000000000000000000000000000000000000000;
        uint epoch = 1;
        RefreshRevokedTokens(epoch, tokens);
//        List memory revokedTokens = GetTokens();
//        console.log("epoch: ", revokedTokens.epoch);
//        console.log("tokens: ");
//        for(uint i=0; i<revokedTokens.tokens.length; i++){
//            console.logBytes32( revokedTokens.tokens[i]);
//        }


    }
}




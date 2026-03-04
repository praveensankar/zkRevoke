package irma

import (
	"github.com/privacybydesign/gabi"
	"github.com/privacybydesign/gabi/big"
	"github.com/privacybydesign/gabi/gabikeys"
	"zkrevoke/irma/internal/common"
)

var testAttributes1 = []*big.Int{
	new(big.Int).SetBytes([]byte("one")),
	new(big.Int).SetBytes([]byte("two")),
	new(big.Int).SetBytes([]byte("three")),
	new(big.Int).SetBytes([]byte("four"))}

func createKeyshareCredential(context, secret, keyshareP *big.Int, attrs []*big.Int, issuer *gabi.Issuer) *gabi.Credential {
	// First create a credential
	keylength := 1024
	nonce1, err := common.RandomBigInt(gabikeys.DefaultSystemParameters[keylength].Lstatzk)
	if err != nil {

	}
	nonce2, err := common.RandomBigInt(gabikeys.DefaultSystemParameters[keylength].Lstatzk)

	cb, err := gabi.NewCredentialBuilder(issuer.Pk, context, secret, nonce2, keyshareP, nil)

	commitMsg, err := cb.CommitToSecretAndProve(nonce1)

	ism, err := issuer.IssueSignature(commitMsg.U, attrs, nil, nonce2, nil)

	cred, err := cb.ConstructCredential(ism, attrs)
	return cred
}

func createCredential(context, secret *big.Int, issuer *gabi.Issuer) *gabi.Credential {
	return createKeyshareCredential(context, secret, nil, testAttributes1, issuer)
}

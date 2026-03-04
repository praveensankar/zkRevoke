package irma

import (
	"github.com/privacybydesign/gabi/big"
	"github.com/privacybydesign/gabi/revocation"
)

func revocationAttrs(w *revocation.Witness) []*big.Int {
	return append(testAttributes1, w.E)
}

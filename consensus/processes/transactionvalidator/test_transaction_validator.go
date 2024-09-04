package transactionvalidator

import (
	"github.com/zilong-dai/karlsen-miner/consensus/model"
	"github.com/zilong-dai/karlsen-miner/consensus/model/testapi"
	"github.com/zilong-dai/karlsen-miner/consensus/utils/txscript"
)

type testTransactionValidator struct {
	*transactionValidator
}

// NewTestTransactionValidator creates an instance of a TestTransactionValidator
func NewTestTransactionValidator(baseTransactionValidator model.TransactionValidator) testapi.TestTransactionValidator {
	return &testTransactionValidator{transactionValidator: baseTransactionValidator.(*transactionValidator)}
}

func (tbv *testTransactionValidator) SigCache() *txscript.SigCache {
	return tbv.sigCache
}

func (tbv *testTransactionValidator) SetSigCache(sigCache *txscript.SigCache) {
	tbv.sigCache = sigCache
}

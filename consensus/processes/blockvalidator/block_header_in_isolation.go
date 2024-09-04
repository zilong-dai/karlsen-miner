package blockvalidator

import (
	"github.com/karlsen-network/karlsend/v2/infrastructure/logger"
	"github.com/karlsen-network/karlsend/v2/util/mstime"
	"github.com/pkg/errors"
	"github.com/zilong-dai/karlsen-miner/consensus/model"
	"github.com/zilong-dai/karlsen-miner/consensus/model/externalapi"
	"github.com/zilong-dai/karlsen-miner/consensus/ruleerrors"
	"github.com/zilong-dai/karlsen-miner/consensus/utils/consensushashing"
	"github.com/zilong-dai/karlsen-miner/consensus/utils/constants"
)

// ValidateHeaderInIsolation validates block headers in isolation from the current
// consensus state
func (v *blockValidator) ValidateHeaderInIsolation(stagingArea *model.StagingArea, blockHash *externalapi.DomainHash) error {
	onEnd := logger.LogAndMeasureExecutionTime(log, "ValidateHeaderInIsolation")
	defer onEnd()

	header, err := v.blockHeaderStore.BlockHeader(v.databaseContext, stagingArea, blockHash)
	if err != nil {
		return err
	}

	//todo : drop this
	//log.Info("blockHash %s - genesisHash %s", blockHash, v.genesisHash)

	if !blockHash.Equal(v.genesisHash) {
		err = v.checkBlockVersion(header)
		if err != nil {
			return err
		}
	}

	err = v.checkBlockTimestampInIsolation(header)
	if err != nil {
		return err
	}

	err = v.checkParentsLimit(header)
	if err != nil {
		return err
	}

	return nil
}

func (v *blockValidator) checkParentsLimit(header externalapi.BlockHeader) error {
	hash := consensushashing.HeaderHash(header)
	if len(header.DirectParents()) == 0 && !hash.Equal(v.genesisHash) {
		return errors.Wrapf(ruleerrors.ErrNoParents, "block has no parents")
	}

	if uint64(len(header.DirectParents())) > uint64(v.maxBlockParents) {
		return errors.Wrapf(ruleerrors.ErrTooManyParents, "block header has %d parents, but the maximum allowed amount "+
			"is %d", len(header.DirectParents()), v.maxBlockParents)
	}
	return nil
}

func (v *blockValidator) checkBlockVersion(header externalapi.BlockHeader) error {
	/*
		if header.Version() != constants.BlockVersion {
			return errors.Wrapf(
				ruleerrors.ErrWrongBlockVersion, "The block version should be %d", constants.BlockVersion)
		}
	*/
	if header.DAAScore() >= v.hfDAAScore && header.Version() != constants.BlockVersionKHashV2 {
		log.Warnf("After HF1 the block version should be %d - block[%d][v%d]", constants.BlockVersionKHashV2, header.DAAScore(), header.Version())
		return errors.Wrapf(ruleerrors.ErrWrongBlockVersion, "The block version should be %d", constants.BlockVersionKHashV2)
	} else if header.DAAScore() < v.hfDAAScore && header.Version() != constants.BlockVersionKHashV1 {
		log.Warnf("Before HF1 the block version should be %d - block[%d][v%d]", constants.BlockVersionKHashV1, header.DAAScore(), header.Version())
		return errors.Wrapf(ruleerrors.ErrWrongBlockVersion, "The block version should be %d", constants.BlockVersionKHashV1)
	}
	return nil
}

func (v *blockValidator) checkBlockTimestampInIsolation(header externalapi.BlockHeader) error {
	blockTimestamp := header.TimeInMilliseconds()
	now := mstime.Now().UnixMilliseconds()
	maxCurrentTime := now + int64(v.timestampDeviationTolerance)*v.targetTimePerBlock.Milliseconds()
	if blockTimestamp > maxCurrentTime {
		return errors.Wrapf(
			ruleerrors.ErrTimeTooMuchInTheFuture, "The block timestamp is in the future.")
	}
	return nil
}

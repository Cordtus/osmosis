package sqs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v20/ingest"
	"github.com/osmosis-labs/osmosis/v20/ingest/sqs/domain"
)

var _ ingest.Ingester = &sqsIngester{}

// sqsIngester is a sidecar query server (SQS) implementation of Ingester.
// It encapsulates all individual SQS ingesters.
type sqsIngester struct {
	txManager     domain.TxManager
	poolsIngester ingest.AtomicIngester
}

// NewSidecarQueryServerIngester creates a new sidecar query server ingester.
// poolsRepository is the storage for pools.
// gammKeeper is the keeper for Gamm pools.
func NewSidecarQueryServerIngester(poolsIngester ingest.AtomicIngester, txManager domain.TxManager) ingest.Ingester {
	return &sqsIngester{
		txManager:     txManager,
		poolsIngester: poolsIngester,
	}
}

// ProcessBlock implements ingest.Ingester.
func (i *sqsIngester) ProcessBlock(ctx sdk.Context) error {
	// Start atomic transaction
	tx := i.txManager.StartTx()

	goCtx := sdk.WrapSDKContext(ctx)

	// Begin by flushing all previous writes
	if err := tx.ClearAll(goCtx); err != nil {
		return err
	}

	// Process block by reading and writing data and ingesting data into sinks
	if err := i.poolsIngester.ProcessBlock(ctx, tx); err != nil {
		return err
	}

	// Flush all writes atomically
	return tx.Exec(goCtx)
}
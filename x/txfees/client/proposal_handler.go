package client

import (
	"github.com/osmosis-labs/osmosis/vv21/x/txfees/client/cli"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	SubmitUpdateFeeTokenProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpdateFeeTokenProposal)
)

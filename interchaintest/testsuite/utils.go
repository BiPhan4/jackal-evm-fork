package testsuite

import (
	"context"

	math "cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
)

const (
	defaultRlyHomeDir = "/home/relayer"
)

// FundAddressChainB sends funds to the given address on chain B.
// The amount sent is 1,000,000,000 of the chain's denom.
func (s *TestSuite) FundAddressChainB(ctx context.Context, address string) {
	s.fundAddress(ctx, s.ChainB, s.UserB.KeyName(), address)
}

// fundAddress sends funds to the given address on the given chain
func (s *TestSuite) fundAddress(ctx context.Context, chain *cosmos.CosmosChain, keyName, address string) {
	err := chain.SendFunds(ctx, keyName, ibc.WalletAmount{
		Address: address,
		Denom:   chain.Config().Denom,
		Amount:  math.NewInt(10_000_000_000),
		// NOTE: not sure why we have to wrap it
		// while srdrk did not, IDE does not
		// complain in his repo, but complains in ours?
		/// 1_000_000_000_000
		///
		//

	})
	s.Require().NoError(err)

	// wait for 2 blocks for the funds to be received
	err = testutil.WaitForBlocks(ctx, 2, chain)
	s.Require().NoError(err)
}

// createChannelWithOverride creates a channel between the two chains with the given port and channel identifiers
// with the override flag so that the channel is created even if it already exists.
func (s *TestSuite) CreateChannelWithOverride(ctx context.Context, srcPort, dstPort string, order ibc.Order, version string) {
	cmd := []string{
		"rly", "tx", "channel", s.PathName, "--src-port", srcPort, "--dst-port", dstPort,
		"--order", order.String(), "--version", version, "--override", "--home", defaultRlyHomeDir,
	}
	res := s.Relayer.Exec(ctx, s.ExecRep, cmd, nil)
	s.Require().NoError(res.Err)
	s.Require().Equal(0, res.ExitCode)
}

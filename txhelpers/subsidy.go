// Copyright (c) 2015-2018 The Decred developers
// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txhelpers

import (
	"github.com/picfight/pfcd/blockchain"
	"github.com/picfight/pfcd/chaincfg"
)

// UltimateSubsidy computes the total subsidy over the entire subsidy
// distribution period of the network.
func UltimateSubsidy(net *chaincfg.Params) int64 {
	if net.SubsidyCalculator != nil {
		return net.SubsidyCalculator().ExpectedTotalNetworkSubsidy().AtomsValue
	}
	params := net.DecredSubsidyParams
	subsidyCache := blockchain.NewSubsidyCache(0, net)

	totalSubsidy := net.BlockOneSubsidy()
	for i := int64(0); ; i++ {
		// Genesis block or first block.
		if i <= 1 {
			continue
		}

		if i%params.SubsidyReductionInterval == 0 {
			numBlocks := params.SubsidyReductionInterval
			// First reduction internal, which is reduction interval - 2 to skip
			// the genesis block and block one.
			if i == params.SubsidyReductionInterval {
				numBlocks -= 2
			}
			height := i - numBlocks

			work := blockchain.CalcBlockWorkSubsidy(subsidyCache, height,
				net.TicketsPerBlock, net)
			stake := blockchain.CalcStakeVoteSubsidy(subsidyCache, height,
				net) * int64(net.TicketsPerBlock)
			tax := blockchain.CalcBlockTaxSubsidy(subsidyCache, height,
				net.TicketsPerBlock, net)
			if (work + stake + tax) == 0 {
				break // all done
			}
			totalSubsidy += ((work + stake + tax) * numBlocks)

			// First reduction internal -- subtract the stake subsidy for blocks
			// before the staking system is enabled.
			if i == params.SubsidyReductionInterval {
				totalSubsidy -= stake * (net.StakeValidationHeight - 2)
			}
		}
	}
	return totalSubsidy
}

// RewardsAtBlock computes the PoW, PoS (per vote), and project fund subsidies
// at for the specified block index, assuming a certain number of votes.
func RewardsAtBlock(blockIdx int64, votes uint16, p *chaincfg.Params) (work, stake, tax int64) {
	subsidyCache := blockchain.NewSubsidyCache(0, p)
	work = blockchain.CalcBlockWorkSubsidy(subsidyCache, blockIdx, votes, p)
	stake = blockchain.CalcStakeVoteSubsidy(subsidyCache, blockIdx, p)
	tax = blockchain.CalcBlockTaxSubsidy(subsidyCache, blockIdx, votes, p)
	return
}

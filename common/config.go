package common

import (
	"math"
	"math/big"

	intf "github.com/arcology-network/evm-adaptor/interface"
	"github.com/ethereum/go-ethereum/common"
	evmcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

// DummyChain implements the ChainContext interface.
type DummyChain struct{}

func (chain *DummyChain) GetHeader(evmcommon.Hash, uint64) *types.Header { return &types.Header{} }
func (chain *DummyChain) Engine() consensus.Engine                       { return nil }

// Config contains all the static settings used in Schedule.
type Config struct {
	ChainConfig *params.ChainConfig
	VMConfig    *vm.Config
	BlockNumber *big.Int    // types.Header.Number
	ParentHash  common.Hash // types.Header.ParentHash
	Time        *big.Int    // types.Header.Time
	Chain       intf.ChainContext
	Coinbase    *evmcommon.Address
	GasLimit    uint64   // types.Header.GasLimit
	Difficulty  *big.Int // types.Header.Difficulty
}

func NewEmptyConfig() *Config {
	cfg := &Config{
		ChainConfig: params.MainnetChainConfig,
		VMConfig:    &vm.Config{},
		BlockNumber: big.NewInt(0),
		ParentHash:  evmcommon.Hash{},
		Time:        big.NewInt(0),
		Coinbase:    &evmcommon.Address{},
		GasLimit:    math.MaxUint64,
		Difficulty:  big.NewInt(0),
	}
	cfg.Chain = new(DummyChain)
	return cfg
}

func NewConfigFromBlockContext(context vm.BlockContext) *Config {
	cfg := &Config{
		ChainConfig: params.MainnetChainConfig,
		VMConfig:    &vm.Config{},
		BlockNumber: context.BlockNumber,
		ParentHash:  evmcommon.Hash{},
		Time:        big.NewInt(int64(context.Time)),
		Coinbase:    &context.Coinbase,
		GasLimit:    context.GasLimit,
		Difficulty:  context.Difficulty,
	}
	cfg.Chain = new(DummyChain)
	return cfg
}
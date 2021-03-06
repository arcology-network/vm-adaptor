package evm

import (
	"math/big"

	"github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
)

type StateDB interface {
	CreateAccount(common.Address)

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int
	GetBalanceNoRecord(common.Address) *big.Int
	SetBalance(common.Address, *big.Int)

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int

	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64

	GetCommittedState(common.Address, common.Hash) common.Hash
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	Suicide(common.Address) bool
	HasSuicided(common.Address) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(common.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(common.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(*types.Log)
	AddPreimage(common.Hash, []byte)

	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) error

	PrepareAccessList(sender common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList)
	AddressInAccessList(addr common.Address) bool
	SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool)
	// AddAddressToAccessList adds the given address to the access list. This operation is safe to perform
	// even if the feature/fork is not active yet
	AddAddressToAccessList(addr common.Address)
	// AddSlotToAccessList adds the given (address,slot) to the access list. This operation is safe to perform
	// even if the feature/fork is not active yet
	AddSlotToAccessList(addr common.Address, slot common.Hash)

	Prepare(thash, bhash common.Hash, ti int)
	GetLogs(hash common.Hash) []*types.Log
	Copy() StateDB

	Set(eac EthAccountCache, esc EthStorageCache)
}

type Account interface {
	GetBalance() *big.Int
	GetNonce() uint64
	GetCodeHash() []byte
}

type EthAccountCache interface {
	GetAccount(string) (Account, error)
	GetCode(string) ([]byte, error)
}

type EthStorageCache interface {
	GetState(string, []byte) []byte
}

type ILog interface {
	GetKey() string
	GetValue() string
}

type KernelAPI interface {
	AddLog(key, value string)
	GetLogs() []ILog
	ClearLogs()
	IsKernelAPI(addr common.Address) bool
	Prepare(thash common.Hash)
	Call(caller, callee common.Address, input []byte, origin common.Address, nonce uint64, blockhash common.Hash) ([]byte, bool)
}

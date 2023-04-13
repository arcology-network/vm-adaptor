package tests

import (
	"math/big"
	"testing"
	"time"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	curstorage "github.com/arcology-network/concurrenturl/v2/storage"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	adaptor "github.com/arcology-network/vm-adaptor/evm"
	"github.com/arcology-network/vm-adaptor/tests"
)

var (
	dstokenCode = "6080604052601260045560006005553480156200001b57600080fd5b50604051602080620033dc833981018060405260208110156200003d57600080fd5b810190808051906020019092919050505033600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503373ffffffffffffffffffffffffffffffffffffffff167fce241d7ca1f669fee44b6fc00b8eba2df3bb514eed0f6f668f8f89096e81ed9460405160405180910390a280600381905550608173ffffffffffffffffffffffffffffffffffffffff1663f02e3aff600160038111156200010457fe5b600260038111156200011257fe5b6040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018460030b60030b81526020018360030b60030b8152602001828103825260098152602001807f62616c616e63654f6600000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b158015620001b157600080fd5b505af1158015620001c6573d6000803e3d6000fd5b50505050608173ffffffffffffffffffffffffffffffffffffffff1663f02e3aff60026003811115620001f557fe5b600260038111156200020357fe5b6040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018460030b60030b81526020018360030b60030b8152602001828103825260098152602001807f616c6c6f77616e636500000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b158015620002a257600080fd5b505af1158015620002b7573d6000803e3d6000fd5b50505050608273ffffffffffffffffffffffffffffffffffffffff166346f81a8760026003811115620002e657fe5b6040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018381526020018281038252600e8152602001807f746f74616c537570706c7941646400000000000000000000000000000000000081525060200192505050600060405180830381600087803b1580156200037257600080fd5b505af115801562000387573d6000803e3d6000fd5b50505050608273ffffffffffffffffffffffffffffffffffffffff166346f81a8760026003811115620003b657fe5b6040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018381526020018281038252600e8152602001807f746f74616c537570706c7953756200000000000000000000000000000000000081525060200192505050600060405180830381600087803b1580156200044257600080fd5b505af115801562000457573d6000803e3d6000fd5b5050505060a173ffffffffffffffffffffffffffffffffffffffff16632dc796886040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808060200180602001838103835260118152602001807f757064617465546f74616c537570706c79000000000000000000000000000000815250602001838103825260198152602001807f757064617465546f74616c537570706c7928737472696e67290000000000000081525060200192505050600060405180830381600087803b1580156200053657600080fd5b505af11580156200054b573d6000803e3d6000fd5b5050505050612e7c80620005606000396000f3fe60806040526004361061013e576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde031461014357806307da68f51461016e578063095ea7b31461018557806313af4035146101f857806318160ddd1461024957806323b872dd14610274578063313ce5671461030757806340c10f191461033257806342966c681461038d5780635ac801fe146103c857806375f12b21146104035780637a9e5e4b146104325780638da5cb5b1461048357806395d89b41146104da5780639dc29fac14610505578063a0712d6814610560578063a9059cbb1461059b578063b753a98c1461060e578063bb35783b14610669578063be9a6555146106e4578063bf7e214f146106fb578063daea85c514610752578063e6a7f3b3146107bb578063f2d5d56b14610883575b600080fd5b34801561014f57600080fd5b506101586108de565b6040518082815260200191505060405180910390f35b34801561017a57600080fd5b506101836108e4565b005b34801561019157600080fd5b506101de600480360360408110156101a857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506109ce565b604051808215151515815260200191505060405180910390f35b34801561020457600080fd5b506102476004803603602081101561021b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610b96565b005b34801561025557600080fd5b5061025e610ce1565b6040518082815260200191505060405180910390f35b34801561028057600080fd5b506102ed6004803603606081101561029757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610ce7565b604051808215151515815260200191505060405180910390f35b34801561031357600080fd5b5061031c611512565b6040518082815260200191505060405180910390f35b34801561033e57600080fd5b5061038b6004803603604081101561035557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611518565b005b34801561039957600080fd5b506103c6600480360360208110156103b057600080fd5b8101908080359060200190929190505050611a16565b005b3480156103d457600080fd5b50610401600480360360208110156103eb57600080fd5b8101908080359060200190929190505050611a23565b005b34801561040f57600080fd5b50610418611acf565b604051808215151515815260200191505060405180910390f35b34801561043e57600080fd5b506104816004803603602081101561045557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611ae2565b005b34801561048f57600080fd5b50610498611c2b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156104e657600080fd5b506104ef611c51565b6040518082815260200191505060405180910390f35b34801561051157600080fd5b5061055e6004803603604081101561052857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611c57565b005b34801561056c57600080fd5b506105996004803603602081101561058357600080fd5b8101908080359060200190929190505050612473565b005b3480156105a757600080fd5b506105f4600480360360408110156105be57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050612480565b604051808215151515815260200191505060405180910390f35b34801561061a57600080fd5b506106676004803603604081101561063157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050612495565b005b34801561067557600080fd5b506106e26004803603606081101561068c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506124a5565b005b3480156106f057600080fd5b506106f96124b6565b005b34801561070757600080fd5b506107106125a1565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561075e57600080fd5b506107a16004803603602081101561077557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506125c6565b604051808215151515815260200191505060405180910390f35b3480156107c757600080fd5b50610881600480360360208110156107de57600080fd5b81019080803590602001906401000000008111156107fb57600080fd5b82018360208201111561080d57600080fd5b8035906020019184600183028401116401000000008311171561082f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506125f9565b005b34801561088f57600080fd5b506108dc600480360360408110156108a657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050612a0e565b005b60055481565b610912336000357fffffffff0000000000000000000000000000000000000000000000000000000016612a1e565b1515610986576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f64732d617574682d756e617574686f72697a656400000000000000000000000081525060200191505060405180910390fd5b60018060146101000a81548160ff0219169083151502179055507fbedf0f4abfe86d4ffad593d9607fe70e83ea706033d44d24b3b6283cf3fc4f6b60405160405180910390a1565b6000600160149054906101000a900460ff16151515610a55576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f64732d73746f702d69732d73746f70706564000000000000000000000000000081525060200191505060405180910390fd5b608173ffffffffffffffffffffffffffffffffffffffff166336f3c77d610a7c3386612c93565b846040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018080602001848152602001838152602001828103825260098152602001807f616c6c6f77616e636500000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b158015610b0f57600080fd5b505af1158015610b23573d6000803e3d6000fd5b505050508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a36001905092915050565b610bc4336000357fffffffff0000000000000000000000000000000000000000000000000000000016612a1e565b1515610c38576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f64732d617574682d756e617574686f72697a656400000000000000000000000081525060200191505060405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fce241d7ca1f669fee44b6fc00b8eba2df3bb514eed0f6f668f8f89096e81ed9460405160405180910390a250565b60025481565b6000600160149054906101000a900460ff16151515610d6e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f64732d73746f702d69732d73746f70706564000000000000000000000000000081525060200191505060405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1614151561100e576000610daf8533612c93565b90506000608173ffffffffffffffffffffffffffffffffffffffff16638d206aad836040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018080602001838152602001828103825260098152602001807f616c6c6f77616e636500000000000000000000000000000000000000000000008152506020019250505060206040518083038186803b158015610e5a57600080fd5b505afa158015610e6e573d6000803e3d6000fd5b505050506040513d6020811015610e8457600080fd5b810190808051906020019092919050505090507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114151561100b57838110151515610f38576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f64732d746f6b656e2d696e73756666696369656e742d617070726f76616c000081525060200191505060405180910390fd5b608173ffffffffffffffffffffffffffffffffffffffff166336f3c77d83610f608488612d46565b6040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018080602001848152602001838152602001828103825260098152602001807f616c6c6f77616e636500000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b158015610ff257600080fd5b505af1158015611006573d6000803e3d6000fd5b505050505b50505b6000608173ffffffffffffffffffffffffffffffffffffffff1663c41eb85a866040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828103825260098152602001807f62616c616e63654f6600000000000000000000000000000000000000000000008152506020019250505060206040518083038186803b1580156110e357600080fd5b505afa1580156110f7573d6000803e3d6000fd5b505050506040513d602081101561110d57600080fd5b81019080805190602001909291905050509050828110151515611198576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f64732d746f6b656e2d696e73756666696369656e742d62616c616e636500000081525060200191505060405180910390fd5b608173ffffffffffffffffffffffffffffffffffffffff16634f7c4f4c866111c08487612d46565b6040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828103825260098152602001807f62616c616e63654f6600000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b15801561127e57600080fd5b505af1158015611292573d6000803e3d6000fd5b50505050608173ffffffffffffffffffffffffffffffffffffffff16634f7c4f4c856113cb608173ffffffffffffffffffffffffffffffffffffffff1663c41eb85a896040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828103825260098152602001807f62616c616e63654f6600000000000000000000000000000000000000000000008152506020019250505060206040518083038186803b15801561138a57600080fd5b505afa15801561139e573d6000803e3d6000fd5b505050506040513d60208110156113b457600080fd5b810190808051906020019092919050505087612dcb565b6040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828103825260098152602001807f62616c616e63654f6600000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b15801561148957600080fd5b505af115801561149d573d6000803e3d6000fd5b505050508373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040518082815260200191505060405180910390a360019150509392505050565b60045481565b611546336000357fffffffff0000000000000000000000000000000000000000000000000000000016612a1e565b15156115ba576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f64732d617574682d756e617574686f72697a656400000000000000000000000081525060200191505060405180910390fd5b600160149054906101000a900460ff1615151561163f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f64732d73746f702d69732d73746f70706564000000000000000000000000000081525060200191505060405180910390fd5b608173ffffffffffffffffffffffffffffffffffffffff16634f7c4f4c83611774608173ffffffffffffffffffffffffffffffffffffffff1663c41eb85a876040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828103825260098152602001807f62616c616e63654f6600000000000000000000000000000000000000000000008152506020019250505060206040518083038186803b15801561173357600080fd5b505afa158015611747573d6000803e3d6000fd5b505050506040513d602081101561175d57600080fd5b810190808051906020019092919050505085612dcb565b6040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828103825260098152602001807f62616c616e63654f6600000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b15801561183257600080fd5b505af1158015611846573d6000803e3d6000fd5b50505050608273ffffffffffffffffffffffffffffffffffffffff1663a0aa9f29826040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018381526020018281038252600e8152602001807f746f74616c537570706c7941646400000000000000000000000000000000000081525060200192505050600060405180830381600087803b1580156118f357600080fd5b505af1158015611907573d6000803e3d6000fd5b5050505060a173ffffffffffffffffffffffffffffffffffffffff166306e354dd6040518163ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018080602001828103825260118152602001807f757064617465546f74616c537570706c79000000000000000000000000000000815250602001915050600060405180830381600087803b1580156119ac57600080fd5b505af11580156119c0573d6000803e3d6000fd5b505050508173ffffffffffffffffffffffffffffffffffffffff167f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885826040518082815260200191505060405180910390a25050565b611a203382611c57565b50565b611a51336000357fffffffff0000000000000000000000000000000000000000000000000000000016612a1e565b1515611ac5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f64732d617574682d756e617574686f72697a656400000000000000000000000081525060200191505060405180910390fd5b8060058190555050565b600160149054906101000a900460ff1681565b611b10336000357fffffffff0000000000000000000000000000000000000000000000000000000016612a1e565b1515611b84576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f64732d617574682d756e617574686f72697a656400000000000000000000000081525060200191505060405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f1abebea81bfa2637f28358c371278fb15ede7ea8dd28d2e03b112ff6d936ada460405160405180910390a250565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60035481565b611c85336000357fffffffff0000000000000000000000000000000000000000000000000000000016612a1e565b1515611cf9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f64732d617574682d756e617574686f72697a656400000000000000000000000081525060200191505060405180910390fd5b600160149054906101000a900460ff16151515611d7e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f64732d73746f702d69732d73746f70706564000000000000000000000000000081525060200191505060405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614151561201e576000611dbf8333612c93565b90506000608173ffffffffffffffffffffffffffffffffffffffff16638d206aad836040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018080602001838152602001828103825260098152602001807f616c6c6f77616e636500000000000000000000000000000000000000000000008152506020019250505060206040518083038186803b158015611e6a57600080fd5b505afa158015611e7e573d6000803e3d6000fd5b505050506040513d6020811015611e9457600080fd5b810190808051906020019092919050505090507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114151561201b57828110151515611f48576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f64732d746f6b656e2d696e73756666696369656e742d617070726f76616c000081525060200191505060405180910390fd5b608173ffffffffffffffffffffffffffffffffffffffff166336f3c77d83611f708487612d46565b6040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018080602001848152602001838152602001828103825260098152602001807f616c6c6f77616e636500000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b15801561200257600080fd5b505af1158015612016573d6000803e3d6000fd5b505050505b50505b6000608173ffffffffffffffffffffffffffffffffffffffff1663c41eb85a846040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828103825260098152602001807f62616c616e63654f6600000000000000000000000000000000000000000000008152506020019250505060206040518083038186803b1580156120f357600080fd5b505afa158015612107573d6000803e3d6000fd5b505050506040513d602081101561211d57600080fd5b810190808051906020019092919050505090508181101515156121a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f64732d746f6b656e2d696e73756666696369656e742d62616c616e636500000081525060200191505060405180910390fd5b608173ffffffffffffffffffffffffffffffffffffffff16634f7c4f4c846121d08486612d46565b6040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828103825260098152602001807f62616c616e63654f6600000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b15801561228e57600080fd5b505af11580156122a2573d6000803e3d6000fd5b50505050608273ffffffffffffffffffffffffffffffffffffffff1663a0aa9f29836040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018381526020018281038252600e8152602001807f746f74616c537570706c7953756200000000000000000000000000000000000081525060200192505050600060405180830381600087803b15801561234f57600080fd5b505af1158015612363573d6000803e3d6000fd5b5050505060a173ffffffffffffffffffffffffffffffffffffffff166306e354dd6040518163ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018080602001828103825260118152602001807f757064617465546f74616c537570706c79000000000000000000000000000000815250602001915050600060405180830381600087803b15801561240857600080fd5b505af115801561241c573d6000803e3d6000fd5b505050508273ffffffffffffffffffffffffffffffffffffffff167fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5836040518082815260200191505060405180910390a2505050565b61247d3382611518565b50565b600061248d338484610ce7565b905092915050565b6124a0338383610ce7565b505050565b6124b0838383610ce7565b50505050565b6124e4336000357fffffffff0000000000000000000000000000000000000000000000000000000016612a1e565b1515612558576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f64732d617574682d756e617574686f72697a656400000000000000000000000081525060200191505060405180910390fd5b6000600160146101000a81548160ff0219169083151502179055507f1b55ba3aa851a46be3b365aee5b5c140edd620d578922f3e8466d2cbd96f954b60405160405180910390a1565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60006125f2827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6109ce565b9050919050565b6000608273ffffffffffffffffffffffffffffffffffffffff1663ce8699036040518163ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018281038252600e8152602001807f746f74616c537570706c7941646400000000000000000000000000000000000081525060200191505060206040518083038186803b15801561269a57600080fd5b505afa1580156126ae573d6000803e3d6000fd5b505050506040513d60208110156126c457600080fd5b8101908080519060200190929190505050905060008090505b818110156127e5576000608273ffffffffffffffffffffffffffffffffffffffff1663f61fe1446040518163ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018281038252600e8152602001807f746f74616c537570706c79416464000000000000000000000000000000000000815250602001915050602060405180830381600087803b15801561278857600080fd5b505af115801561279c573d6000803e3d6000fd5b505050506040513d60208110156127b257600080fd5b810190808051906020019092919050505090506127d160025482612dcb565b6002819055505080806001019150506126dd565b50608273ffffffffffffffffffffffffffffffffffffffff1663ce8699036040518163ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018281038252600e8152602001807f746f74616c537570706c7953756200000000000000000000000000000000000081525060200191505060206040518083038186803b15801561288557600080fd5b505afa158015612899573d6000803e3d6000fd5b505050506040513d60208110156128af57600080fd5b8101908080519060200190929190505050905060008090505b818110156129d0576000608273ffffffffffffffffffffffffffffffffffffffff1663f61fe1446040518163ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018281038252600e8152602001807f746f74616c537570706c79537562000000000000000000000000000000000000815250602001915050602060405180830381600087803b15801561297357600080fd5b505af1158015612987573d6000803e3d6000fd5b505050506040513d602081101561299d57600080fd5b810190808051906020019092919050505090506129bc60025482612d46565b6002819055505080806001019150506128c8565b507f91cf34a58a9ed220d9b072639106d1d07251e1791426d2560ab2d6c4ca19a1836002546040518082815260200191505060405180910390a15050565b612a19823383610ce7565b505050565b60003073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415612a5d5760019050612c8d565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415612abc5760019050612c8d565b600073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415612b1b5760009050612c8d565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b70096138430856040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001935050505060206040518083038186803b158015612c4f57600080fd5b505afa158015612c63573d6000803e3d6000fd5b505050506040513d6020811015612c7957600080fd5b810190808051906020019092919050505090505b92915050565b60008282604051602001808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c010000000000000000000000000281526014018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166c01000000000000000000000000028152601401925050506040516020818303038152906040528051906020012060019004905092915050565b60008282840391508111151515612dc5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f64732d6d6174682d7375622d756e646572666c6f77000000000000000000000081525060200191505060405180910390fd5b92915050565b60008282840191508110151515612e4a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f64732d6d6174682d6164642d6f766572666c6f7700000000000000000000000081525060200191505060405180910390fd5b9291505056fea165627a7a72305820da44e3d5b3d5cc2ff99efb133b00a95fe008fde7f2e42597a28ddd774eb7551f0029"
)

func TestDSTokenMint(t *testing.T) {
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	// db := urlcommon.NewTransientDB(persistentDB)
	db := persistentDB

	url := concurrenturl.NewConcurrentUrl(db)
	api := adaptor.NewAPI(db, url)
	statedb := adaptor.NewStateDB(api, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(tests.Coinbase)
	statedb.CreateAccount(owner)
	statedb.AddBalance(owner, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)

	// Deploy DSToken.
	eu, config := prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt := deploy(eu, config, owner, 0, dsTokenV2Code, []byte{32}, []byte{4}, []byte("TESTxxxxxxxxxxxxxxxxxxxxxxxxxxxx"))
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	dstokenAddress := receipt.ContractAddress
	t.Log(dstokenAddress)

	// Call mint twice.
	N := 10
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{1})
	totalTransitions := []urlcommon.UnivalueInterface{}
	txs := []uint32{}
	begin := time.Now()
	for i := 0; i < N; i++ {
		tdb := curstorage.NewTransientDB(db)
		url = concurrenturl.NewConcurrentUrl(tdb)
		api := adaptor.NewAPI(tdb, url)
		statedb := adaptor.NewStateDB(api, tdb, url)

		config := MainConfig()
		config.Coinbase = &coinbase
		config.BlockNumber = new(big.Int).SetUint64(10000001)
		config.Time = new(big.Int).SetUint64(10000001)

		eu := adaptor.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, tdb, url)
		transitions, receipt = run(eu, config, &owner, &dstokenAddress, uint64(i+1), false, "mint(address,uint256)", []byte{byte((i + 1) / 65536), byte((i + 1) / 256), byte((i + 1) % 256)}, []byte{1})
		if i <= 1 {
			t.Log("\n", FormatTransitions(transitions))
			t.Log(receipt)
		}
		if receipt.Status != 1 {
			t.Log(receipt)
			t.Fail()
			return
		}
		totalTransitions = append(totalTransitions, transitions...)
		txs = append(txs, uint32(i+2))

		// if (i+1)%10000 == 0 {
		// 	// t.Log("\n" + FormatTransitions(transitions))
		// 	// t.Log(receipt)
		// 	t.Log("time for exec: ", time.Since(begin))

		// 	begin = time.Now()
		// 	url = concurrenturl.NewConcurrentUrl(db)
		// 	url.Commit(totalTransitions, txs)
		// 	t.Log("time for commit: ", time.Since(begin))
		// 	begin = time.Now()
		// 	totalTransitions = []urlcommon.UnivalueInterface{}
		// 	txs = []uint32{}
		// }
	}
	t.Log("time for exec: ", time.Since(begin))
	begin = time.Now()

	// Call defer function.
	tdb := curstorage.NewTransientDB(db)
	url = concurrenturl.NewConcurrentUrl(tdb)
	url.Import(totalTransitions)
	url.PostImport()
	url.Commit(txs)
	t.Log("time for commit: ", time.Since(begin))
	begin = time.Now()

	url = concurrenturl.NewConcurrentUrl(tdb)
	api = adaptor.NewAPI(tdb, url)
	statedb = adaptor.NewStateDB(api, tdb, url)

	config = MainConfig()
	config.Coinbase = &coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000001)
	config.Time = new(big.Int).SetUint64(10000001)

	eu = adaptor.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, tdb, url)
	transitions, receipt = run(eu, config, &owner, &dstokenAddress, uint64(N+1), false, "updateTotalSupply(string)", []byte{32}, []byte{17}, append([]byte("updateTotalSupply"), []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}...))
	if receipt.Status != 1 {
		t.Log(receipt)
		t.Fail()
		return
	}
	t.Log("time for defer: ", time.Since(begin))
	t.Log("\n" + FormatTransitions(transitions))
	begin = time.Now()

	totalTransitions = append(totalTransitions, transitions...)
	txs = append(txs, uint32(N+2))

	// Commit mint and defer.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(totalTransitions)
	url.PostImport()
	url.Commit(txs)
	t.Log("time for final commit: ", time.Since(begin))
}
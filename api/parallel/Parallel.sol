// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../noncommutative/bytes/Bytes.sol";

contract Parallel is Bytes {
    uint256 numThreads = 1;
    constructor (uint256 threads) Bytes(address(0xb0)) {
        numThreads = threads; 
    }

    function run() public {       
        foreach(abi.encode(numThreads));
    }
}

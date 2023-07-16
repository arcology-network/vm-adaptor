// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;


import "../base/Base.sol";

contract String is Base {
    constructor() Base(address(0x84)) {}

    function push(string memory elem) public virtual{ //9e c6 69 25
         Base.insert(Base.rand(), abi.encode(elem));
    }    

    function pop() public virtual returns(string memory) { // 80 26 32 97
        return abi.decode(Base.popBack(), (string));  
    }

    function get(uint256 idx) public virtual returns(string memory)  { // 31 fe 88 d0
        return abi.decode(Base.getIndex(idx), (string));  
    }

    function set(uint256 idx, string memory elem) public { // 7a fa 62 38
        Base.setIndex(idx, abi.encode(elem));
    }
}

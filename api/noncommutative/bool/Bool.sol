// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../base/Base.sol";

contract Bool is Base {
    constructor() Base(address(0x84)) {}

    function push(bool elem) public virtual { //9e c6 69 25
        Base.insert(Base.rand(), abi.encode(elem));
    }    

    function pop() public virtual returns(bool) { // 80 26 32 97
        return abi.decode(Base.popBack(), (bool));  
    }

    function get(uint256 idx) public virtual  returns(bool)  { // 31 fe 88 d0
        return abi.decode(Base.getIndex(idx), (bool));  
    }

    function set(uint256 idx, bool elem) public { // 7a fa 62 38
        Base.setIndex(idx, abi.encode(elem));    
    }
}

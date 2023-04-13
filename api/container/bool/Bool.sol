pragma solidity ^0.5.0;

import "./Base.sol";

contract Bool {
    Base base;

    constructor  () public {  base = new Base(); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(bool) { 
        return abi.decode(abi.decode(base.pop(), (bytes)), (bool));  
    }

    function push(bool elem) public { 
       base.push(abi.encodeWithSignature("push(bytes, bytes)",  base.id(), abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(bool)  { 
        return abi.decode(abi.decode(base.get(idx), (bytes)), (bool));
    }

    function set(uint256 idx, bool elem) public {
        base.set(abi.encodeWithSignature("set(bytes, uint256, bytes)", base.id(), idx, abi.encode(elem)));        
    }
}
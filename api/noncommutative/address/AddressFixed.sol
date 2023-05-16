pragma solidity ^0.5.0;

import "./Address.sol";

contract AddressFixed {
    Address array;
    constructor  (uint length, address initialV) public {  
        array = new Address(); 
        for (uint i = 0; i < length; i ++) {
            array.push(initialV);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(address)  {return array.get(idx);}
    function set(uint256 idx, address elem) public { array.set(idx, elem); }
}

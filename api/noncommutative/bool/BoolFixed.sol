pragma solidity ^0.5.0;

import "./Bool.sol";

contract BoolFixed {
    Bool array;
    constructor  (uint length, bool initialV) public {  
        array = new Bool(); 
        for (uint i = 0; i < length; i ++) {
            array.push(initialV);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(bool)  {return array.get(idx);}
    function set(uint256 idx, bool elem) public { array.set(idx, elem); }
}

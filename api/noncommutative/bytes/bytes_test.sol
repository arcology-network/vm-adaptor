pragma solidity ^0.5.0;

import "./Bytes.sol";

contract ByteTest {
    Bytes container = new Bytes();
    
    constructor() public {     
        require(container.length() == 0); 
 
        bytes memory arr1 = '0x1000000000000000000000000000000000000000000000000000000000000001';
        bytes memory arr2 = '0x2000000000000000000000000000000000000000000000000000000000000002';

        container.push(arr1);  
        container.push(arr1); 

        require(container.length() == 2); 


        require(keccak256(container.get(1)) == keccak256(arr1));

        container.set(1, arr2);       

        require(keccak256(container.get(0)) == keccak256(arr1));
        require(keccak256(container.get(1)) == keccak256(arr2));
        require(keccak256(container.pop()) == keccak256(arr2));

        container.pop();
        require(container.length() == 0); 
    }
}
pragma solidity ^0.5.0;

import "./Int256Cumulative.sol";
import "./Threading.sol";

contract ThreadingInt64 {
    Int256Cumulative container = new Int256Cumulative(0, 100);

   //  function call() public {
   //     Threading mp = new Threading();
   //     mp.add(address(this), abi.encodeWithSignature("add(int64)", 1));
   //     mp.add(address(this), abi.encodeWithSignature("add(int64)", 2));      
       


   //     mp.run(1);
   //     assert(container.length() == 1 ); 
   //  }

   //  function push(int64 elem) public { //9e c6 69 25
   //     container.push(elem);
   //  }  
}
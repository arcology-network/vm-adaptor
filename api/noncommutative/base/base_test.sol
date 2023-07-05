// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract BaseTest {    
    address constant public API = address(0x84); 

    uint[] public arr2 = [1, 2, 3];
    bytes private id;

    event logMsg(string message);

    constructor() {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("new()"));       
        require(success, "Bytes.New() Failed");
        id = data;
 
        bytes memory byteArray = new bytes(75);
        for (uint  i = 0; i < 75; i ++) {
            byteArray[i] = 0x41;
        }
    
        require(peek() == 0);  
        require(length() == 0); 
        push(byteArray);  
        push(byteArray);          
        require(length() == 2); 
        require(peek() == 0);  

        bytes memory stored = get(1);
        require(stored.length == byteArray.length);
        for (uint  i = 0; i < byteArray.length; i ++) {
            require(stored[i] == byteArray[i]);
        }

        bytes memory elems = new bytes(5);
        for (uint  i = 0; i < elems.length; i ++) {
            elems[i] = 0xaa;
        }
        set(1, elems);
       
        stored = get(0);
        require(stored.length == byteArray.length);
        for (uint  i = 0; i < byteArray.length; i ++) {
            require(stored[i] == byteArray[i]);
        }

        stored = get(1);
        require(stored.length == elems.length); 
        for (uint  i = 0; i < elems.length; i ++) {
            require(stored[i] == elems[i]);
        }

        stored = pop();
        for (uint  i = 0; i < elems.length; i ++) {
            require(stored[i] == elems[i]);
        }
        require(length() == 1); 
        require(peek() == 0);  
    }

    function call() public{ 
        require(peek() == 1); 
        pop();
        require(peek() == 1); 
    }

    function peek() public returns(uint256) {
        (,bytes memory data) = address(API).call(abi.encodeWithSignature("peek()"));
        if (data.length > 0) {
            return abi.decode(data, (uint256));   
        }
        return 0;   
    }

    function length() public returns(uint256) {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("length()"));
        require(success, "Bytes.length() Failed");
        return  abi.decode(data, (uint256));
    }

    function pop() public returns(bytes memory) {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("pop()"));
        require(success, "Bytes.pop() Failed");
        return abi.decode(data, (bytes)); 
    }

    function push(bytes memory elem) public {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("push(bytes)", elem));
        require(success, "Bytes.push() Failed");
    }   

    function get(uint256 idx) public returns(bytes memory)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("get(uint256)", idx));
        require(success, "Bytes.get() Failed");
        return abi.decode(data, (bytes));  
    }

    function set(uint256 idx, bytes memory elem) public {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("set(uint256,bytes)", idx, elem));
        require(success, "Bytes.set() Failed");
    }
}

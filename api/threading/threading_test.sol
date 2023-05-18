pragma solidity ^0.5.0;

import "./Threading.sol";

contract ThreadingTest {
    function call() public  { 
       bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";

       Threading mp = new Threading();
       mp.add(address(this), abi.encodeWithSignature("hasher(bytes)", data));
       mp.add(address(this), abi.encodeWithSignature("hasher(bytes)", data));
       assert(mp.length() == 2);

       mp.del(0);
       assert(mp.length() == 1);

       (bool success,) = address(address(0x90)).call(abi.encodeWithSignature("run(uint8)", 1));   
       assert(success);

       (,bytes memory hash) = mp.get(0);
       bytes32 hash32 = bytesToBytes32(hash); 
       assert(hash32 == keccak256(data));

       assert(mp.length() == 1);
       mp.clear();
       assert(mp.length() == 0);
    }

    function hasher(bytes memory data) pure public returns(bytes32){
      return keccak256(data);
    }

    function bytesToBytes32(bytes memory b) private pure returns (bytes32) {
        bytes32 out;
        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}


contract Benchmarking {
    function reverse(bytes memory message) public pure {
        bytes memory reversed = new bytes(message.length);
        for(uint i=0;i<message.length;i++){
            reversed[ message.length - i - 1] = message[i];
        }
        keccak256(reversed); 
    }

    function testReverseString10k() public {
        string memory message = "roaA1Cjl0SOW9WwjnVEqt9wPAEyvothOcw6HwIsmW5sKz2YlonlBGJY0sSGGF8KYBoPcpbwqkBEfjakmRxz0sWNuXomwQvEUEI1ZssEzwOOeXY1dxVmmUGRK2rdwVsQaP0ZzR7oP6rQMSLvWpz1jheOU4yZ0fh7lhTgWBpCstwbaKBXCHpTYM4ipHFo8kGwMnKHv2XlvGR7Jj8e4XO3p4K3BgWgEPZ1jg0QXGhIkDnHI8tNlHCfLJZzCgShbQWDhvlFfdHeFqJpBBao44gHgXzYuHHssPzBhmSsha2ZzanFKaoSGKWuJ2n7lPqmBdy98EwLClBb2LgtdeIBRRZ0TIZkRrWeHITf5oxkLyjOeRRGBDkacA25dPtCfAWb6YGqMbLZthWiPepm5enYQe4nvYIj6bGejeTHKzknmLCHO0HoRlpjmhALh2445Kk2IWLkESX0ugfahuKPhCBhCEsUh0g5mTxUOCZFCNkFViSqmgUTQLqlXU4JxzkIPrOvNiv12mJRRkfxZ5bJMeOqudERviXzIj1Onb1JddtnCzV47jJflkhLRmLxGfbMaxl0t2YwEMME7qpwNzaGPWC0yt5WgHafDu7y6WvQRX61Tsi185bDwk2369QiM30Q6tn3KISdynXspaeAf4xh2wNRlYl5a2TQyye8W3yVep96sUdW9sntA1kZeIiVbxY9VUUzGythYNlvb3X0zLaDfbNKiLmhr4uqBLxcue1oWvKz3PHP2FyqmxAn1dtjzcbIQOGLGrygayXtwqyRpplhuCnjjoEj8WBBL48vqNAlyLECwTZVrbA9SVHUgxAQCuFKjf4iN2TtVwCcoqZyJtNiqjfH1v8K7nPD5LjtFzAc2mVefOllbadsK1fAqKKtgWIUArHd1z3j85nqzw0SbYof4nrk1GglGh7XgPpyD6HlLi3ynzdIR6SRBuUS7TTA1I1RIjTNmhxD1AClx39RFFHodsGpKo1L0NC6ETl971LzlmQqFoFuJsPeC1iIfMvgiOTQEoLxHZ5k5hvj3xEdzKMdQIVyrybFfl285Ev9S3n6B0ye9D8wsdfrJ1wEeYJWLi4BOd7GXHzoo4Vc2jR4QU5vEGFNIu7mkLea3S4MqLx8rWfr94TKfrQda7XOMdQjVswIDAp7MjEIOPOXcEw1ieooNS0CC2R9XpU8yLIZ4zXpS6e9Kgj4sHoKAMYTZOb6TysslCKlSvZt7PNC2K4A8ZNEeXwdfVhbXEt752711gu5oBbGfgooDySJnkAPWR21tYKRAr0uTNWt9Bzi9eo9k7V3Rw2Yl6ZXRc4k8gqYwLp4k9rzRGwCBRKBOqho85uOuDj4EoZrK67SYuCjDtlG9Pf2Y7BAf2JM6yGIWXhqyQabICvsmhBWcx4riVmikzae5novWFM6wF1rzKvJK2N9dbotEmy6vp1Vq6IW70XycJXNFlYrdGpwl9XXFUt7WMiFFgqPvaiKBtnnMB0AP2f2XkvM0oAkqSF445L0NFLR4j4B6iSeLdLlHysH6hXpu7f4lsEWsCpKfJcBftxSmKI8bt3RI8QPSlOb59ujvXAQ6uQVcRLdkZbhAMYFr7rJrFbXkke3LMZkEjWK7mqowFGR7au9dUFINLtA29964arhSvfAzbuI3e02MGc6li2gdG1NXC5KS2EqTeJd1iOhGdoJnwBcOiB2toeTHnGHwQ5yv8LK6e2rO0CFzNemzZRoVXOKq0wxviYK4xylHbhcLACG8PkZ1KXKfO9VWoC22Ya1fzD6tg2M8zis4ZXqGexgi0rL684fbrIcmvQIL1D3rQORWyNaONuTsVz3jWLu3XdiFLXUKt2SXkJmZQE1nVLYxaEOO5yfAMdWiOMNXNsVT4nUUGeBVmQwOT6XpBRtRfA4EmRTAljn10p7wJFRpjRoft2CWL91OGQCMsomMTpDog98OTMKsYPL8qgvva6odKGbEUVWVjba6eIFwBqvYJNonNvKfM9YnWAFAE5YHOW9NsnnOoJ1ke2JyhnwFuyuKdFKyM1XjuS6OBfP4nxfKVIysp0uyaIS8d4Pzy7gM5rurycADTrWigP6urOQ2ug90QMP4KSvpEYSrtbKAhrIZ8tq5tdlnTzVV90Q6qD22FESEwH8aGnPUMbPzijt34z7V6F0b077qLVgl5bGoJ1OovF2hHbuama5oqv8RPjxyk9JOuepIVsIAD9S1hBtEdH7UaMUZUsb0we48yr7A5WWSH4hKI5bBwzlXNtKS14FuLD1akozrhNpCDBJ8EIrAQiskRBd8g7OZljpk4Qvb20BFIAol6gqI1WyqXReKXVeN4X7KI6vsZrtGNNAv6pTCFrusm8871KIUIy7zR32BvyIZegY5TfJY1CX74XhLkUsrZze08K1LzkhDSC8w7xZfSppD617fO70lvGvMrky8lZFRu6VT3oXvtvsv8ls8CRx0uIWBQLWC94hwZtA9pPUDJpQJrtjiT1e6DDsSU8Ytna280hbCPPK5jImAHDvmy8pBUAPo8b5cCQV1Is5WWctc5DrvbPbq66248BW7mWqZ8C1AQJMsRquouQZdw3ziPv9RJunVyGo3MiHQtFuDJO2mMjHLWvXFcciOe1Wycf9FP6mujwBtfYNB9BtOqfxiedapVIV5zFcHtPYJEEcIxvC52RwVn1AxBg2YB0Ff77N4WPMCSz7p0s1qEXEvb67A9zSq0qeK1uzJ9h09sVI19CDqv65W33wEmhpIkBvAZAC6J8AgCe0TCqHlVL7Hb4FbrHX7gzRzHz4A27Ky5Z0KyW67zJpLp9jbo66EqiQdPjqyUP1HRY0NpwYxjabf86BFLzeP17gs0esr7875FKfZyvblPm48rRLHGE6UoECwFc9sEldDG1DIJGkB6CLHu8tyskev9xmMzUergHulbVWalGqZRUG3GppxHVq3cNcPHuUUOqXZoY9Eip3J9eTBu41uNXtD0ru6V4FmPxiyQCz4T0OFoSseOOSfDnA9kYZ25ohuoJ5eC2bPv9fTuZ1nU5U7elSGf3CJHIV3xy5M9WCMb8jViVes5sjSplV2HFNe9XliSOVoah396u8bs33iOETA83bO7sHbAfHBzW8SoQ6OvsyUb4Kidqj81c7QQKAWCg7IRaP6gevnbz7VBAbucMll2ezTnifFBY3GU75sRlpn45yHG7EO26l4K2N0UNN4qDItABKVN8mP6nEo4r91oFC91lZI1erkcfxJahQT0KKjIbCSdxCVwyx9bDgBq67MbPVuX39QdgWUT0m5iVyQcVpzfokr66xKgcBDNvESwwqjRkx1ci1KUOCAOGDLcZUJfk6l5qfvlO5H2dF7NqO79DJJhLgSRoNJf93CIdgBDAfemWPjvbllieqtqJUXph1wOVwmottHIBSHwsrTK5SfmxYhlnawGf7xR6W1UemL4CusDinIabe8cBQEBoLMzeKV1kwcjGQbExkPDoYXODIXcbueuY7eetdXDm506t1WXEm0DfN8Idz2Ol0MzTvDW1l6PjGFsGpH3V5gCV6JYfoWHMp1LHVb9Hfxn9JmwjKS4etZtDUudTFbmZqxw8BX1ajPLVVymn8ImlZy8VXcytAmFsJIwpMK7X1SatJGvqqKkhpjH3oh32ChE3MSQDlEYFhrjucgassI5oNT139MLnbRkVSf81cnwD4jObP6bTEndXGw3NDEAOQxppklB1lVvPfb76PPPXsyNBPiQGko8AHkO1a4i8K6IlRVXyYiUPWmYSzfqcl2IjZ04IIY4jJ4L0hwS5CkqJXFlqkl9syarXhX9hlrtoEYbAfcaTrLmMMkK89Nq1PkhPmAZSNaGNLIiBrEeDn3znnjH6MK4rWiWMQMmLVddpvXAQtYwgF3EUpPRwZPcjpDI06pdYq0Qr2S7nZm32ej9WE1EegSPOxzq7AQAsHM0KTLQtUi6JdOr1dOJF5EdybIeGS2ziAd6afKdN6f6SDWLzY9Usgcd15C0Vyr99SePVj5veSpvK02qEJ4E59agV2dRJuxg1ewv01a9xvn058kJUKGaOFmoA5cepzN71ty6xtGn9qJWDDVnlUZJhXmVMb3jOFXka41InjWreqijZv3qDZlVH1ynx1kwie9fMGTtYYqazZTXDJXOm0BcQWeOiS7EPqEK7V0VUIqpF7L4BYfOIgJpy5QfEBjv2ADShlq73CbYv24bLTu5xlisVCYNVj81vDxYqjorgm7D0PlDFmZpat9SjdFxWkynKNmWWdFb1BxAKHrQ5oF5vZbPOYywpwveNzuJNdXxDzSSxwjKE4fqUWPgeFPj58nGsjqFCLH8rPLop3Q1b9RoYQYkNjXW7aoj2MQIVGfjUr8yBP2sIgppillOF4cUehwjrE0rk4TR7Fr4GHuKkNEqxOixByoP9JswwEJSXGsN4OBawp8WBh8d2CSi764hyqcYPAkxvIBtnYa0VM8Js8SYZnT6wtA8SCfXyTvfcLx9sZ2HnpRCjqdUagCTKVa37kVldWd9nfMhxGaT9FDd83Qn7pu0HJ6X4M5jddXNgwefj9z9BhZZe0eZn04EQDVBUQHhTn185RthlrUrgW8uRo2yHXq2oOYTzHc8RusuSL7xBj2RoGSfVbJDZnRmEHbKBlWek1QvLYj6U9aanHvcK9bClEdahuTIbhNmZviMo18qBLgGvhTzoJIw5hRu2wZnkZ0P4utvIgyFOCFa7P1BrLOqpRVQqbNfvU3Vh0s0bL0gEYVodHAfrP09EdLBX2FtlV1Jh6fz1mOjK2w9oEZgaZyHOIob3pQ6Xxfx44l8mEovngtqQqDbczpgHmpfLsPYna5JayHT4IMYytYKtE3BnmGQMqH6g8O0lipl0mtgeKCWjSyKC0TMXFomMWEfIIynDHaZE043lni370ek2iD2MYe9D53LIejuagnrknhM9HUV6uA6JYF6Exx7678JcXxVyZZfVFfYvdwUXeSAM5l0h403YzuQm3CEFKhcdnq1LOpIPQCghV6l3rMEeCvjzEYeIAPyRXuGnfjSrx6q1YPFMXd";
        Threading mp = new Threading();
        for (int i = 0; i < 1000; i ++) { 
            mp.add(address(this), abi.encodeWithSignature("reverse(bytes)", address(this), message));
        }
        mp.run(16);
    }
}
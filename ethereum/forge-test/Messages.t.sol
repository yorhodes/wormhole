// test/Messages.sol
// SPDX-License-Identifier: Apache 2

pragma solidity ^0.8.0;

import "../contracts/Messages.sol";
import "../contracts/Setters.sol";
import "../contracts/Structs.sol";
import "forge-std/Test.sol";

contract TestMessages is Messages, Test, Setters {
  address constant testGuardianPub = 0xbeFA429d57cD18b7F8A4d91A2da9AB4AF05d0FBe;

  // A valid VM with one signature from the testGuardianPublic key
  bytes validVM = hex"01000000000100867b55fec41778414f0683e80a430b766b78801b7070f9198ded5e62f48ac7a44b379a6cf9920e42dbd06c5ebf5ec07a934a00a572aefc201e9f91c33ba766d900000003e800000001000b0000000000000000000000000000000000000000000000000000000000000eee00000000000005390faaaa";

  function testQuorum() public {
    assertEq(quorum(0), 1);
    assertEq(quorum(1), 1);
    assertEq(quorum(2), 2);
    assertEq(quorum(3), 3);
    assertEq(quorum(4), 3);
    assertEq(quorum(5), 4);
    assertEq(quorum(6), 5);
    assertEq(quorum(7), 5);
    assertEq(quorum(8), 6);
    assertEq(quorum(9), 7);
    assertEq(quorum(10), 7);
    assertEq(quorum(11), 8);
    assertEq(quorum(12), 9);
    assertEq(quorum(19), 13);
    assertEq(quorum(20), 14);
  }

  function testQuorumCanAlwaysBeReached(uint numGuardians) public {
    if (numGuardians == 0) {
      return;
    }
    if (numGuardians >= 256) {
      vm.expectRevert("too many guardians");
    }
    // test that quorums is never greater than the number of guardians
    assert(quorum(numGuardians) <= numGuardians);
  }

  // This test ensures that submitting invalid signatures for non-existent
  // guardians fails.
  //
  // The main purpose of this test is to ensure that there's no surprising
  // behaviour arising from solidity's handling of invalid signatures and out of
  // bounds memory access. In particular, pubkey recovery of an invalid
  // signature returns 0, and in some cases out of bounds memory access also
  // just returns 0.
  function testOutOfBoundsSignature() public {
    // Initialise a guardian set with a single guardian.
    address[] memory keys = new address[](1);
    keys[0] = testGuardianPub;
    Structs.GuardianSet memory guardianSet = Structs.GuardianSet(keys, 0);
    require(quorum(guardianSet.keys.length) == 1, "Quorum should be 1");

    // Two invalid signatures, for guardian index 2 and 3 respectively.
    // These guardian indices are out of bounds for the guardian set.
    bytes32 message = "hello";
    Structs.Signature memory bad1 = Structs.Signature(message, 0, 0, 2);
    Structs.Signature memory bad2 = Structs.Signature(message, 0, 0, 3);
    // ecrecover on an invalid signature returns 0 instead of reverting
    require(ecrecover(message, bad1.v, bad1.r, bad1.s) == address(0), "ecrecover should return the 0 address for an invalid signature");

    Structs.Signature[] memory badSigs = new Structs.Signature[](2);
    badSigs[0] = bad1;
    badSigs[1] = bad2;
    vm.expectRevert(bytes("guardian index out of bounds"));
    verifySignatures(0, badSigs, guardianSet);
  }

  // This test checks the possibility of getting a unsigned message verified through verifyVM
  function testHashMismatchedVMIsNotVerified() public {
    // Set the initial guardian set
    address[] memory initialGuardians = new address[](1);
    initialGuardians[0] = testGuardianPub;

    // Create a guardian set
    Structs.GuardianSet memory initialGuardianSet = Structs.GuardianSet({
      keys: initialGuardians,
      expirationTime: 0
    });

    storeGuardianSet(initialGuardianSet, uint32(0));

    // Confirm that the test VM is valid
    (Structs.VM memory parsedValidVm, bool valid, string memory reason) = this.parseAndVerifyVM(validVM);
    require(valid, reason);
    assertEq(valid, true);
    assertEq(reason, "");

    // Manipulate the payload of the vm
    Structs.VM memory invalidVm = parsedValidVm;
    invalidVm.payload = abi.encodePacked(
        parsedValidVm.payload,
        "malicious bytes in payload"
    );

    // Confirm that the verifyVM fails on invalid VM
    (valid, reason) = this.verifyVM(invalidVm);
    assertEq(valid, false);
    assertEq(reason, "vm.hash doesn't match body");
  }
}

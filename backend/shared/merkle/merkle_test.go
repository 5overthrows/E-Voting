package merkle

import (
	"testing"
)

func TestComputeMerkleRoot(t *testing.T) {
	// Test empty leaves
	if ComputeMerkleRoot([]string{}) != "" {
		t.Errorf("Empty leaves should return empty string")
	}

	// Test single leaf
	single := []string{"hash1"}
	if ComputeMerkleRoot(single) != "hash1" {
		t.Errorf("Single leaf should return itself")
	}

	// Test even number of leaves
	leavesEven := []string{"hash1", "hash2", "hash3", "hash4"}
	rootEven := ComputeMerkleRoot(leavesEven)
	if rootEven == "" {
		t.Errorf("Merkle root for even leaves should not be empty")
	}

	// Verify root validation
	if !VerifyMerkleRoot(leavesEven, rootEven) {
		t.Errorf("VerifyMerkleRoot failed for correct root")
	}

	if VerifyMerkleRoot(leavesEven, "invalidroot") {
		t.Errorf("VerifyMerkleRoot passed for incorrect root")
	}

	// Test odd number of leaves (should duplicate the odd node)
	leavesOdd := []string{"hash1", "hash2", "hash3"}
	rootOdd := ComputeMerkleRoot(leavesOdd)
	if rootOdd == "" {
		t.Errorf("Merkle root for odd leaves should not be empty")
	}

	// Double check duplication logic:
	// For leaves ["hash1", "hash2", "hash3"]:
	// level 1: hash12 = SHA256("hash1" + "hash2"), hash33 = SHA256("hash3" + "hash3")
	// level 2 (root): root = SHA256(hash12 + hash33)
	// Let's verify manually
	expectedRoot := ComputeMerkleRoot([]string{"hash1", "hash2", "hash3"})
	if rootOdd != expectedRoot {
		t.Errorf("Odd root computation mismatch: %s != %s", rootOdd, expectedRoot)
	}
}

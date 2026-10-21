package merkle

import (
	"crypto/sha256"
	"fmt"
)

// ComputeMerkleRoot calculates the Merkle Root from a slice of leaf hashes (hex strings)
func ComputeMerkleRoot(leaves []string) string {
	if len(leaves) == 0 {
		return ""
	}

	level := make([]string, len(leaves))
	copy(level, leaves)

	for len(level) > 1 {
		var nextLevel []string
		for i := 0; i < len(level); i += 2 {
			if i+1 < len(level) {
				combined := level[i] + level[i+1]
				h := sha256.Sum256([]byte(combined))
				nextLevel = append(nextLevel, fmt.Sprintf("%x", h))
			} else {
				// Odd node: duplicate it
				combined := level[i] + level[i]
				h := sha256.Sum256([]byte(combined))
				nextLevel = append(nextLevel, fmt.Sprintf("%x", h))
			}
		}
		level = nextLevel
	}

	return level[0]
}

// VerifyMerkleRoot verifies if the recomputed Merkle root matches the expected root
func VerifyMerkleRoot(leaves []string, expectedRoot string) bool {
	computedRoot := ComputeMerkleRoot(leaves)
	return computedRoot == expectedRoot
}

package fcrmerkletrie

import (
	"crypto/sha256"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
)

// FCRMerkleProof is the proof of a single cid in a merkle trie
type FCRMerkleProof struct {
	path  [][]byte
	index []int64
}

// VerifyCID is used to verify a given cid and a given root matches the proof
func (mp *FCRMerkleProof) VerifyCID(cid *cid.ContentID, root string) bool {
	currentHash, _ := cid.CalculateHash()
	for i, path := range mp.path {
		hashFunc := sha256.New()
		if mp.index[i] == 1 {
			hashFunc.Write(append(currentHash, path...))
		} else {
			hashFunc.Write(append(path, currentHash...))
		}
		currentHash = hashFunc.Sum(nil)
	}
	return string(currentHash) == root
}

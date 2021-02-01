package fcrmerkletrie

import (
	"github.com/cbergoon/merkletree"
)

// FCRMerkleTrie is used to store
type FCRMerkleTrie struct {
	trie *merkletree.MerkleTree
}

// CreateMerkleTrie creates a merkle trie from a list of cids
func CreateMerkleTrie(contents []merkletree.Content) (*FCRMerkleTrie, error) {
	trie, err := merkletree.NewTree(contents)
	if err != nil {
		return nil, err
	}
	return &FCRMerkleTrie{trie: trie}, nil
}

// GetMerkleRoot returns the merkle root of the trie
func (mt *FCRMerkleTrie) GetMerkleRoot() string {
	return string(mt.trie.MerkleRoot())
}

// GenerateMerkleProof gets the merkle proof for a given cid
func (mt *FCRMerkleTrie) GenerateMerkleProof(content merkletree.Content) (*FCRMerkleProof, error) {
	path, index, err := mt.trie.GetMerklePath(content)
	if err != nil {
		return nil, err
	}
	return &FCRMerkleProof{path: path, index: index}, nil
}

package fcrmerkletrie

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/cid"
	"github.com/cbergoon/merkletree"
)

// FCRMerkleTrie is used to store
type FCRMerkleTrie struct {
	trie *merkletree.MerkleTree
}

// CreateMerkleTrie creates a merkle trie from a list of cids
func CreateMerkleTrie(cids []cid.ContentID) (*FCRMerkleTrie, error) {
	size := len(cids)
	list := make([]merkletree.Content, size)
	for i := 0; i < size; i++ {
		list[i] = &cids[i]
	}
	trie, err := merkletree.NewTree(list)
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
func (mt *FCRMerkleTrie) GenerateMerkleProof(cid *cid.ContentID) (*FCRMerkleProof, error) {
	path, index, err := mt.trie.GetMerklePath(cid)
	if err != nil {
		return nil, err
	}
	return &FCRMerkleProof{path: path, index: index}, nil
}

package ethtype

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"slices"

	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/internal/kzg4844"
	"github.com/donutnomad/eths/rlp"
)

// BlobTxSidecar contains the blobs of a blob transaction.
type BlobTxSidecar struct {
	Version     byte                 // Version
	Blobs       []kzg4844.Blob       // Blobs needed by the blob pool
	Commitments []kzg4844.Commitment // Commitments needed by the blob pool
	Proofs      []kzg4844.Proof      // Proofs needed by the blob pool
}

// NewBlobTxSidecar initialises the BlobTxSidecar object with the provided parameters.
func NewBlobTxSidecar(version byte, blobs []kzg4844.Blob, commitments []kzg4844.Commitment, proofs []kzg4844.Proof) *BlobTxSidecar {
	return &BlobTxSidecar{
		Version:     version,
		Blobs:       blobs,
		Commitments: commitments,
		Proofs:      proofs,
	}
}

// BlobHashes computes the blob hashes of the given blobs.
func (sc *BlobTxSidecar) BlobHashes() []common.Hash {
	hasher := sha256.New()
	h := make([]common.Hash, len(sc.Commitments))
	for i := range sc.Blobs {
		h[i] = kzg4844.CalcBlobHashV1(hasher, &sc.Commitments[i])
	}
	return h
}

// CellProofsAt returns the cell proofs for blob with index idx.
// This method is only valid for sidecars with version 1.
func (sc *BlobTxSidecar) CellProofsAt(idx int) ([]kzg4844.Proof, error) {
	if sc.Version != BlobSidecarVersion1 {
		return nil, fmt.Errorf("cell proof unsupported, version: %d", sc.Version)
	}
	if idx < 0 || idx >= len(sc.Blobs) {
		return nil, fmt.Errorf("cell proof out of bounds, index: %d, blobs: %d", idx, len(sc.Blobs))
	}
	index := idx * kzg4844.CellProofsPerBlob
	if len(sc.Proofs) < index+kzg4844.CellProofsPerBlob {
		return nil, fmt.Errorf("cell proof is corrupted, index: %d, proofs: %d", idx, len(sc.Proofs))
	}
	return sc.Proofs[index : index+kzg4844.CellProofsPerBlob], nil
}

// ToV1 converts the BlobSidecar to version 1, attaching the cell proofs.
func (sc *BlobTxSidecar) ToV1() error {
	if sc.Version == BlobSidecarVersion1 {
		return nil
	}
	if sc.Version == BlobSidecarVersion0 {
		proofs := make([]kzg4844.Proof, 0, len(sc.Blobs)*kzg4844.CellProofsPerBlob)
		for i := range sc.Blobs {
			cellProofs, err := kzg4844.ComputeCellProofs(&sc.Blobs[i])
			if err != nil {
				return err
			}
			proofs = append(proofs, cellProofs...)
		}
		sc.Version = BlobSidecarVersion1
		sc.Proofs = proofs
	}
	return nil
}

// encodedSize computes the RLP size of the sidecar elements. This does NOT return the
// encoded size of the BlobTxSidecar, it's just a helper for tx.Size().
func (sc *BlobTxSidecar) encodedSize() uint64 {
	var blobs, commitments, proofs uint64
	for i := range sc.Blobs {
		blobs += rlp.BytesSize(sc.Blobs[i][:])
	}
	for i := range sc.Commitments {
		commitments += rlp.BytesSize(sc.Commitments[i][:])
	}
	for i := range sc.Proofs {
		proofs += rlp.BytesSize(sc.Proofs[i][:])
	}
	return rlp.ListSize(blobs) + rlp.ListSize(commitments) + rlp.ListSize(proofs)
}

// ValidateBlobCommitmentHashes checks whether the given hashes correspond to the
// commitments in the sidecar
func (sc *BlobTxSidecar) ValidateBlobCommitmentHashes(hashes []common.Hash) error {
	if len(sc.Commitments) != len(hashes) {
		return fmt.Errorf("invalid number of %d blob commitments compared to %d blob hashes", len(sc.Commitments), len(hashes))
	}
	hasher := sha256.New()
	for i, vhash := range hashes {
		computed := CalcBlobHashV1(hasher, &sc.Commitments[i])
		if vhash != computed {
			return fmt.Errorf("blob %d: computed hash %#x mismatches transaction one %#x", i, computed, vhash)
		}
	}
	return nil
}

// Copy returns a deep-copied BlobTxSidecar object.
func (sc *BlobTxSidecar) Copy() *BlobTxSidecar {
	return &BlobTxSidecar{
		Version: sc.Version,

		// The element of these slice is fix-size byte array,
		// therefore slices.Clone will actually deep copy by value.
		Blobs:       slices.Clone(sc.Blobs),
		Commitments: slices.Clone(sc.Commitments),
		Proofs:      slices.Clone(sc.Proofs),
	}
}

// CalcBlobHashV1 calculates the 'versioned blob hash' of a commitment.
// The given hasher must be a sha256 hash instance, otherwise the result will be invalid!
func CalcBlobHashV1(hasher hash.Hash, commit *kzg4844.Commitment) (vh [32]byte) {
	if hasher.Size() != 32 {
		panic("wrong hash size")
	}
	hasher.Reset()
	hasher.Write(commit[:])
	hasher.Sum(vh[:0])
	vh[0] = 0x01 // version
	return vh
}

const (
	// BlobSidecarVersion0 includes a single proof for verifying the entire blob
	// against its commitment. Used when the full blob is available and needs to
	// be checked as a whole.
	BlobSidecarVersion0 = byte(0)

	// BlobSidecarVersion1 includes multiple cell proofs for verifying specific
	// blob elements (cells). Used in scenarios like data availability sampling,
	// where only portions of the blob are verified individually.
	BlobSidecarVersion1 = byte(1)
)

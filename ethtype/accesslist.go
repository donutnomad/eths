package ethtype

import (
	"github.com/donutnomad/eths/ecommon"
)

// AccessList is an EIP-2930 access list.
type AccessList []AccessTuple

// AccessTuple is the element type of an access list.
type AccessTuple struct {
	Address     ecommon.Address `json:"address"`
	StorageKeys []ecommon.Hash  `json:"storageKeys"`
}

package ethtype

import (
	"github.com/donutnomad/eths/ecommon"
)

//
//func (a AccessList) To() ethTypes.AccessList {
//	var ret = make(ethTypes.AccessList, len(a))
//	for i, item := range a {
//		tuple := ethTypes.AccessTuple{
//			Address: item.Address.To(),
//		}
//		for _, k := range item.StorageKeys {
//			tuple.StorageKeys = append(tuple.StorageKeys, common.Hash(k))
//		}
//		ret[i] = tuple
//	}
//	return ret
//}

// AccessTuple is the element type of an access list.
type AccessTuple struct {
	Address     ecommon.Address `json:"address"`
	StorageKeys []ecommon.Hash  `json:"storageKeys"`
}

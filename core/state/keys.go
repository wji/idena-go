package state

import (
	"encoding/binary"
	"github.com/idena-network/idena-go/common"
	"github.com/pkg/errors"
	dbm "github.com/tendermint/tm-db"
)

var (
	//global db keys
	currentStateDbPrefixKey             = []byte{0x1}
	currentIdentityStateDbPrefixKey     = []byte{0x2}
	preliminaryIdentityStateDbPrefixKey = []byte{0x3}

	//state prefixes
	stateDbPrefixBytes            = []byte{0x1}
	identityStateDbPrefixBytes    = []byte{0x2}
	preliminaryStateDbPrefixBytes = []byte{0x3}

	//state db prefixes and keys
	addressPrefix            = []byte{0x1}
	identityPrefix           = []byte{0x2}
	globalKey                = []byte{0x3}
	statusSwitchKey          = []byte{0x4}
	contractStorePrefix      = []byte{0x5}
	delegationSwitchKey      = []byte{0x6}
	delayedOfflinePenaltyKey = []byte{0x7}
	burntCoinsPrefix         = []byte{0x8}
	contractCodePrefix       = []byte{0x9}
)

var (
	StateDbKeys         = &stateDbKeys{}
	IdentityStateDbKeys = &identityStateDbPrefix{}
)

type stateDbKeys struct {
}

func (s *stateDbKeys) LoadDbPrefix(db dbm.DB) ([]byte, error) {
	p, err := db.Get(currentStateDbPrefixKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get value")
	}
	if p == nil {
		p = s.BuildDbPrefix(0)
		b := db.NewBatch()
		s.SaveDbPrefix(b, p)
		err := b.WriteSync()
		return p, errors.Wrap(err, "failed to write value")
	}
	return p, nil
}

func (s *stateDbKeys) SaveDbPrefix(b dbm.Batch, prefix []byte) {
	b.Set(currentStateDbPrefixKey, prefix)
}

func (s *stateDbKeys) BuildDbPrefix(height uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, height)

	return append(stateDbPrefixBytes, b...)
}

func (s *stateDbKeys) IdentityKey(addr common.Address) []byte {
	return append(identityPrefix, addr[:]...)
}

func (s *stateDbKeys) IdentityKeyToAddress(key []byte) common.Address {
	addr := common.Address{}
	addr.SetBytes(key[len(identityPrefix):])
	return addr
}

func (s *stateDbKeys) AddressKey(addr common.Address) []byte {
	return append(addressPrefix, addr[:]...)
}

func (s *stateDbKeys) AddressKeyToAddress(key []byte) common.Address {
	addr := common.Address{}
	addr.SetBytes(key[len(addressPrefix):])
	return addr
}

func (s *stateDbKeys) GlobalKey() []byte {
	return globalKey
}

func (s *stateDbKeys) StatusSwitchKey() []byte {
	return statusSwitchKey
}

func (s *stateDbKeys) DelegationSwitchKey() []byte {
	return delegationSwitchKey
}

func (s *stateDbKeys) DelayedOfflinePenaltyKey() []byte {
	return delayedOfflinePenaltyKey
}

func (s *stateDbKeys) ContractStoreKey(address common.Address, key []byte) []byte {
	return append(append(contractStorePrefix, address[:]...), key...)
}

func (s *stateDbKeys) BurntCoinsKey(height uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, height)
	return append(burntCoinsPrefix, b...)
}

func (s *stateDbKeys) BurntCoinsKeyToHeight(key []byte) uint64 {
	return binary.BigEndian.Uint64(key[len(burntCoinsPrefix):])
}

func (s *stateDbKeys) ContractCodeKey(codeHash common.Hash) []byte {
	return append(contractCodePrefix, codeHash[:]...)
}


type identityStateDbPrefix struct {
}

func (s *identityStateDbPrefix) LoadDbPrefix(db dbm.DB, preliminary bool) ([]byte, error) {
	key := currentIdentityStateDbPrefixKey
	if preliminary {
		key = preliminaryIdentityStateDbPrefixKey
	}
	p, err := db.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get value")
	}
	if p == nil {
		p = s.buildDbPrefix(0)
		b := db.NewBatch()
		s.SaveDbPrefix(b, p, preliminary)
		err := b.WriteSync()
		return p, errors.Wrapf(err, "failed to write value")
	}
	return p, nil
}

func (s *identityStateDbPrefix) buildDbPrefix(height uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, height)
	return append(identityStateDbPrefixBytes, b...)
}

func (s *identityStateDbPrefix) SaveDbPrefix(batch dbm.Batch, prefix []byte, preliminary bool) {
	key := currentIdentityStateDbPrefixKey
	if preliminary {
		key = preliminaryIdentityStateDbPrefixKey
	}
	batch.Set(key, prefix)
}

package dao

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type BearerTokenBloomFilter struct {
	filter *bloom.BloomFilter
}

func NewBearerTokenBloomFilter(filter *bloom.BloomFilter) *BearerTokenBloomFilter {
	return &BearerTokenBloomFilter{filter: filter}
}

func (bf *BearerTokenBloomFilter) AddKeys(keys ...string) {
	if bf.filter == nil {
		// disabled bloom filter, do nothing
		return
	}

	for _, key := range keys {
		bf.filter.AddString(key)
	}
}

func (bf *BearerTokenBloomFilter) CheckKey(key string) bool {
	if bf.filter == nil {
		// disabled bloom filter, always return true
		return true
	}

	return bf.filter.TestString(key)
}

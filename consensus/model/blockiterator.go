package model

import "github.com/zilong-dai/karlsen-miner/consensus/model/externalapi"

// BlockIterator is an iterator over blocks according to some order.
type BlockIterator interface {
	First() bool
	Next() bool
	Get() (*externalapi.DomainHash, error)
	Close() error
}

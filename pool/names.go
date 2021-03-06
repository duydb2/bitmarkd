// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pool

// Notes:
// 1. each separate pool has a single byte prefix (to spread the keys in LevelDB)
// 2. digest = sha256(sha256(data)) i.e. must be compatible with Bitcoin merkle tree
// 3. prop = registration-transfer-digest
//
// Blocks:
//
//   B<block-number>       - block store (already mined blocks) = header + cbLength + coinbase + count + merkle tree of transactions
//
// Transactions:
//
//   T<tx-digest>          - packed transaction data
//   S<tx-digest>          - state: byte[expired(E), unpaid(U), available(A), mined(M)] ++ int64[the U/A table count value]
//   U<count>              - transaction-digest ++ int64[timestamp] (pool of unpaid transactions for checking)
//   A<count>              - transaction-digest (pool of payment confirmed transactions, available for mining)
//
// Assets:
//
//   I<assetIndex>         - transaction-digest (to locate the AssetData transaction)
//
// Ownership:
//
//   O<bmtran-digest>      - owner public key ++ registration digest (to check current ownership of property)
//
//   K<pubkey><tx-digest>  - byte[asset(A), Bitmark Issue(I), bitmark transfer(T)] ++ asset index
//                           (to list current ownership of asset/issue/bitmark)
//
// Networking:
//
//   P<IP:port>            - P2P: ZMQ public-key
//   R<IP:port>            - RPC: certificate-fingerprint
//   C<fingerprint>        - raw certificate

// type for pool name
type nameb byte

// Names of the pools
const (
	// networking pools
	Peers        = nameb('P')
	RPCs         = nameb('R')
	Certificates = nameb('C')

	// transaction data pools
	TransactionData  = nameb('T')
	TransactionState = nameb('S')

	// transaction index pools
	UnpaidIndex    = nameb('U')
	AvailableIndex = nameb('A')

	// asset
	AssetData = nameb('I')

	// ownership indexes
	OwnerIndex = nameb('O')

	// blocks
	BlockData = nameb('B')

	// just for testing
	TestData = nameb('Z')
)

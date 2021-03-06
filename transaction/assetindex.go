// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package transaction

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/bitmark-inc/bitmarkd/fault"
)

// limits
const (
	assetIndexPrefix     = "BMA0"
	assetIndexPrefixSize = len(linkPrefix)
	AssetIndexSize       = 64
)

// the type for an iasset index - same as block digest
// stored as little endian byte array
// represented as big endian hex value for print
// represented as little endian base64 text for JSON encoding
type AssetIndex [AssetIndexSize]byte

// create an asset index from a byte slice
//
// double SHA2-512 Hash
func NewAssetIndex(record []byte) AssetIndex {
	roundOne := sha512.Sum512(record)
	roundTwo := sha512.Sum512(roundOne[:])
	return AssetIndex(roundTwo)
}

// internal function to return a reversed byte order copy of a asset index
func (d AssetIndex) reversed() []byte {
	result := make([]byte, AssetIndexSize)
	for i := 0; i < AssetIndexSize; i += 1 {
		result[i] = d[AssetIndexSize-1-i]
	}
	return result
}

// internal function to return a reversed byte order copy of a asset index
func (assetIndex *AssetIndex) reversedFromBytes(buffer []byte) {
	for i := 0; i < AssetIndexSize; i += 1 {
		assetIndex[i] = buffer[AssetIndexSize-1-i]
	}
}

// convert a binary assetIndex to byte slice
func (assetIndex AssetIndex) Bytes() []byte {
	return assetIndex[:]
}

// convert a binary assetIndex to big endian hex string for use by the fmt package (for %s)
func (assetIndex AssetIndex) String() string {
	return hex.EncodeToString(assetIndex.reversed())
}

// convert a binary assetIndex to big endian hex string for use by the fmt package (for %#v)
func (assetIndex AssetIndex) GoString() string {
	return "<asset:" + hex.EncodeToString(assetIndex.reversed()) + ">"
}

// convert a big endian hex text representation to a assetIndex for use by the format package scan routines
func (assetIndex *AssetIndex) Scan(state fmt.ScanState, verb rune) error {
	count := 0
	token, err := state.Token(true, func(c rune) bool {
		if count < assetIndexPrefixSize { // allow anything for prefix - exact check later
			count += 1
			return true
		}
		if c >= '0' && c <= '9' {
			return true
		}
		if c >= 'A' && c <= 'F' {
			return true
		}
		if c >= 'a' && c <= 'f' {
			return true
		}
		return false
	})
	if nil != err {
		return err
	}
	expectedSize := assetIndexPrefixSize + hex.EncodedLen(AssetIndexSize)
	actualSize := len(token)
	if expectedSize != actualSize {
		return fault.ErrNotAssetIndex
	}

	if assetIndexPrefix != string(token[:assetIndexPrefixSize]) {
		return fault.ErrNotAssetIndex
	}

	buffer := make([]byte, AssetIndexSize)
	byteCount, err := hex.Decode(buffer, token[assetIndexPrefixSize:])
	if nil != err {
		return err
	}

	if AssetIndexSize != byteCount {
		return fault.ErrNotAssetIndex
	}
	assetIndex.reversedFromBytes(buffer)
	return nil
}

// convert a binary assetIndex to little endian base64 text for JSON
//
// ***** possibly re-use MarshalText to save code duplication
// ***** but would that cost more buffer copying?
func (assetIndex AssetIndex) MarshalJSON() ([]byte, error) {
	// stage the prefix and assetIndex
	stageSize := assetIndexPrefixSize + len(assetIndex)
	stage := make([]byte, stageSize)
	copy(stage, assetIndexPrefix)

	for i := 0; i < AssetIndexSize; i += 1 {
		stage[assetIndexPrefixSize+i] = assetIndex[i]
	}

	// encode to base64
	size := 2 + base64.StdEncoding.EncodedLen(stageSize)
	buffer := make([]byte, size)
	buffer[0] = '"'
	buffer[size-1] = '"'
	base64.StdEncoding.Encode(buffer[1:], stage)
	return buffer, nil
}

// convert a little endian base64 string to a assetIndex for conversion from JSON
func (assetIndex *AssetIndex) UnmarshalJSON(s []byte) error {
	// length = '"' + characters + '"'
	last := len(s) - 1
	if '"' != s[0] || '"' != s[last] {
		return fault.ErrInvalidCharacter
	}
	return assetIndex.UnmarshalText(s[1:last])
}

// convert assetIndex to little endian base64 text
//
// ***** possibly use NewEncoder and byte buffer to save copy
func (assetIndex AssetIndex) MarshalText() ([]byte, error) {
	// stage the prefix and assetIndex
	stageSize := assetIndexPrefixSize + len(assetIndex)
	stage := make([]byte, stageSize)
	copy(stage, assetIndexPrefix)

	for i := 0; i < AssetIndexSize; i += 1 {
		stage[assetIndexPrefixSize+i] = assetIndex[i]
	}

	// encode to base64
	size := base64.StdEncoding.EncodedLen(stageSize)
	buffer := make([]byte, size)
	base64.StdEncoding.Encode(buffer, stage)
	return buffer, nil
}

// convert little endian base64 text into a assetIndex
func (assetIndex *AssetIndex) UnmarshalText(s []byte) error {
	buffer := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	byteCount, err := base64.StdEncoding.Decode(buffer, s)
	if nil != err {
		return err
	}
	if assetIndexPrefixSize+AssetIndexSize != byteCount {
		return fault.ErrNotAssetIndex
	}
	if assetIndexPrefix != string(buffer[:assetIndexPrefixSize]) {
		return fault.ErrNotAssetIndex
	}

	for i := 0; i < AssetIndexSize; i += 1 {
		assetIndex[i] = buffer[assetIndexPrefixSize+i]
	}
	return nil
}

// convert and validate little endian binary byte slice to a assetIndex
func AssetIndexFromBytes(assetIndex *AssetIndex, buffer []byte) error {
	if AssetIndexSize != len(buffer) {
		return fault.ErrNotAssetIndex
	}
	for i := 0; i < AssetIndexSize; i += 1 {
		assetIndex[i] = buffer[i]
	}
	return nil
}

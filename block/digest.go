// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package block

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/bitmark-inc/bitmarkd/fault"
	"math/big"
)

// number of bytes in the digest
const DigestSize = 32

// the type for a digest
// stored as little endian byte array
// represented as big endian hex value for print
// represented as little endian base64 text for JSON encoding
type Digest [DigestSize]byte

// create a digest from a byte slice
//
// BTC compatible double SHA2-256 Hash
func NewDigest(record []byte) Digest {
	roundOne := sha256.Sum256(record)
	roundTwo := sha256.Sum256(roundOne[:])
	return Digest(roundTwo)
}

// convert the hash to its equivalent big.Int
func (digest Digest) Cmp(difficulty *big.Int) int {
	bigEndian := reversed(digest)
	result := new(big.Int)
	return result.SetBytes(bigEndian[:]).Cmp(difficulty)
}

// internal function to return a reversed byte order copy of a digest
func reversed(d Digest) []byte {
	result := make([]byte, DigestSize)
	for i := 0; i < DigestSize; i += 1 {
		result[i] = d[DigestSize-1-i]
	}
	return result
}

// convert a binary digest to BTC little endian word swapped hex string for use by miners
func (digest Digest) BtcHex() string {
	l := len(digest)
	buffer := make([]byte, l)
	for i := 0; i < l; i += 4 {
		buffer[i+0] = digest[i+3]
		buffer[i+1] = digest[i+2]
		buffer[i+2] = digest[i+1]
		buffer[i+3] = digest[i+0]
	}
	return hex.EncodeToString(buffer)
}

// convert a binary digest to hex string for use by the stratum miner
//
// the stored version and the output string are both little endian
func (digest Digest) MinerHex() string {
	return hex.EncodeToString(digest[:])
}

// convert a binary digest to hex string for use by the fmt package (for %s)
//
// the stored version is in little endian, but the output string is big endian
func (digest Digest) String() string {
	return hex.EncodeToString(reversed(digest))
}

// convert a binary digest to big endian hex string for use by the fmt package (for %#v)
func (digest Digest) GoString() string {
	return "<sha256*2:" + hex.EncodeToString(reversed(digest)) + ">"
}

// convert a big endian hex representation to a digest for use by the format package scan routines
func (digest *Digest) Scan(state fmt.ScanState, verb rune) error {
	token, err := state.Token(true, func(c rune) bool {
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
	buffer := make([]byte, hex.DecodedLen(len(token)))
	byteCount, err := hex.Decode(buffer, token)
	if nil != err {
		return err
	}

	for i, v := range buffer[:byteCount] {
		digest[DigestSize-1-i] = v
	}
	return nil
}

// convert a binary digest to little endian base64 for JSON
func (digest Digest) MarshalJSON() ([]byte, error) {
	size := 2 + base64.StdEncoding.EncodedLen(len(digest))
	buffer := make([]byte, size)
	buffer[0] = '"'
	buffer[size-1] = '"'
	base64.StdEncoding.Encode(buffer[1:], digest[:])
	return buffer, nil
}

// convert a little endian base64 string to a digest for conversion from JSON
func (digest *Digest) UnmarshalJSON(s []byte) error {
	// length = '"' + characters + '"'
	last := len(s) - 1
	if '"' != s[0] || '"' != s[last] {
		return fault.ErrInvalidCharacter
	}
	return digest.UnmarshalText(s[1:last])
}

// convert digest to little endian base64 text
func (digest Digest) MarshalText() ([]byte, error) {
	size := base64.StdEncoding.EncodedLen(len(digest))
	buffer := make([]byte, size)
	base64.StdEncoding.Encode(buffer, digest[:])
	return buffer, nil
}

// convert little endian base 64 text into a digest
func (digest *Digest) UnmarshalText(s []byte) error {
	buffer := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	byteCount, err := base64.StdEncoding.Decode(buffer, s)
	if nil != err {
		return err
	}
	for i, v := range buffer[:byteCount] {
		digest[i] = v
	}
	return nil
}

// convert and validate little endian binary byte slice to a digest
func DigestFromBytes(digest *Digest, buffer []byte) error {
	if DigestSize != len(buffer) {
		return fault.ErrNotLink
	}
	for i := 0; i < DigestSize; i += 1 {
		digest[i] = buffer[i]
	}
	return nil
}

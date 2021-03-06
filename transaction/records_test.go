// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package transaction_test

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/agl/ed25519"
	"github.com/bitmark-inc/bitmarkd/fault"
	"github.com/bitmark-inc/bitmarkd/mode"
	"github.com/bitmark-inc/bitmarkd/transaction"
	"github.com/bitmark-inc/bitmarkd/util"
	"reflect"
	"testing"
)

// to print a keypair for future tests
func TestGenerateKeypair(t *testing.T) {
	generate := false

	// generate = true // (uncomment to get a new key pair)

	if generate {
		// display key pair and fail the test
		// use the displayed values to modify data below
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if nil != err {
			t.Errorf("key pair generation error: %v", err)
			return
		}
		t.Errorf("*** GENERATED:\n%s", formatBytes("publicKey", publicKey[:]))
		t.Errorf("*** GENERATED:\n%s", formatBytes("privateKey", privateKey[:]))
		return
	}
}

// to hold a keypair for testing
type keyPair struct {
	publicKey  [32]byte
	privateKey [64]byte
}

// public/private keys from above generate

var registrant = keyPair{
	publicKey: [...]byte{
		0x7a, 0x81, 0x92, 0x56, 0x5e, 0x6c, 0xa2, 0x35,
		0x80, 0xe1, 0x81, 0x59, 0xef, 0x30, 0x73, 0xf6,
		0xe2, 0xfb, 0x8e, 0x7e, 0x9d, 0x31, 0x49, 0x7e,
		0x79, 0xd7, 0x73, 0x1b, 0xa3, 0x74, 0x11, 0x01,
	},
	privateKey: [...]byte{
		0x66, 0xf5, 0x28, 0xd0, 0x2a, 0x64, 0x97, 0x3a,
		0x2d, 0xa6, 0x5d, 0xb0, 0x53, 0xea, 0xd0, 0xfd,
		0x94, 0xca, 0x93, 0xeb, 0x9f, 0x74, 0x02, 0x3e,
		0xbe, 0xdb, 0x2e, 0x57, 0xb2, 0x79, 0xfd, 0xf3,
		0x7a, 0x81, 0x92, 0x56, 0x5e, 0x6c, 0xa2, 0x35,
		0x80, 0xe1, 0x81, 0x59, 0xef, 0x30, 0x73, 0xf6,
		0xe2, 0xfb, 0x8e, 0x7e, 0x9d, 0x31, 0x49, 0x7e,
		0x79, 0xd7, 0x73, 0x1b, 0xa3, 0x74, 0x11, 0x01,
	},
}

var issuer = keyPair{
	publicKey: [...]byte{
		0x9f, 0xc4, 0x86, 0xa2, 0x53, 0x4f, 0x17, 0xe3,
		0x67, 0x07, 0xfa, 0x4b, 0x95, 0x3e, 0x3b, 0x34,
		0x00, 0xe2, 0x72, 0x9f, 0x65, 0x61, 0x16, 0xdd,
		0x7b, 0x01, 0x8d, 0xf3, 0x46, 0x98, 0xbd, 0xc2,
	},
	privateKey: [...]byte{
		0xf3, 0xf7, 0xa1, 0xfc, 0x33, 0x10, 0x71, 0xc2,
		0xb1, 0xcb, 0xbe, 0x4f, 0x3a, 0xee, 0x23, 0x5a,
		0xae, 0xcc, 0xd8, 0x5d, 0x2a, 0x80, 0x4c, 0x44,
		0xb5, 0xc6, 0x03, 0xb4, 0xca, 0x4d, 0x9e, 0xc0,
		0x9f, 0xc4, 0x86, 0xa2, 0x53, 0x4f, 0x17, 0xe3,
		0x67, 0x07, 0xfa, 0x4b, 0x95, 0x3e, 0x3b, 0x34,
		0x00, 0xe2, 0x72, 0x9f, 0x65, 0x61, 0x16, 0xdd,
		0x7b, 0x01, 0x8d, 0xf3, 0x46, 0x98, 0xbd, 0xc2,
	},
}

var ownerOne = keyPair{
	publicKey: [...]byte{
		0x27, 0x64, 0x0e, 0x4a, 0xab, 0x92, 0xd8, 0x7b,
		0x4a, 0x6a, 0x2f, 0x30, 0xb8, 0x81, 0xf4, 0x49,
		0x29, 0xf8, 0x66, 0x04, 0x3a, 0x84, 0x1c, 0x38,
		0x14, 0xb1, 0x66, 0xb8, 0x89, 0x44, 0xb0, 0x92,
	},
	privateKey: [...]byte{
		0xc7, 0xae, 0x9f, 0x22, 0x32, 0x0e, 0xda, 0x65,
		0x02, 0x89, 0xf2, 0x64, 0x7b, 0xc3, 0xa4, 0x4f,
		0xfa, 0xe0, 0x55, 0x79, 0xcb, 0x6a, 0x42, 0x20,
		0x90, 0xb4, 0x59, 0xb3, 0x17, 0xed, 0xf4, 0xa1,
		0x27, 0x64, 0x0e, 0x4a, 0xab, 0x92, 0xd8, 0x7b,
		0x4a, 0x6a, 0x2f, 0x30, 0xb8, 0x81, 0xf4, 0x49,
		0x29, 0xf8, 0x66, 0x04, 0x3a, 0x84, 0x1c, 0x38,
		0x14, 0xb1, 0x66, 0xb8, 0x89, 0x44, 0xb0, 0x92,
	},
}

var ownerTwo = keyPair{
	publicKey: [...]byte{
		0xa1, 0x36, 0x32, 0xd5, 0x42, 0x5a, 0xed, 0x3a,
		0x6b, 0x62, 0xe2, 0xbb, 0x6d, 0xe4, 0xc9, 0x59,
		0x48, 0x41, 0xc1, 0x5b, 0x70, 0x15, 0x69, 0xec,
		0x99, 0x99, 0xdc, 0x20, 0x1c, 0x35, 0xf7, 0xb3,
	},
	privateKey: [...]byte{
		0x8f, 0x83, 0x3e, 0x58, 0x30, 0xde, 0x63, 0x77,
		0x89, 0x4a, 0x8d, 0xf2, 0xd4, 0x4b, 0x17, 0x88,
		0x39, 0x1d, 0xcd, 0xb8, 0xfa, 0x57, 0x22, 0x73,
		0xd6, 0x2e, 0x9f, 0xcb, 0x37, 0x20, 0x2a, 0xb9,
		0xa1, 0x36, 0x32, 0xd5, 0x42, 0x5a, 0xed, 0x3a,
		0x6b, 0x62, 0xe2, 0xbb, 0x6d, 0xe4, 0xc9, 0x59,
		0x48, 0x41, 0xc1, 0x5b, 0x70, 0x15, 0x69, 0xec,
		0x99, 0x99, 0xdc, 0x20, 0x1c, 0x35, 0xf7, 0xb3,
	},
}

// helper to make an address
func makeAddress(publicKey *[32]byte) *transaction.Address {
	return &transaction.Address{
		AddressInterface: &transaction.ED25519Address{
			Test:      true,
			PublicKey: publicKey,
		},
	}
}

// test the packing/unpacking of registration record
//
// NOTE: This test change testing mode, so _MUST be before others
//
// ensures that pack->unpack returns the same original value
func TestPackAssetData(t *testing.T) {

	registrantAddress := makeAddress(&registrant.publicKey)

	r := transaction.AssetData{
		Description: "Just the description",
		Name:        "Item's Name",
		Fingerprint: "0123456789abcdef",
		Registrant:  registrantAddress,
	}

	expected := []byte{
		0x01, 0x14, 0x4a, 0x75, 0x73, 0x74, 0x20, 0x74,
		0x68, 0x65, 0x20, 0x64, 0x65, 0x73, 0x63, 0x72,
		0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x0b, 0x49,
		0x74, 0x65, 0x6d, 0x27, 0x73, 0x20, 0x4e, 0x61,
		0x6d, 0x65, 0x10, 0x30, 0x31, 0x32, 0x33, 0x34,
		0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63,
		0x64, 0x65, 0x66, 0x21, 0x13, 0x7a, 0x81, 0x92,
		0x56, 0x5e, 0x6c, 0xa2, 0x35, 0x80, 0xe1, 0x81,
		0x59, 0xef, 0x30, 0x73, 0xf6, 0xe2, 0xfb, 0x8e,
		0x7e, 0x9d, 0x31, 0x49, 0x7e, 0x79, 0xd7, 0x73,
		0x1b, 0xa3, 0x74, 0x11, 0x01,
	}

	expectedTxId := transaction.Link{
		0xd5, 0x1d, 0x38, 0xd5, 0x81, 0xaf, 0x18, 0xd4,
		0x47, 0xf5, 0x37, 0xd1, 0xb1, 0x54, 0x9c, 0x92,
		0x88, 0x52, 0xca, 0x1f, 0xdc, 0x09, 0x4f, 0x23,
		0x56, 0xc4, 0xf6, 0x9c, 0xe7, 0x17, 0x92, 0x3b,
	}

	expectedAssetIndex := transaction.AssetIndex{
		0x37, 0xde, 0x31, 0x95, 0xb8, 0x79, 0xcf, 0x56,
		0x92, 0xa5, 0x64, 0xe7, 0xf1, 0x9b, 0xd6, 0x20,
		0x3e, 0x3f, 0xd8, 0xcf, 0xfa, 0x80, 0x64, 0xeb,
		0xe8, 0x6e, 0xc8, 0xff, 0xf8, 0xdb, 0xf3, 0x26,
		0xff, 0xa4, 0xd6, 0x5c, 0xfe, 0xf8, 0x35, 0x0d,
		0xb4, 0xd7, 0x2d, 0x93, 0x40, 0x6f, 0x54, 0xfa,
		0x0d, 0x06, 0xce, 0x98, 0x00, 0x5a, 0x93, 0x99,
		0x95, 0xed, 0x05, 0xcc, 0x34, 0xfb, 0x73, 0x44,
	}

	// manually sign the record and attach signature to "expected"
	signature := ed25519.Sign(&registrant.privateKey, expected)
	r.Signature = signature[:]
	//t.Logf("signature: %#v", r.Signature)
	l := util.ToVarint64(uint64(len(signature)))
	expected = append(expected, l...)
	expected = append(expected, signature[:]...)

	// test the packer
	packed, err := r.Pack(registrantAddress)
	if nil != err {
		t.Errorf("pack error: %v", err)
	}

	// if either of above fail we will have the message _without_ a signature
	if !bytes.Equal(packed, expected) {
		t.Errorf("pack record: %x  expected: %x", packed, expected)
		t.Errorf("*** GENERATED Packed:\n%s", formatBytes("expected", packed))
		return
	}

	// check txIds
	txId := packed.MakeLink()

	if txId != expectedTxId {
		t.Errorf("pack tx id: %#v  expected: %#v", txId, expectedTxId)
		t.Errorf("*** GENERATED tx id:\n%s", formatBytes("expectedTxId", txId.Bytes()))
	}

	// check asset index
	assetIndex := r.AssetIndex()

	if assetIndex != expectedAssetIndex {
		t.Errorf("pack asset index: %#v  expected: %#v", assetIndex, expectedAssetIndex)
		t.Errorf("*** GENERATED asset index:\n%s", formatBytes("expectedAssetIndex", assetIndex.Bytes()))
	}

	// check test-network detection
	mode.SetTesting(false) // default so will not panic
	if _, err := packed.Unpack(); err != fault.ErrWrongNetworkForPublicKey {
		t.Errorf("expected 'wrong network for public key' but got error: %v", err)
	}
	mode.SetTesting(true) // enter test mode - ONLY ALLOWED ONCE (or panic will occur

	// test the unpacker
	unpacked, err := packed.Unpack()
	if nil != err {
		t.Errorf("unpack error: %v", err)
		return
	}

	reg, ok := unpacked.(*transaction.AssetData)
	if !ok {
		t.Errorf("did not unpack to AssetData")
		return
	}

	// display a JSON version for information
	item := struct {
		HexTxId   string
		TxId      transaction.Link
		HexAsset  string
		Asset     transaction.AssetIndex
		AssetData *transaction.AssetData
	}{
		HexTxId:   txId.String(),
		TxId:      txId,
		HexAsset:  assetIndex.String(),
		Asset:     assetIndex,
		AssetData: reg,
	}
	b, err := json.MarshalIndent(item, "", "  ")
	if nil != err {
		t.Errorf("json error: %v", err)
		return
	}

	t.Logf("AssetData: JSON: %s", b)

	// check that structure is preserved through Pack/Unpack
	// note reg is a pointer here
	if !reflect.DeepEqual(r, *reg) {
		t.Errorf("different, original: %v  recovered: %v", r, *reg)
		return
	}
}

// test the packing/unpacking of Bitmark issue record
//
// ensures that pack->unpack returns the same original value
func TestPackBitmarkIssue(t *testing.T) {

	issuerAddress := makeAddress(&issuer.publicKey)

	var asset transaction.AssetIndex
	_, err := fmt.Sscan("BMA04473fb34cc05ed9599935a0098ce060dfa546f40932dd7b40d35f8fe5cd6a4ff26f3dbf8ffc86ee8eb6480facfd83f3e20d69bf1e764a59256cf79b89531de37", &asset)
	if nil != err {
		t.Errorf("hex to link error: %v", err)
		return
	}

	r := transaction.BitmarkIssue{
		AssetIndex: asset,
		Owner:      issuerAddress,
		Nonce:      99,
	}

	expected := []byte{
		0x02, 0x40, 0x37, 0xde, 0x31, 0x95, 0xb8, 0x79,
		0xcf, 0x56, 0x92, 0xa5, 0x64, 0xe7, 0xf1, 0x9b,
		0xd6, 0x20, 0x3e, 0x3f, 0xd8, 0xcf, 0xfa, 0x80,
		0x64, 0xeb, 0xe8, 0x6e, 0xc8, 0xff, 0xf8, 0xdb,
		0xf3, 0x26, 0xff, 0xa4, 0xd6, 0x5c, 0xfe, 0xf8,
		0x35, 0x0d, 0xb4, 0xd7, 0x2d, 0x93, 0x40, 0x6f,
		0x54, 0xfa, 0x0d, 0x06, 0xce, 0x98, 0x00, 0x5a,
		0x93, 0x99, 0x95, 0xed, 0x05, 0xcc, 0x34, 0xfb,
		0x73, 0x44, 0x21, 0x13, 0x9f, 0xc4, 0x86, 0xa2,
		0x53, 0x4f, 0x17, 0xe3, 0x67, 0x07, 0xfa, 0x4b,
		0x95, 0x3e, 0x3b, 0x34, 0x00, 0xe2, 0x72, 0x9f,
		0x65, 0x61, 0x16, 0xdd, 0x7b, 0x01, 0x8d, 0xf3,
		0x46, 0x98, 0xbd, 0xc2, 0x63,
	}

	expectedTxId := transaction.Link{
		0xf6, 0x97, 0x48, 0x5e, 0xe8, 0xdd, 0xd7, 0x6f,
		0x8a, 0x3e, 0xf8, 0xb2, 0xac, 0x4b, 0x3f, 0xc7,
		0xfa, 0x77, 0xe4, 0xee, 0x1c, 0xb3, 0x18, 0x53,
		0x50, 0xc2, 0xa4, 0x72, 0x78, 0xb2, 0x04, 0x75,
	}

	// manually sign the record and attach signature to "expected"
	signature := ed25519.Sign(&issuer.privateKey, expected)
	r.Signature = signature[:]
	l := util.ToVarint64(uint64(len(signature)))
	expected = append(expected, l...)
	expected = append(expected, signature[:]...)

	// test the packer
	packed, err := r.Pack(issuerAddress)
	if nil != err {
		t.Errorf("pack error: %v", err)
	}

	// if either of above fail we will have the message _without_ a signature
	if !bytes.Equal(packed, expected) {
		t.Errorf("pack record: %x  expected: %x", packed, expected)
		t.Errorf("*** GENERATED Packed:\n%s", formatBytes("expected", packed))
		return
	}

	// check txId
	txId := packed.MakeLink()

	if txId != expectedTxId {
		t.Errorf("pack tx id: %#v  expected: %x", txId, expectedTxId)
		t.Errorf("*** GENERATED tx id:\n%s", formatBytes("expectedTxId", txId.Bytes()))
		return
	}

	// test the unpacker
	unpacked, err := packed.Unpack()
	if nil != err {
		t.Errorf("unpack error: %v", err)
		return
	}

	bmt, ok := unpacked.(*transaction.BitmarkIssue)
	if !ok {
		t.Errorf("did not unpack to BitmarkIssue")
		return
	}

	// display a JSON version for information
	item := struct {
		HexTxId      string
		TxId         transaction.Link
		BitmarkIssue *transaction.BitmarkIssue
	}{
		txId.String(),
		txId,
		bmt,
	}
	b, err := json.MarshalIndent(item, "", "  ")
	if nil != err {
		t.Errorf("json error: %v", err)
		return
	}

	t.Logf("Bitmark Issue: JSON: %s", b)

	// check that structure is preserved through Pack/Unpack
	// note reg is a pointer here
	if !reflect.DeepEqual(r, *bmt) {
		t.Errorf("different, original: %v  recovered: %v", r, *bmt)
		return
	}
}

// test the packing/unpacking of Bitmark transfer record
//
// transfer from issue
// ensures that pack->unpack returns the same original value
func TestPackBitmarkTransferOne(t *testing.T) {

	issuerAddress := makeAddress(&issuer.publicKey)
	ownerOneAddress := makeAddress(&ownerOne.publicKey)

	var link transaction.Link
	_, err := fmt.Sscan("BMK07504b27872a4c2505318b31ceee477fac73f4bacb2f83e8a6fd7dde85e4897f6", &link)
	if nil != err {
		t.Errorf("hex to link error: %v", err)
		return
	}

	r := transaction.BitmarkTransfer{
		Link:  link,
		Owner: ownerOneAddress,
	}

	expected := []byte{
		0x03, 0x20, 0xf6, 0x97, 0x48, 0x5e, 0xe8, 0xdd,
		0xd7, 0x6f, 0x8a, 0x3e, 0xf8, 0xb2, 0xac, 0x4b,
		0x3f, 0xc7, 0xfa, 0x77, 0xe4, 0xee, 0x1c, 0xb3,
		0x18, 0x53, 0x50, 0xc2, 0xa4, 0x72, 0x78, 0xb2,
		0x04, 0x75, 0x21, 0x13, 0x27, 0x64, 0x0e, 0x4a,
		0xab, 0x92, 0xd8, 0x7b, 0x4a, 0x6a, 0x2f, 0x30,
		0xb8, 0x81, 0xf4, 0x49, 0x29, 0xf8, 0x66, 0x04,
		0x3a, 0x84, 0x1c, 0x38, 0x14, 0xb1, 0x66, 0xb8,
		0x89, 0x44, 0xb0, 0x92,
	}

	expectedTxId := transaction.Link{
		0x01, 0xc9, 0x12, 0xe3, 0xa1, 0xe5, 0x0e, 0xde,
		0x0b, 0xc2, 0xaf, 0x34, 0x86, 0x81, 0x10, 0x64,
		0x8d, 0x57, 0x75, 0xc2, 0xe4, 0xff, 0xab, 0x6b,
		0xc5, 0xc6, 0x76, 0x4c, 0x5e, 0x86, 0xfd, 0x0b,
	}

	// manually sign the record and attach signature to "expected"
	signature := ed25519.Sign(&issuer.privateKey, expected)
	r.Signature = signature[:]
	l := util.ToVarint64(uint64(len(signature)))
	expected = append(expected, l...)
	expected = append(expected, signature[:]...)

	// test the packer
	packed, err := r.Pack(issuerAddress)
	if nil != err {
		t.Errorf("pack error: %v", err)
	}

	// if either of above fail we will have the message _without_ a signature
	if !bytes.Equal(packed, expected) {
		t.Errorf("pack record: %x  expected: %x", packed, expected)
		t.Errorf("*** GENERATED Packed:\n%s", formatBytes("expected", packed))
		return
	}

	// check txId
	txId := packed.MakeLink()

	if txId != expectedTxId {
		t.Errorf("pack txId: %#v  expected: %x", txId, expectedTxId)
		t.Errorf("*** GENERATED txId:\n%s", formatBytes("expectedTxId", txId.Bytes()))
		return
	}

	// test the unpacker
	unpacked, err := packed.Unpack()
	if nil != err {
		t.Errorf("unpack error: %v", err)
		return
	}

	bmt, ok := unpacked.(*transaction.BitmarkTransfer)
	if !ok {
		t.Errorf("did not unpack to BitmarkTransfer")
		return
	}

	// display a JSON version for information
	item := struct {
		HexTxId         string
		TxId            transaction.Link
		BitmarkTransfer *transaction.BitmarkTransfer
	}{
		txId.String(),
		txId,
		bmt,
	}
	b, err := json.MarshalIndent(item, "", "  ")
	if nil != err {
		t.Errorf("json error: %v", err)
		return
	}

	t.Logf("Bitmark Transfer: JSON: %s", b)

	// check that structure is preserved through Pack/Unpack
	// note reg is a pointer here
	if !reflect.DeepEqual(r, *bmt) {
		t.Errorf("different, original: %v  recovered: %v", r, *bmt)
		return
	}
}

// test the packing/unpacking of Bitmark transfer record
//
// test transfer to transfer
// ensures that pack->unpack returns the same original value
func TestPackBitmarkTransferTwo(t *testing.T) {

	ownerOneAddress := makeAddress(&ownerOne.publicKey)
	ownerTwoAddress := makeAddress(&ownerTwo.publicKey)

	var link transaction.Link
	_, err := fmt.Sscan("BMK00bfd865e4c76c6c56babffe4c275578d6410818634afc20bde0ee5a1e312c901", &link)
	if nil != err {
		t.Errorf("hex to link error: %v", err)
		return
	}

	r := transaction.BitmarkTransfer{
		Link:  link,
		Owner: ownerTwoAddress,
	}

	expected := []byte{
		0x03, 0x20, 0x01, 0xc9, 0x12, 0xe3, 0xa1, 0xe5,
		0x0e, 0xde, 0x0b, 0xc2, 0xaf, 0x34, 0x86, 0x81,
		0x10, 0x64, 0x8d, 0x57, 0x75, 0xc2, 0xe4, 0xff,
		0xab, 0x6b, 0xc5, 0xc6, 0x76, 0x4c, 0x5e, 0x86,
		0xfd, 0x0b, 0x21, 0x13, 0xa1, 0x36, 0x32, 0xd5,
		0x42, 0x5a, 0xed, 0x3a, 0x6b, 0x62, 0xe2, 0xbb,
		0x6d, 0xe4, 0xc9, 0x59, 0x48, 0x41, 0xc1, 0x5b,
		0x70, 0x15, 0x69, 0xec, 0x99, 0x99, 0xdc, 0x20,
		0x1c, 0x35, 0xf7, 0xb3,
	}

	expectedTxId := transaction.Link{
		0x6c, 0x5c, 0x0e, 0x43, 0xf5, 0x98, 0x06, 0xe5,
		0x19, 0x79, 0x3d, 0xd6, 0x16, 0x48, 0x73, 0x18,
		0x43, 0x06, 0x8a, 0x00, 0x93, 0xdb, 0x7b, 0x07,
		0x76, 0x6f, 0x7f, 0x7f, 0x2d, 0x18, 0xb9, 0x19,
	}

	// manually sign the record and attach signature to "expected"
	signature := ed25519.Sign(&ownerOne.privateKey, expected)
	r.Signature = signature[:]
	l := util.ToVarint64(uint64(len(signature)))
	expected = append(expected, l...)
	expected = append(expected, signature[:]...)

	// test the packer
	packed, err := r.Pack(ownerOneAddress)
	if nil != err {
		t.Errorf("pack error: %v", err)
	}

	// if either of above fail we will have the message _without_ a signature
	if !bytes.Equal(packed, expected) {
		t.Errorf("pack record: %x  expected: %x", packed, expected)
		t.Errorf("*** GENERATED Packed:\n%s", formatBytes("expected", packed))
		return
	}

	// check txId
	txId := packed.MakeLink()

	if txId != expectedTxId {
		t.Errorf("pack txId: %#v  expected: %x", txId, expectedTxId)
		t.Errorf("*** GENERATED txId:\n%s", formatBytes("expectedTxId", txId.Bytes()))
		return
	}

	// test the unpacker
	unpacked, err := packed.Unpack()
	if nil != err {
		t.Errorf("unpack error: %v", err)
		return
	}

	bmt, ok := unpacked.(*transaction.BitmarkTransfer)
	if !ok {
		t.Errorf("did not unpack to BitmarkTransfer")
		return
	}

	// display a JSON version for information
	item := struct {
		HexTxId         string
		TxId            transaction.Link
		BitmarkTransfer *transaction.BitmarkTransfer
	}{
		txId.String(),
		txId,
		bmt,
	}
	b, err := json.MarshalIndent(item, "", "  ")
	if nil != err {
		t.Errorf("json error: %v", err)
		return
	}

	t.Logf("Bitmark Transfer: JSON: %s", b)

	// check that structure is preserved through Pack/Unpack
	// note reg is a pointer here
	if !reflect.DeepEqual(r, *bmt) {
		t.Errorf("different, original: %v  recovered: %v", r, *bmt)
		return
	}
}

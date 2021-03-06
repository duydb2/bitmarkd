// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// error instances
//
// Provides a single instance of errors to allow easy comparison
package fault

// error base
type GenericError string

// to allow for different classes of errors
type ExistsError GenericError
type InvalidError GenericError
type LengthError GenericError
type NotFoundError GenericError
type ProcessError GenericError
type RecordError GenericError

// common errors - keep in alphabetic order
var (
	ErrAlreadyInitialised            = ExistsError("already initialised")
	ErrAssetNotFound                 = NotFoundError("asset not found")
	ErrBlockNotFound                 = NotFoundError("block not found")
	ErrCannotDecodeAddress           = RecordError("cannot decode address")
	ErrCertificateFileAlreadyExists  = ExistsError("certificate file already exists")
	ErrCertificateNotFound           = NotFoundError("certificate not found")
	ErrChecksumMismatch              = ProcessError("checksum mismatch")
	ErrCountMismatch                 = ProcessError("count mismatch")
	ErrConnectingToSelfForbidden     = ProcessError("connecting to self forbidden")
	ErrDescriptionTooLong            = LengthError("name too long")
	ErrFingerprintTooLong            = LengthError("fingerprint too long")
	ErrInsufficientPayment           = InvalidError("insufficient payment")
	ErrInvalidBlock                  = InvalidError("invalid block")
	ErrInvalidBlockHeader            = InvalidError("invalid block header")
	ErrInvalidCoinbase               = InvalidError("invalid coinbase")
	ErrInvalidCount                  = InvalidError("invalid count")
	ErrInvalidCharacter              = InvalidError("invalid character")
	ErrInvalidCurrency               = InvalidError("invalid currency")
	ErrInvalidIPAddress              = InvalidError("invalid IP Address")
	ErrInvalidKeyLength              = InvalidError("invalid key length")
	ErrInvalidKeyType                = InvalidError("invalid key type")
	ErrInvalidLength                 = InvalidError("invalid length")
	ErrInvalidLoggerChannel          = InvalidError("invalid logger channel")
	ErrInvalidPortNumber             = InvalidError("invalid port number")
	ErrInvalidRemote                 = InvalidError("invalid remote: expected 'z85',IP:Port")
	ErrInvalidSignature              = InvalidError("invalid signature")
	ErrInvalidTransactionChain       = InvalidError("invalid transaction chain")
	ErrInvalidType                   = InvalidError("invalid type")
	ErrInvalidVersion                = InvalidError("invalid version")
	ErrKeyFileAlreadyExists          = ExistsError("key file already exists")
	ErrKeyFileNotFound               = NotFoundError("key file not found")
	ErrKeyNotFound                   = NotFoundError("key not found")
	ErrLinkNotFound                  = NotFoundError("link not found")
	ErrLinksToUnconfirmedTransaction = InvalidError("links to unconfirmed transaction")
	ErrMessagingTerminated           = ProcessError("messaging terminated")
	ErrNameTooLong                   = LengthError("name too long")
	ErrNoPaymentToMiner              = InvalidError("no payment to miner")
	ErrNotABitmarkPayment            = InvalidError("not a bitmark payment")
	ErrNotAssetIndex                 = RecordError("not asset index")
	ErrNotCurrentOwner               = RecordError("not current owner")
	ErrNotInitialised                = NotFoundError("not initialised")
	ErrNotLink                       = RecordError("not link")
	ErrNotPublicKey                  = RecordError("not public key")
	ErrNotTransactionType            = RecordError("not transaction type")
	ErrNotTransactionPack            = RecordError("not transaction pack")
	ErrPaymentAddressMissing         = NotFoundError("payment address missing")
	ErrPeerAlreadyExists             = ExistsError("peer already exists")
	ErrPeerNotFound                  = NotFoundError("peer not found")
	ErrSignatureTooLong              = LengthError("signature too long")
	ErrTransactionAlreadyExists      = ExistsError("transaction already exists")
	ErrWrongNetworkForPublicKey      = InvalidError("wrong network for public key")
)

// the error interface base method
func (e GenericError) Error() string { return string(e) }

// the error interface methods
func (e ExistsError) Error() string   { return string(e) }
func (e InvalidError) Error() string  { return string(e) }
func (e LengthError) Error() string   { return string(e) }
func (e NotFoundError) Error() string { return string(e) }
func (e ProcessError) Error() string  { return string(e) }
func (e RecordError) Error() string   { return string(e) }

// determine the class of an error
func IsErrExists(e error) bool   { _, ok := e.(ExistsError); return ok }
func IsErrInvalid(e error) bool  { _, ok := e.(InvalidError); return ok }
func IsErrLength(e error) bool   { _, ok := e.(LengthError); return ok }
func IsErrNotFound(e error) bool { _, ok := e.(NotFoundError); return ok }
func IsErrProcess(e error) bool  { _, ok := e.(ProcessError); return ok }
func IsErrRecord(e error) bool   { _, ok := e.(RecordError); return ok }

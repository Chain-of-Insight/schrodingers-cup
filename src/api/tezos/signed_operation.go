package tezos

import (
	"crypto"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/xerrors"

	"github.com/anchorageoss/tezosprotocol/v2"
)

func VerifyOperation(message []byte, signature tezosprotocol.Signature, publicKey crypto.PublicKey) error {
	return verifyGeneric(tezosprotocol.OperationWatermark, message, signature, publicKey)
}

func verifyGeneric(watermark tezosprotocol.Watermark, message []byte, signature tezosprotocol.Signature, publicKey crypto.PublicKey) error {
	// prepend the tezos operation watermark
	byteesWithWatermark := append([]byte{byte(watermark)}, message...)

	// hash
	payloadHash := blake2b.Sum256(byteesWithWatermark)

	// verify signature over hash
	sigPrefix, sigBytes, err := tezosprotocol.Base58CheckDecode(string(signature))
	if err != nil {
		return xerrors.Errorf("failed to decode signature: %s: %w", signature, err)
	}
	var ok bool
	switch key := publicKey.(type) {
	case ed25519.PublicKey:
		if sigPrefix != tezosprotocol.PrefixEd25519Signature && sigPrefix != tezosprotocol.PrefixGenericSignature {
			return xerrors.Errorf("signature type %s does not match public key type %T", sigPrefix, publicKey)
		}
		ok = ed25519.Verify(key, payloadHash[:], sigBytes)
	default:
		return xerrors.Errorf("unsupported public key type: %T", publicKey)
	}
	if !ok {
		return xerrors.Errorf("invalid signature %s for public key %s", signature, publicKey)
	}
	return nil
}

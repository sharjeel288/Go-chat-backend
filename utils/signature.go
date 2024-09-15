package utils

import (
	"ChaiLabs/config"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// VerifySignature verifies a message signature and returns a boolean indicating whether it is valid.
func VerifySignature(signature, address, message string) bool {
	signatureBytes := common.FromHex(signature)
	messageHash := crypto.Keccak256Hash([]byte(message))
	signer, err := crypto.SigToPub(messageHash.Bytes(), signatureBytes)
	if err != nil {
		log.Fatal(err)
	}
	recoveredAddress := crypto.PubkeyToAddress(*signer)
	return strings.EqualFold(recoveredAddress.Hex(), address)
}

// CreateSignature creates a signature for a message and returns the signature and the address of the signer.
func CreateSignature(message string) (string, string, error) {
	privateKeyECDSA, err := crypto.HexToECDSA(config.PrivateKey())
	if err != nil {
		return "", "", err
	}

	signatureBytes, err := crypto.Sign(crypto.Keccak256Hash([]byte(message)).Bytes(), privateKeyECDSA)
	if err != nil {
		return "", "", err
	}

	signature := hex.EncodeToString(signatureBytes)
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()
	return signature, address, nil
}

// CreateMessageForSigning creates a message that users will sign to verify ownership of their address.
func CreateMessageForSigning(address string) string {
	message := fmt.Sprintf(
		"%s\nPlease sign this message to verify ownership of your address %s and login to Chailab.\n%d",
		config.AuthMessage(), address, time.Now().Unix(),
	)
	return message
}

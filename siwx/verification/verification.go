package verification

import (
    "errors"
    "fmt"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/nnn-community/go-siwx/siwx/db"
    "github.com/nnn-community/go-siwx/siwx/session"
    "github.com/nnn-community/go-siwx/siwx/utils"
    "strings"
)

type Verification struct {
    Success bool
    Result  *session.Siwx
}

func VerifyUser(message string, signature string) (*Verification, error) {
    if message == "" {
        return &Verification{Success: false}, errors.New("SiwxMessage is undefined")
    }

    address, err := utils.GetAddressFromMessage(message)
    chainId, err := utils.GetChainIdFromMessage(message)

    if err != nil {
        return &Verification{Success: false}, errors.New("failed to extract address or chain ID")
    }

    // Handle CAIP-2 formatted chain ID
    if strings.Contains(chainId, ":") {
        parts := strings.Split(chainId, ":")
        chainId = parts[1]
    }

    isValid, err := VerifyMessage(message, address, signature)

    if err != nil || !isValid {
        return &Verification{Success: false}, errors.New("invalid signature")
    }

    var userId string

    err = db.Connection.Raw(`
        SELECT user_id
		FROM user_addresses
		WHERE address = ?
        LIMIT 1
    `, strings.ToLower(address)).First(&userId).Error

    if err != nil || userId == "" {
        var groupID string

        db.Connection.Raw(`
            SELECT id
            FROM groups
            WHERE "default" = TRUE
            LIMIT 1
        `).First(&groupID)

        db.Connection.Raw(`
            INSERT INTO users (id, group_id)
            VALUES (gen_random_uuid(), ?)
            RETURNING id
        `, groupID).Scan(&userId)

        db.Connection.Exec(`
            INSERT INTO user_addresses (id, user_id, address, master)
            VALUES (gen_random_uuid(), ?, ?, TRUE)
        `, userId, strings.ToLower(address))
    }

    return &Verification{
        Success: true,
        Result: &session.Siwx{
            Address: address,
            ChainId: chainId,
        },
    }, nil
}

func VerifyMessage(message, address, signature string) (bool, error) {
    // 1. Recreate the prefixed hash of the message (EIP-191)
    msg := []byte(message)
    prefixedMsg := crypto.Keccak256(
        []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)),
    )

    // 2. Decode the signature hex (assumed 0x-prefixed)
    sigBytes, err := hexToBytes(signature)

    if err != nil {
        return false, err
    }

    if len(sigBytes) != 65 {
        return false, errors.New("signature must be 65 bytes")
    }

    // 3. Adjust the recovery id (V) if needed
    // Ethereum signatures have v as 27 or 28; go-ethereum expects 0 or 1
    if sigBytes[64] >= 27 {
        sigBytes[64] -= 27
    }

    // 4. Recover the public key from the signature and hashed message
    pubKey, err := crypto.SigToPub(prefixedMsg, sigBytes)

    if err != nil {
        return false, err
    }

    // 5. Compute the address from the public key
    recoveredAddr := crypto.PubkeyToAddress(*pubKey)

    // 6. Compare recovered address with expected address
    return common.HexToAddress(address) == recoveredAddr, nil
}

func hexToBytes(str string) ([]byte, error) {
    if len(str) >= 2 && str[0:2] == "0x" {
        str = str[2:]
    }

    return common.FromHex("0x" + str), nil
}

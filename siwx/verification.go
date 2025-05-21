package siwx

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "strings"
)

type Verification struct {
    Success bool
    Result  *VerificationResult
}

type VerificationResult struct {
    Address string
    ChainID string
}

func VerifyUser(db *sql.DB, message string, signature string) (*Verification, error) {
    ctx := context.Background()

    if message == "" {
        return &Verification{Success: false}, errors.New("SiwxMessage is undefined")
    }

    address, err := GetAddressFromMessage(message)
    if err != nil {
        return &Verification{Success: false}, fmt.Errorf("failed to extract address: %w", err)
    }

    chainId, err := GetChainIdFromMessage(message)
    if err != nil {
        return &Verification{Success: false}, fmt.Errorf("failed to extract chain ID: %w", err)
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

    addressLower := strings.ToLower(address)

    // Check if the address already exists with a user
    var userID sql.NullInt64
    query := `SELECT user_id FROM user_addresses WHERE address = $1 LIMIT 1`
    err = db.QueryRowContext(ctx, query, addressLower).Scan(&userID)
    if err != nil && err != sql.ErrNoRows {
        return &Verification{Success: false}, fmt.Errorf("DB query error: %w", err)
    }

    if !userID.Valid {
        // Get the default group ID
        var groupID int
        err := db.QueryRowContext(ctx, `SELECT id FROM groups WHERE is_default = true LIMIT 1`).Scan(&groupID)
        if err != nil {
            return &Verification{Success: false}, fmt.Errorf("could not find default group: %w", err)
        }

        // Create user and address (in a transaction)
        tx, err := db.BeginTx(ctx, nil)
        if err != nil {
            return &Verification{Success: false}, fmt.Errorf("failed to begin transaction: %w", err)
        }
        defer tx.Rollback()

        var newUserID int
        err = tx.QueryRowContext(ctx,
            `INSERT INTO users (active, group_id) VALUES (true, $1) RETURNING id`,
            groupID,
        ).Scan(&newUserID)
        if err != nil {
            return &Verification{Success: false}, fmt.Errorf("failed to insert user: %w", err)
        }

        _, err = tx.ExecContext(ctx,
            `INSERT INTO user_addresses (address, master, user_id) VALUES ($1, true, $2)`,
            addressLower, newUserID,
        )
        if err != nil {
            return &Verification{Success: false}, fmt.Errorf("failed to insert address: %w", err)
        }

        if err = tx.Commit(); err != nil {
            return &Verification{Success: false}, fmt.Errorf("transaction commit failed: %w", err)
        }
    }

    return &Verification{
        Success: true,
        Result: &VerificationResult{
            Address: address,
            ChainID: chainId,
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
    if common.HexToAddress(address) != recoveredAddr {
        return false, nil // signature invalid
    }

    return true, nil
}

func hexToBytes(str string) ([]byte, error) {
    if len(str) >= 2 && str[0:2] == "0x" {
        str = str[2:]
    }

    return common.FromHex("0x" + str), nil
}

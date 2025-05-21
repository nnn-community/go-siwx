package siwx

import (
    "errors"
    "regexp"
    "strings"
)

func GetAddressFromMessage(message string) (string, error) {
    re := regexp.MustCompile(`0x[a-fA-F0-9]{40}`)
    match := re.FindString(message)

    if match == "" {
        return "", errors.New("invalid address")
    }

    return strings.ToLower(match), nil
}

func GetChainIdFromMessage(message string) (string, error) {
    re := regexp.MustCompile(`Chain ID:\s*(\d+)`)
    match := re.FindStringSubmatch(message)

    if len(match) < 2 {
        return "eip155:1", nil
    }

    return "eip155:" + match[1], nil
}

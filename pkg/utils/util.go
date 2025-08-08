package utils

import (
	"fmt"
	"math/big"
	"strings"
)

func HexToBigInt(hexStr string) (*big.Int, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x") // "0x" 제거
	n := new(big.Int)
	_, ok := n.SetString(hexStr, 16) // 16진수로 파싱
	if !ok {
		return nil, fmt.Errorf("invalid hex string: %s", hexStr)
	}
	return n, nil
}

func DivideBy(a *big.Int, b *big.Int) string {
	bf := new(big.Float).SetInt(a) // balance → big.Float
	denom := new(big.Float).SetInt(b)
	return new(big.Float).Quo(bf, denom).String()
}

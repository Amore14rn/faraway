package pow

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"

	"github.com/Amore14rn/faraway/app/internal/pkg/config"
)

// const zeroByte = '0'

type ProofOfWork struct {
	Version    int
	ZerosCount int
	Date       int64
	Resource   string
	Rand       string
	Counter    int
}

func (p ProofOfWork) Format() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", p.Version, p.ZerosCount, p.Date, p.Resource, p.Rand, p.Counter)
}

func CalculateSHA1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func ValidateHash(hash string, zerosCount int) bool {
	requiredPrefix := strings.Repeat("0", zerosCount)
	return strings.HasPrefix(hash, requiredPrefix)
}

func (p ProofOfWork) ComputeProofOfWork(maxIterations int) (ProofOfWork, error) {
	cfg := config.GetConfig()
	if maxIterations < 0 {
		return p, fmt.Errorf("max iterations cannot be negative")
	}

	startTime := time.Now()
	for p.Counter <= maxIterations || maxIterations <= 0 {
		header := p.Format()
		hash := CalculateSHA1Hash(header)
		if ValidateHash(hash, p.ZerosCount) {
			return p, nil
		}
		p.Counter++

		if time.Since(startTime).Seconds() > float64(cfg.HashCash.ChallengeLifetime) {
			break
		}
	}
	return p, fmt.Errorf("max iterations exceeded")
}

package otp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/bits"
	"strconv"
	"testing"
	"time"
)

type otpComponents struct {
	seed      string
	hash      [32]byte
	truncated uint32
	x1        int64
	otp       int64
}

func calcComponents(username string, t time.Time) otpComponents {
	lcmOnce.Do(loadLCMParams)

	interval := t.Unix() / 30
	seed := fmt.Sprintf("%s-%d", username, interval)
	hash := sha256.Sum256([]byte(seed))

	hashString := hex.EncodeToString(hash[:])
	truncatedHash := hashString[:8]
	x0, _ := strconv.ParseInt(truncatedHash, 16, 64)
	x1 := (lcmA*x0 + lcmC) % lcmM
	otpVal := x1 % 1000000

	return otpComponents{
		seed:      seed,
		hash:      hash,
		truncated: uint32(x0),
		x1:        x1,
		otp:       otpVal,
	}
}

func hammingDistanceBytes(a, b []byte) int {
	max := min(len(b), len(a))
	count := 0
	for i := range max {
		count += bits.OnesCount8(a[i] ^ b[i])
	}
	return count
}

func hammingDistanceInt64(a, b int64, bitLength int) int {
	if bitLength <= 0 {
		return 0
	}

	if bitLength >= 64 {
		return bits.OnesCount64(uint64(a) ^ uint64(b))
	}

	mask := (uint64(1) << bitLength) - 1
	diff := (uint64(a) ^ uint64(b)) & mask
	return bits.OnesCount64(diff)
}

func logAvalanche(t *testing.T, label string, base, variant otpComponents) {
	hashDist := hammingDistanceBytes(base.hash[:], variant.hash[:])
	truncDist := hammingDistanceInt64(int64(base.truncated), int64(variant.truncated), 32)

	x1Bits := bits.Len64(uint64(lcmM - 1))
	if x1Bits == 0 {
		x1Bits = 1
	}
	x1Dist := hammingDistanceInt64(base.x1, variant.x1, x1Bits)

	otpBits := bits.Len64(999999)
	otpDist := hammingDistanceInt64(base.otp, variant.otp, otpBits)

	t.Logf("[%s] Seed A: %s", label, base.seed)
	t.Logf("[%s] Seed B: %s", label, variant.seed)
	t.Logf("[%s] SHA-256 Hamming: %d/256 (%.2f%%)", label, hashDist, float64(hashDist)/256.0*100)
	t.Logf("[%s] Trunc32 Hamming: %d/32 (%.2f%%)", label, truncDist, float64(truncDist)/32.0*100)
	t.Logf("[%s] LCM X1 Hamming: %d/%d (%.2f%%)", label, x1Dist, x1Bits, float64(x1Dist)/float64(x1Bits)*100)
	t.Logf("[%s] OTP(6-digit) Hamming: %d/%d (%.2f%%)", label, otpDist, otpBits, float64(otpDist)/float64(otpBits)*100)

	if base.otp == variant.otp {
		t.Logf("[%s] CATATAN: OTP kebetulan sama (kemungkinan kecil).", label)
	}
}

func TestAvalancheEffectLCM(t *testing.T) {
	baseUsername := "m_anaska_"
	baseTime := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)

	base := calcComponents(baseUsername, baseTime)

	// Variasi 1: ubah satu karakter username
	mutatedUsername := baseUsername[:len(baseUsername)-1] + "6"
	variantUsername := calcComponents(mutatedUsername, baseTime)

	// Variasi 2: geser waktu ke interval berikutnya (+30 detik)
	variantTime := calcComponents(baseUsername, baseTime.Add(30*time.Second))

	t.Logf("=== AVALANCHE EFFECT (LCM + SHA-256) ===")
	logAvalanche(t, "Ubah-Username-1Karakter", base, variantUsername)
	logAvalanche(t, "Ubah-Interval-Waktu", base, variantTime)
	t.Logf("========================================")
}

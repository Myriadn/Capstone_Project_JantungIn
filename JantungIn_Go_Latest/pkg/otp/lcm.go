package otp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Default LCM params (POSIX / Numerical Recipes)
const (
	defaultLCMA int64 = 1103515245
	defaultLCMC int64 = 12345
	defaultLCMM int64 = 2147483648
)

var (
	lcmOnce sync.Once
	lcmA    int64
	lcmC    int64
	lcmM    int64
)

// GenerateOTP membuat 6-digit OTP berdasarkan username dan waktu (interval 30 detik).
func GenerateOTP(username string, t time.Time) string {
	lcmOnce.Do(loadLCMParams)

	interval := t.Unix() / 30                           // Time Interval 30 detik
	seedStr := fmt.Sprintf("%s-%d", username, interval) // Username + Time Interval
	hash := sha256.Sum256([]byte(seedStr))              // Hashing dengan SHA-256
	hashString := hex.EncodeToString(hash[:])           // Ubah byte array menjadi string heksadesimal (total 64 karakter)
	truncatedHash := hashString[:8]                     // Mengambil 8 karakter pertama dari hash

	// if strings.EqualFold(strings.TrimSpace(os.Getenv("APP_ENV")), "development") {
	// 	log.Printf("[OTP-DEBUG] seed=%s truncatedHash=%s", seedStr, truncatedHash)
	// }

	x0, _ := strconv.ParseInt(truncatedHash, 16, 64) // Konversi 8 karakter heksadesimal sebagai nilai seed
	x1 := (lcmA*x0 + lcmC) % lcmM                    // Rumus LCM: X_{n+1} = (a * X_n + c) mod m
	otpVal := x1 % 1000000                           // Konversi ke 6-digit OTP

	// Mengembalikan format 6 digit angka (ditambah 0 di depan jika nilainya < 100000)
	return fmt.Sprintf("%06d", otpVal)
}

func loadLCMParams() {
	lcmA = getEnvInt64("OTP_LCM_A", defaultLCMA)
	lcmC = getEnvInt64("OTP_LCM_C", defaultLCMC)
	lcmM = getEnvInt64("OTP_LCM_M", defaultLCMM)

	if lcmA <= 0 {
		lcmA = defaultLCMA
	}
	if lcmC <= 0 {
		lcmC = defaultLCMC
	}
	if lcmM <= 0 {
		lcmM = defaultLCMM
	}
}

func getEnvInt64(key string, defaultValue int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return defaultValue
	}

	return parsed
}

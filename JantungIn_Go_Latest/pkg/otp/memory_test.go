package otp

import (
	"testing"
	"time"
)

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	username := "testuser"
	otpCode := "123456"

	t.Logf("=== MEMULAI TEST VIRTUAL MEMORY (IN-MEMORY) ===")

	// Simpan OTP dengan durasi 1 detik
	t.Logf("[1] Menyimpan OTP '%s' untuk user '%s' (TTL: 1 Detik)", otpCode, username)
	store.Save(username, otpCode, 1*time.Second)

	// Cek apakah OTP valid ketika belum kadaluarsa
	isValid := store.Validate(username, otpCode)
	if !isValid {
		t.Errorf("Harusnya OTP valid karena dimasukkan kode yang benar dan belum expired")
	} else {
		t.Logf("Skenario 1 Berhasil: OTP berhasil divalidasi.")
	}

	// Cek OTP yang sudah dipakai tidak boleh dipakai lagi
	t.Logf("[2] Mencoba memakai OTP yang sama untuk kedua kalinya...")
	isStillValid := store.Validate(username, otpCode)
	if isStillValid {
		t.Errorf("Harusnya OTP sudah tidak valid setelah digunakan (One-Time)")
	} else {
		t.Logf("Skenario 2 Berhasil: OTP ditolak karena sifatnya One-Time (langsung terhapus setelah sukses dipakai).")
	}

	// Test Expired
	username2 := "user_expired"
	otpCode2 := "654321"

	t.Logf("[3] Skenario Kedaluwarsa (Expired) - Menyimpan OTP '%s' untuk user '%s' (TTL: 1 Detik)", otpCode2, username2)
	store.Save(username2, otpCode2, 1*time.Second)

	t.Logf("[3] Menunggu selama 1.5 detik agar OTP kedaluwarsa...")
	time.Sleep(1500 * time.Millisecond)

	t.Logf("[3] Mencoba memvalidasi OTP yang sudah kedaluwarsa...")
	isExpiredValid := store.Validate(username2, otpCode2)
	if isExpiredValid {
		t.Errorf("Harusnya OTP tidak valid karena sudah lewat dari batas TTL di memori")
	} else {
		t.Logf("Skenario 3 Berhasil: OTP ditolak karena sudah lebih dari 1 detik (Terhapus oleh sistem).")
	}
	t.Logf("==================================================")
}

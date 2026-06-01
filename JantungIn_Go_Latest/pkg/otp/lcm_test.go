package otp

import (
	"testing"
	"time"
)

func TestGenerateOTP(t *testing.T) {
	username := "anis_setiawan445"

	// Waktu awal untuk pengetesan
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	t.Logf("=== MEMULAI TEST ALGORITMA OTP (LCM + SHA-256) ===")
	t.Logf("Username yang diuji : %s", username)
	t.Logf("Waktu Dasar (T0)    : %v", baseTime)

	// OTP untuk detik 0-29 harus SAMA
	otp1 := GenerateOTP(username, baseTime)
	t.Logf("[Interval 1] OTP detik ke-0  : %s", otp1)

	otp2 := GenerateOTP(username, baseTime.Add(15*time.Second)) // +15 detik
	t.Logf("[Interval 1] OTP detik ke-15 : %s", otp2)

	otp3 := GenerateOTP(username, baseTime.Add(29*time.Second)) // +29 detik
	t.Logf("[Interval 1] OTP detik ke-29 : %s", otp3)

	if otp1 != otp2 || otp1 != otp3 {
		t.Errorf("OTP di dalam satu blok waktu (30s) harusnya sama.")
	} else {
		t.Logf("Berhasil: OTP tetap konstan selama 30 detik pertama.")
	}

	// Memastikan panjang OTP persis 6 karakter
	if len(otp1) != 6 {
		t.Errorf("Panjang OTP harus 6 digit, didapat: %d digit", len(otp1))
	}

	// OTP di blok waktu BERIKUTNYA harus BERBEDA
	otpNextBlock := GenerateOTP(username, baseTime.Add(31*time.Second))
	t.Logf("[Interval 2] OTP detik ke-31 : %s (Blok waktu baru)", otpNextBlock)

	if otp1 == otpNextBlock {
		t.Errorf("OTP dari blok 30 detik berbeda tidak boleh sama.")
	} else {
		t.Logf("Berhasil: OTP berubah setelah melewati 30 detik.")
	}

	// OTP untuk username yang BERBEDA pada waktu yang SAMA juga harus BERBEDA
	otpOtherUser := GenerateOTP("dosen_penguji", baseTime)
	t.Logf("[User Beda] OTP user 'dosen_penguji' pada T0: %s", otpOtherUser)
	if otp1 == otpOtherUser {
		t.Errorf("OTP dari user berbeda tidak boleh sama.")
	} else {
		t.Logf("Berhasil: OTP unik untuk setiap username walau di detik yang sama.")
	}
	t.Logf("==================================================")
}

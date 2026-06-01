package otp

import (
	"sync"
	"time"
)

// OTPItem struktur untuk menyimpan OTP berserta batas kedaluwarsanya di memory
type OTPItem struct {
	Code      string
	ExpiresAt time.Time
}

// MemoryStore mengatur penyimpanan OTP menggunakan In-Memory (Virtual Memory RAM)
type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]OTPItem
}

// NewMemoryStore inisialisasi penyimpanan baru
func NewMemoryStore() *MemoryStore {
	m := &MemoryStore{
		store: make(map[string]OTPItem),
	}
	// Anda juga dapat memanggil goroutine untuk membersihkan memory secara periodik di background
	go m.cleanupRoutine()
	return m
}

// Save menyimpan OTP untuk username dengan masa aktif (duration), umumnya 30 detik
func (m *MemoryStore) Save(username, code string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[username] = OTPItem{
		Code:      code,
		ExpiresAt: time.Now().Add(duration),
	}
}

// Validate mengecek apakah OTP yang dikirim user sesuai dan belum kadaluarsa
func (m *MemoryStore) Validate(username, code string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.store[username]
	if !exists {
		return false // OTP tidak ditemukan di memory
	}

	// Cek apakah sudah melebihi masa aktif
	if time.Now().After(item.ExpiresAt) {
		delete(m.store, username) // Hapus dari memori (Virtual Memory)
		return false
	}

	// Cek apakah kodenya cocok
	if item.Code == code {
		// Jika sukses divalidasi, maka OTP langsung dihapus (One-Time)
		delete(m.store, username)
		return true
	}

	return false
}

// cleanupRoutine berjalan di background (goroutine)
// untuk menghapus OTP yang sudah expire agar RAM tidak penuh
func (m *MemoryStore) cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Second) // Setiap 5 detik akan melakukan sapuan (garbage collection manual)
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for user, item := range m.store {
			if now.After(item.ExpiresAt) {
				delete(m.store, user)
			}
		}
		m.mu.Unlock()
	}
}

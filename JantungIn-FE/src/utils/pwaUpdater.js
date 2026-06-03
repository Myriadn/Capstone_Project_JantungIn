/**
 * PWA Update Handler
 *
 * Modul ini menangani pembaruan otomatis service worker dan
 * menawarkan opsi bagi pengguna untuk memperbarui aplikasi
 */

// Jangka waktu cek pembaruan dalam milidetik (15 menit)
const UPDATE_CHECK_INTERVAL = 15 * 60 * 1000

/**
 * Mengelola pembaruan service worker dan cache aplikasi
 */
class PWAUpdater {
  constructor() {
    this.registration = null
    this.updateFound = false
    this.updateAvailable = false
    this.checkInterval = null
    this.callbacks = {
      onUpdateFound: [],
      onUpdateReady: [],
      onNoUpdateAvailable: [],
      onOffline: [],
    }
  }

  /**
   * Inisialisasi pembaruan service worker
   */
  async init() {
    if (!('serviceWorker' in navigator)) {
      console.log('Service worker tidak didukung')
      return
    }

    try {
      // Dapatkan registrasi service worker
      this.registration = await navigator.serviceWorker.getRegistration()

      if (!this.registration) {
        console.log('Tidak ada service worker yang terdaftar')
        return
      }

      // Pantau perubahan status service worker
      this.trackWorkerStateChanges()

      // Mulai pemeriksaan pembaruan berkala
      this.startPeriodicUpdateCheck()

      // Periksa langsung pembaruan
      this.checkForUpdate()
    } catch (error) {
      console.error('Kesalahan saat inisialisasi PWA updater:', error)
    }
  }

  /**
   * Menambahkan callback event
   * @param {string} event Nama event
   * @param {Function} callback Fungsi callback
   */
  on(event, callback) {
    if (this.callbacks[event]) {
      this.callbacks[event].push(callback)
    }
  }

  /**
   * Memicu event dengan data tertentu
   * @param {string} event Nama event
   * @param {any} data Data untuk callback
   */
  trigger(event, data) {
    if (this.callbacks[event]) {
      this.callbacks[event].forEach((callback) => callback(data))
    }
  }

  /**
   * Memantau perubahan status service worker
   */
  trackWorkerStateChanges() {
    if (!this.registration) return

    // Deteksi pembaruan service worker
    this.registration.addEventListener('updatefound', () => {
      const newWorker = this.registration.installing
      this.updateFound = true
      this.trigger('onUpdateFound', { worker: newWorker })

      console.log('Service worker baru ditemukan, menunggu instalasi...')

      newWorker.addEventListener('statechange', () => {
        if (newWorker.state === 'installed') {
          if (navigator.serviceWorker.controller) {
            // Service worker baru diinstal, tetapi service worker lama masih mengontrol halaman
            console.log('Update tersedia! Refresh halaman untuk menggunakan versi terbaru.')
            this.updateAvailable = true
            this.trigger('onUpdateReady', { worker: newWorker })
          } else {
            // Service worker diinstal untuk pertama kali
            console.log('Konten di-cache untuk penggunaan offline.')
          }
        }
      })
    })

    // Deteksi saat service worker diperbarui
    navigator.serviceWorker.addEventListener('controllerchange', () => {
      if (this.updateAvailable) {
        console.log('Service worker baru aktif, reload halaman...')
        this.updateAvailable = false // biar gak loop
        window.location.reload()
      }
    })
  }

  /**
   * Mulai pemeriksaan pembaruan secara berkala
   */
  startPeriodicUpdateCheck() {
    this.checkInterval = setInterval(() => {
      if (navigator.onLine) {
        this.checkForUpdate()
      }
    }, UPDATE_CHECK_INTERVAL)
  }

  /**
   * Periksa pembaruan service worker
   */
  async checkForUpdate() {
    if (!this.registration || !navigator.onLine) {
      if (!navigator.onLine) {
        this.trigger('onOffline')
      }
      return
    }

    try {
      console.log('Memeriksa pembaruan service worker...')
      await this.registration.update()

      if (!this.updateFound) {
        console.log('Tidak ada pembaruan yang tersedia')
        this.trigger('onNoUpdateAvailable')
      }
    } catch (error) {
      console.error('Kesalahan saat memeriksa pembaruan:', error)
    }
  }

  /**
   * Paksa perbarui dan refresh aplikasi
   */
  applyUpdate() {
    if (!this.registration || !this.updateAvailable) return

    const worker = this.registration.waiting

    if (worker) {
      // Kirim pesan ke service worker untuk memperbarui
      worker.postMessage({ type: 'SKIP_WAITING' })
    } else {
      // Jika tidak ada worker yang menunggu, coba refresh halaman saja
      window.location.reload()
    }
  }

  /**
   * Bersihkan interval pengecekan pembaruan
   */
  cleanup() {
    if (this.checkInterval) {
      clearInterval(this.checkInterval)
    }
  }
}

// Ekspor instance tunggal
export const pwaUpdater = new PWAUpdater()

// Tambahkan listener di service worker untuk menerima pesan SKIP_WAITING
if ('serviceWorker' in navigator && navigator.serviceWorker.controller) {
  navigator.serviceWorker.addEventListener('message', (event) => {
    if (event.data && event.data.type === 'UPDATE_AVAILABLE') {
      console.log('Service worker memberitahukan pembaruan tersedia')
      pwaUpdater.updateAvailable = true
      pwaUpdater.trigger('onUpdateReady', { source: 'worker-message' })
    }
  })
}

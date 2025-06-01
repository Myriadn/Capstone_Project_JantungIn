# 🫀 JantungIn API

<p align="center">
  <img src="./public/logo.png" alt="JantungIn Logo" width="300"/>
</p>

<p align="center">
  <b>Solusi Cerdas untuk Deteksi Dini Risiko Penyakit Kardiovaskular</b>
</p>

JantungIn API adalah layanan backend RESTful untuk aplikasi JantungIn yang berfokus pada penilaian risiko penyakit kardiovaskular. API ini mendukung otentikasi pengguna, manajemen data pasien, dan prediksi risiko kardiovaskular menggunakan model pembelajaran mesin.

## ✨ Fitur Utama

- 🔐 **Otentikasi Pengguna** - Daftar, masuk, dan kelola profil
- 🔍 **Prediksi Risiko Penyakit Kardiovaskular** - Menggunakan TensorFlow.js untuk analisis data kesehatan
- 📊 **Pelacakan Riwayat Diagnosis** - Menyimpan dan menampilkan riwayat diagnosis pasien
- 🌐 **API RESTful** - Dengan endpoint terstandarisasi
- 💾 **Database PostgreSQL** - Penyimpanan data yang handal
- 🚀 **Opsi Deployment yang Aman** - Dukungan untuk AWS, Railway, dan Docker

## 🛠️ Teknologi yang Digunakan

- **Node.js** - Lingkungan runtime JavaScript
- **Hapi.js** - Framework web yang kuat dan fleksibel
- **PostgreSQL** - Sistem manajemen database relasional
- **JWT** - JSON Web Token untuk otentikasi yang aman
- **TensorFlow.js** - Library pembelajaran mesin untuk prediksi

## 📋 Prasyarat

- **Node.js v14+** dan npm
- **PostgreSQL** database
- **Git** untuk kontrol versi

## 🚀 Instalasi

1. Clone repositori:

```bash
git clone <repository-url>
cd JantungIn_API
```

2. Install dependensi:

```bash
npm install
```

3. Siapkan variabel lingkungan:

```bash
cp .env.example .env
```

Edit file `.env` untuk mengatur kredensial database dan opsi konfigurasi lainnya.

## 💾 Konfigurasi Database

### Konfigurasi PostgreSQL

1. Konfigurasi pengaturan database PostgreSQL di file `.env`:

```
DB_DIALECT=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=jantungin
DB_USER=username_anda
DB_PASSWORD=password_anda
```

2. Pastikan server PostgreSQL berjalan dan database telah dibuat.

## 🏃‍♂️ Menjalankan Aplikasi

### Mode Pengembangan

```bash
npm run dev
```

### Mode Produksi

```bash
npm start
```

## 📡 Endpoint API

### 🔐 Autentikasi

- `POST /api/auth/register` - Mendaftarkan pengguna baru
- `POST /api/auth/login` - Login pengguna
- `GET /api/auth/profile` - Mendapatkan profil pengguna
- `PUT /api/auth/profile` - Memperbarui profil pengguna

### 🩺 Diagnosis

- `POST /api/diagnosis` - Membuat diagnosis baru
- `GET /api/diagnosis/history` - Mendapatkan riwayat diagnosis pengguna
- `GET /api/diagnosis/{id}` - Mendapatkan diagnosis spesifik

### 👩‍💼 Admin (Khusus Dokter)

- `GET /api/admin/users` - Mendapatkan daftar semua pengguna
- `GET /api/admin/diagnoses` - Mendapatkan semua data diagnosis
- `PUT /api/admin/users/{id}` - Memperbarui data pengguna

## 🧪 Pengujian

```bash
npm test
```

## 🚀 Deployment

### 🚂 Deployment di Railway

JantungIn API dapat dengan mudah di-deploy di [Railway](https://railway.app):

1. Hubungkan repositori GitHub Anda ke Railway
2. Siapkan variabel lingkungan yang diperlukan:
   - `PORT` (default: 3000)
   - `NODE_ENV=production`
   - `JWT_SECRET`
   - String koneksi database (berdasarkan plugin PostgreSQL Railway)
   - Variabel lingkungan spesifik aplikasi lainnya
3. Deploy aplikasi Anda

### ☁️ Langkah-langkah Deployment AWS

1. Siapkan instance EC2
2. Konfigurasi variabel lingkungan untuk produksi
3. Gunakan process manager seperti PM2 untuk mengelola aplikasi Node.js:

```bash
npm install -g pm2
pm2 start src/server.js --name jantungin-api
```

### 🐳 Deployment dengan Docker

Dockerfile disertakan untuk mengkontainerisasi aplikasi. Build dan jalankan dengan:

```bash
docker build -t jantungin-api .
docker run -p 3000:3000 -e NODE_ENV=production jantungin-api
```

## 🔧 Pemecahan Masalah

### 🚨 Masalah Umum saat Deployment

1. **Kesalahan Middleware** - Pastikan fungsi middleware kompatibel dengan Hapi.js. Middleware gaya Express (`app.use()`) tidak akan berfungsi dengan Hapi.js.

2. **Masalah Koneksi Database** - Verifikasi kredensial database dan pastikan database dapat diakses dari lingkungan deployment Anda.

3. **Variabel Lingkungan** - Periksa bahwa semua variabel lingkungan yang diperlukan diatur dengan benar di platform deployment Anda.

4. **Konflik Port** - Pastikan port yang ditentukan dalam aplikasi Anda tersedia dan tidak diblokir oleh firewall.

## 📈 Arsitektur Sistem

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │
│   Frontend  │────▶│  JantungIn  │────▶│  PostgreSQL │
│   (Vue.js)  │     │     API     │     │  Database   │
│             │◀────│  (Hapi.js)  │◀────│             │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │ TensorFlow  │
                    │    Model    │
                    └─────────────┘
```

## 📝 Alur Kerja

1. Pengguna melakukan autentikasi melalui frontend
2. Frontend mengirimkan data kesehatan ke API
3. API memproses data menggunakan model TensorFlow.js
4. Hasil prediksi dikembalikan ke pengguna
5. Data diagnosis disimpan dalam database

## 📜 Lisensi

Proyek ini dilisensikan di bawah Lisensi ISC.

## 👨‍💻 Kontributor

- Myriadn - Pengembangan awal dan desain API

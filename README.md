# Go Findest Id REST API

Proyek ini adalah implementasi REST API menggunakan Go (Golang) dengan framework gin.

## Prasyarat

Pastikan Anda telah menginstal hal-hal berikut sebelum melanjutkan:

- Go (Golang) versi terbaru: https://golang.org/doc/install
- Database PostgreSQL

## Langkah-langkah Setup Proyek

### 1. Clone Repository

Clone repositori ini ke mesin lokal Anda:

```bash
git clone https://github.com/irfanyahyaa/go-findest-rest-api.git
cd go-findest-rest-api
```

### 2. Setup Database
Sebelum menjalankan aplikasi, Anda perlu membuat database dan mengatur koneksi ke database tersebut. Ikuti langkah-langkah berikut:

#### PostgreSQL:
1. Buat database baru di PostgreSQL:
```sql
CREATE DATABASE findest;
```
2. Update konfigurasi koneksi database di file `.env` atau file konfigurasi yang relevan dengan pengaturan Anda.

### 3. Migrasi Database
Untuk memigrasikan database, anda bisa langsung menjalankan aplikasinya. Karena repo ini sudah menghandle hal tersebut menggunakan `autoMigrate` milik `GORM`

### 4. Install Dependencies
Pasang semua dependencies yang diperlukan dengan perintah berikut:
```bash
go mod tidy
```

### 5. Menjalankan Aplikasi
Setelah semua konfigurasi selesai, jalankan aplikasi dengan perintah:
```bash
go run main.go
```
API akan berjalan di port default yaitu `8080` (http://localhost:8080).

## Pengujian
Untuk menjalankan pengujian, gunakan perintah berikut:
```bash
./test-all.sh
```
Di dalam file `test-all.sh` tersebut sudah disediakan `script` untuk menjalankan unit test yang ditujukan di package `controller`, begitu juga dengan `coverage` dari unit test yang dilakukan. Hasil `coverage` bisa dilihat di `coverage.html`

## Dokumentasi API
Untuk dokumentasi API saya sudah menyediakan dalam bentuk `Postman Collection` di link berikut: 
https://documenter.getpostman.com/view/27765876/2sAYX5MPHY#intro
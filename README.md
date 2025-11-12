# Template Struktur Golang
![Go](https://img.shields.io/badge/Go-1.21-blue.svg)
![Gin](https://img.shields.io/badge/Gin-1.11.0-lightgrey.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

Proyek ini adalah **template backend REST API menggunakan Go dengan Gin**, dengan struktur rapi mirip Laravel.  
Template ini sudah termasuk folder `controllers`, `models`, `services`, `routes`, `utils`, `config`, file `.env`, dan `config.go` untuk load environment variables.

---

## 📂 Struktur Proyek

```plaintext
myapp/
├── .env
├── go.mod
├── cmd/
│   └── main.go
└── app/
    ├── config/
    │   └── config.go
    ├── controllers/
    │   └── home_controller.go
    ├── models/
    │   └── user.go
    ├── services/
    │   └── user_service.go
    ├── utils/
    │   └── hash.go
    └── routes/
        └── routes.go
```

markdown
   

---

## ⚡ Fitur

- Struktur proyek **terpisah dan rapi** (mirip MVC)
- Template file untuk:
  - Controller
  - Model
  - Service
  - Utils
  - Routes
  - Config
- File `.env` standar untuk server, database, JWT, dan logging
- Bisa langsung dijalankan:
```bash
go run cmd/main.go
```
Gin framework sudah terinstall otomatis (github.com/gin-gonic/gin)

🚀 Cara Menjalankan
Clone template Gin Generator (lokasi bebas, misal di D:\GinGenerator):

Buka folder di Explorer, misal D:\GinGenerator

Di folder itu sudah ada setup.go

Masuk ke folder projek tujuan (misal di D:\BERKAS MAGANG\KOMINFO MERAUKE\JadwalGubernur)

Klik kanan di folder → pilih “Open in Terminal”

Atau tekan Shift + Klik kanan → Open command window here

Jalankan generator untuk membuat proyek baru

Ketik perintah di CMD (contoh):

```bash
   
go run D:\GinGenerator\setup.go
```
Masukkan:

Nama proyek → misal: myapp

Path folder tujuan → tekan Enter jika ingin di folder CMD saat ini

Hasilnya:

Folder myapp otomatis dibuat di lokasi tujuan

Struktur proyek lengkap dengan semua folder (cmd, app/routes, dll) dan template file (main.go, .env, config.go)

Jalankan server Gin

Masuk ke folder proyek:

```bash
   
cd myapp
```
Jalankan server:

bash
   
go run cmd/main.go
Buka browser / Postman ke:

http
   
http://localhost:8080/
Response:

json
   
{
  "message": "Hello from Gin!"
}
Cara melihat lokasi clone template

Buka folder yang berisi setup.go → klik direktori di atas →   path

Gunakan path ini di CMD saat menjalankan:

```bash
   
go run D:\Lokasi\GinGenerator\setup.go
```
Gunakan path ini di CMD saat menjalankan:

```bash
   
go run D:\Lokasi\GinGenerator\setup.go
```

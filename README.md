# Todo API (Go) — Dokumentasi Penggunaan

Dokumentasi ini menjelaskan cara instalasi, menjalankan, dan menggunakan Todo REST API yang dibuat dengan Go, GORM (SQLite) dan JWT.

Isi dokumentasi:
- Ringkasan
- Prasyarat
- Instalasi & menjalankan server
- Konfigurasi
- Daftar endpoint lengkap (request + response contoh)
- Contoh `curl`, `bash` (ambil token), dan contoh JavaScript `fetch`
- UI web testing
- Troubleshooting & tips

---

## Ringkasan

API ini menyediakan fitur dasar manajemen todo untuk tiap pengguna:

- Register (mendaftar akun)
- Login (mengembalikan JWT token)
- CRUD Todo (Create, Read, Update, Delete) — hanya dapat diakses bila memiliki token

Server default: `http://localhost:8080`

Base path API: `/api`

Database default: SQLite file `todo.db` (di folder proyek)

---

## Prasyarat

- Go 1.21+ (direkomendasikan 1.21 atau yang lebih baru)
- `curl` (untuk contoh) atau Postman / HTTP client lainnya
- (opsional) `jq` untuk parsing JSON di terminal

---

## Instalasi & Jalankan

1. Masuk ke folder proyek:

```bash
cd /home/harsya/Unduhan/projects
```

2. Unduh dependensi dan rapikan module:

```bash
go mod tidy
```

3. Jalankan server:

```bash
go run main.go
```

Atau buat binary dan jalankan:

```bash
go build -o todo-api .
./todo-api
```

Server akan menampilkan:

```
✅ Database terhubung: todo.db
🚀 Server berjalan di http://localhost:8080
```

Jika port `8080` sudah digunakan, hentikan proses lain yang memakai port tersebut atau ubah `config.GetConfig()` untuk port lain.

---

## Konfigurasi

Konfigurasi sederhana ada di file `config/config.go`. Nilai default:

- Port: `8080`
- JWTSecret: `rahasia-jwt-super-aman-ganti-ini-di-produksi`
- DBPath: `todo.db`

Untuk produksi, ganti nilai ini agar `JWTSecret` aman dan gunakan database server jika perlu.

---

## Endpoint API

Semua endpoint berada di prefix `/api`.

1) POST /api/register — Registrasi

- Auth: Tidak
- Body JSON:

```json
{
  "name": "Nama",
  "email": "email@example.com",
  "password": "rahasia"
}
```

- Response (201 Created):

```json
{
  "message": "Registrasi berhasil",
  "user": { "id": 1, "name": "Nama", "email": "...", "created_at": "..." }
}
```

2) POST /api/login — Login

- Auth: Tidak
- Body JSON:

```json
{
  "email": "email@example.com",
  "password": "rahasia"
}
```

- Response (200 OK):

```json
{
  "message": "Login berhasil",
  "token": "<JWT_TOKEN>",
  "user": { "id": 1, "name": "Nama", "email": "..." }
}
```

3) GET /api/todos — Ambil semua todo milik user

- Auth: Ya (Header `Authorization: Bearer <token>`)
- Response (200 OK):

```json
{
  "todos": [ { "id":1, "user_id":1, "title":"...", "completed":false, ... } ],
  "total": 1
}
```

4) GET /api/todos/{id} — Ambil satu todo

- Auth: Ya
- Params: `id` (path)
- Response (200 OK): JSON todo object atau 404 jika tidak ditemukan

5) POST /api/todos — Buat todo baru

- Auth: Ya
- Body JSON:

```json
{
  "title": "Judul",
  "description": "Opsional"
}
```

- Response (201 Created):

```json
{ "message": "Todo berhasil dibuat", "todo": { ... } }
```

6) PUT /api/todos/{id} — Update todo

- Auth: Ya
- Body JSON (kirim hanya field yang mau diubah):

```json
{
  "title": "Judul baru",
  "description": "...",
  "completed": true
}
```

- Response (200 OK)

7) DELETE /api/todos/{id} — Hapus todo

- Auth: Ya
- Response (200 OK)

---

## Contoh Penggunaan (curl)

1) Register (contoh):

```bash
curl -i -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Budi","email":"budi@example.com","password":"rahasia123"}'
```

2) Login dan ambil token (dengan `jq` jika tersedia):

```bash
# Dengan jq
token=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"budi@example.com","password":"rahasia123"}' | jq -r .token)

# Tanpa jq (grep+sed)
token=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"budi@example.com","password":"rahasia123"}' | grep -oP '"token":"\K[^"]+')

echo "Token: $token"
```

3) Buat todo (pakai token):

```bash
curl -i -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $token" \
  -d '{"title":"Belajar Go","description":"Contoh penggunaan API"}'
```

4) Ambil daftar todo:

```bash
curl -i -H "Authorization: Bearer $token" http://localhost:8080/api/todos
```

5) Update todo:

```bash
curl -i -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $token" \
  -d '{"completed":true}'
```

6) Hapus todo:

```bash
curl -i -X DELETE http://localhost:8080/api/todos/1 \
  -H "Authorization: Bearer $token"
```

---

## Contoh JavaScript (fetch)

Login + panggil protected endpoint:

```javascript
async function loginAndFetchTodos() {
  const loginRes = await fetch('http://localhost:8080/api/login', {
    method: 'POST', headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({email: 'budi@example.com', password: 'rahasia123'})
  });
  const loginJson = await loginRes.json();
  const token = loginJson.token;

  const todosRes = await fetch('http://localhost:8080/api/todos', {
    headers: { 'Authorization': 'Bearer ' + token }
  });
  const todos = await todosRes.json();
  console.log(todos);
}

loginAndFetchTodos();
```

---

## UI Web untuk testing

File `web/index.html` sudah disertakan dan disajikan oleh server root (`/`). Buka `http://localhost:8080` di browser untuk menggunakan UI sederhana yang menyediakan register/login dan CRUD todo.

---

## Troubleshooting

- 404 pada `GET /api/register` atau `GET /api/login` → normal karena endpoint hanya menerima `POST`.
- Jika `go run main.go` gagal: pastikan port `8080` tersedia. Cari proses yang memakai port:

```bash
ss -ltnp | grep ':8080' || true
```

Matikan proses tersebut atau ubah port di `config/config.go`.

- Jika error saat build terkait `cgo` / `sqlite3`, pastikan toolchain C (gcc) terpasang. Di Linux biasanya:

```bash
sudo apt install build-essential
```

atau setara distro Anda.

- Reset database: hapus file `todo.db` lalu jalankan ulang server (migrasi otomatis akan membuat tabel baru):

```bash
rm todo.db
go run main.go
```

---

## Keamanan & Produksi

- Jangan gunakan `JWTSecret` default di produksi. Ambil dari environment variable.
- Batasi CORS di produksi (jangan pakai `*`).
- Pertimbangkan mengganti SQLite dengan PostgreSQL/MySQL untuk aplikasi nyata.
- Tambahkan rate-limiting dan logging terpusat untuk produksi.

---

## Testing otomatis & build

- Jalankan unit-check (tidak ada file test saat ini):

```bash
go test ./...
```

- Build release:

```bash
go build -o todo-api .
```

---

## Contoh skrip lengkap (bash) — register, login, create, list

```bash
#!/usr/bin/env bash
set -euo pipefail

BASE=http://localhost:8080/api
EMAIL=test+cli@example.com
PASS=rahasia123

curl -s -X POST "$BASE/register" -H 'Content-Type: application/json' \
  -d "{\"name\":\"CLI Test\",\"email\":\"$EMAIL\",\"password\":\"$PASS\"}"

TOKEN=$(curl -s -X POST "$BASE/login" -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASS\"}" | grep -oP '"token":"\K[^"]+')

echo "Token: $TOKEN"

curl -s -X POST "$BASE/todos" -H 'Content-Type: application/json' -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Dari skrip CLI","description":"created via script"}' | jq

curl -s -H "Authorization: Bearer $TOKEN" "$BASE/todos" | jq
```

---

Jika Anda ingin, saya bisa:
- Menambahkan file `POSTMAN_COLLECTION.json` untuk import ke Postman
- Menambahkan `Makefile` atau skrip `run.sh` untuk dev convenience
- Menambahkan environment variable support (`PORT`, `JWT_SECRET`, `DB_PATH`)

Mau saya tambahkan salah satu dari hal tersebut sekarang?  

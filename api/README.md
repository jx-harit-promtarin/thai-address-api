# Thai Address API

REST API สำหรับข้อมูลที่อยู่ประเทศไทย (จังหวัด / อำเภอ / ตำบล / ภูมิภาค)
เขียนด้วย Go ทำหน้าที่เป็น proxy ไปยัง [Supabase](https://supabase.com/) REST API
โดยเก็บ `apikey` ไว้ฝั่ง server เพื่อไม่ให้หลุดไปถึง client

## Endpoints

ทุก endpoint รับเฉพาะ method `GET` และคืนค่าเป็น JSON

| Method | Path | คำอธิบาย | Supabase table |
| ------ | ---- | -------- | -------------- |
| GET | `/api/v1/provinces` | รายการจังหวัด | `province` |
| GET | `/api/v1/districts` | รายการอำเภอ/เขต | `district` |
| GET | `/api/v1/sub-districts` | รายการตำบล/แขวง | `sub_district` |
| GET | `/api/v1/geographies` | รายการภูมิภาค | `geography` |

### Query parameters (filter)

บาง endpoint รองรับการ filter ผ่าน query string โดยส่งได้หลายตัวพร้อมกัน
(ค่าจะถูก AND เข้าด้วยกัน) — filter ที่เป็น `*_name_*` จะ join ไปยังตารางที่เกี่ยวข้อง

ทุก endpoint รองรับ `id` เพื่อ filter ตาม id ของ record นั้น ๆ

| Endpoint | Query param | คำอธิบาย |
| -------- | ----------- | -------- |
| `/api/v1/provinces` | `id` | filter จังหวัดตาม id |
| | `geography_id` | filter จังหวัดตาม id ของภูมิภาค |
| | `geography_name` | filter จังหวัดตามชื่อภูมิภาค |
| `/api/v1/districts` | `id` | filter อำเภอตาม id |
| | `province_id` | filter อำเภอตาม id ของจังหวัด |
| | `province_name_th` | filter อำเภอตามชื่อจังหวัด (ไทย) |
| | `province_name_en` | filter อำเภอตามชื่อจังหวัด (อังกฤษ) |
| `/api/v1/sub-districts` | `id` | filter ตำบลตาม id |
| | `district_id` | filter ตำบลตาม id ของอำเภอ |
| | `district_name_th` | filter ตำบลตามชื่ออำเภอ (ไทย) |
| | `district_name_en` | filter ตำบลตามชื่ออำเภอ (อังกฤษ) |
| | `province_id` | filter ตำบลตาม id ของจังหวัด (ผ่านอำเภอ) |
| `/api/v1/geographies` | `id` | filter ภูมิภาคตาม id |

### ตัวอย่าง

```bash
# รายการจังหวัดทั้งหมด
curl "http://localhost:8080/api/v1/provinces"

# filter ตาม id (ใช้ได้ทุก endpoint)
curl "http://localhost:8080/api/v1/provinces?id=1"
curl "http://localhost:8080/api/v1/geographies?id=2"

# จังหวัดในภูมิภาคที่ระบุ
curl "http://localhost:8080/api/v1/provinces?geography_id=2"
curl "http://localhost:8080/api/v1/provinces?geography_name=ภาคกลาง"

# อำเภอในจังหวัดที่ระบุ
curl "http://localhost:8080/api/v1/districts?province_id=1"
curl "http://localhost:8080/api/v1/districts?province_name_th=กรุงเทพมหานคร"
curl "http://localhost:8080/api/v1/districts?province_name_en=Bangkok"

# ตำบลในอำเภอที่ระบุ
curl "http://localhost:8080/api/v1/sub-districts?district_id=1"
curl "http://localhost:8080/api/v1/sub-districts?district_name_th=เขตพระนคร"
```

## Environment variables

| ตัวแปร | จำเป็น | คำอธิบาย |
| ------ | ------ | -------- |
| `SUPABASE_URL` | ✅ | URL ของโปรเจกต์ Supabase เช่น `https://xxxx.supabase.co` |
| `SUPABASE_KEY` | ✅ | anon key ของ Supabase (ใช้ฝั่ง server เท่านั้น) |
| `PORT` | ❌ | port ที่ให้ server ฟัง (ค่าเริ่มต้น `8080`) |

หาก `SUPABASE_URL` หรือ `SUPABASE_KEY` ไม่ถูกตั้งค่า server จะไม่ start

### ไฟล์ `.env`

ตอน start server จะโหลดค่าจากไฟล์ `.env` ในไดเรกทอรีนี้อัตโนมัติ (ถ้ามี)
สร้างไฟล์ `api/.env`:

```env
SUPABASE_URL=https://xxxx.supabase.co
SUPABASE_KEY=your-anon-key
PORT=8080
```

> ⚠️ อย่า commit ไฟล์ `.env` ขึ้น git — ควรเพิ่ม `.env` ไว้ใน `.gitignore`

## การรัน

### รันด้วย Go โดยตรง

```bash
# ตั้งค่าใน .env แล้วรันได้เลย
go run main.go

# หรือกำหนดผ่าน environment variable โดยตรง
export SUPABASE_URL="https://xxxx.supabase.co"
export SUPABASE_KEY="your-anon-key"
go run main.go
```

### รันด้วย Docker

```bash
# build image
docker build -t thai-address-api .

# run container
docker run -p 8080:8080 \
  -e SUPABASE_URL="https://xxxx.supabase.co" \
  -e SUPABASE_KEY="your-anon-key" \
  thai-address-api
```

## Requirements

- Go 1.22 ขึ้นไป (สำหรับรันแบบ local)
- Docker (สำหรับรันแบบ container)

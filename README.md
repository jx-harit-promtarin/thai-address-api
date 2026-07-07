# Thai Address API

A free, open REST API for Thailand's administrative divisions —
province (จังหวัด), district (อำเภอ), sub-district (ตำบล), and geography/region (ภูมิภาค) —
plus world-level continent (ทวีป) and country (ประเทศ) reference data with ISO 3166 codes.

> REST API สาธารณะแบบเปิดสำหรับข้อมูลเขตการปกครองของไทย —
> จังหวัด, อำเภอ/เขต, ตำบล/แขวง และภูมิภาค
> พร้อมข้อมูลอ้างอิงระดับโลก ทวีป (continent) และประเทศ (country) พร้อมโค้ด ISO 3166

A small [Go](api/) service acts as a proxy in front of [Supabase](https://supabase.com/)
(PostgreSQL + PostgREST), keeping the Supabase `apikey` server-side so it never reaches the client.
Deployed on Google Cloud Run.

> เบื้องหลังเป็นเซอร์วิส [Go](api/) ตัวเล็ก ๆ ทำหน้าที่เป็น proxy อยู่หน้า [Supabase](https://supabase.com/)
> (PostgreSQL + PostgREST) โดยเก็บ `apikey` ไว้ฝั่ง server ไม่ให้หลุดไปถึง client — deploy อยู่บน Google Cloud Run

Data sourced from [kongvut/thai-province-data](https://github.com/kongvut/thai-province-data) (MIT License).

> ข้อมูลนำมาจาก [kongvut/thai-province-data](https://github.com/kongvut/thai-province-data) (สัญญาอนุญาต MIT)

☕ If this API is useful to you, consider [supporting the project](#support--tip--สนับสนุนโปรเจกต์).
&nbsp;·&nbsp; ถ้า API นี้มีประโยชน์ ฝาก [สนับสนุนโปรเจกต์](#support--tip--สนับสนุนโปรเจกต์) ด้วยนะครับ

## Base URL

```
https://thai-address-api-373901862529.asia-southeast1.run.app
```

All endpoints below are relative to this base URL. No API key required.
&nbsp;·&nbsp; ทุก endpoint ด้านล่างต่อท้ายจาก base URL นี้ และเรียกใช้ได้เลยโดยไม่ต้องมี API key

## Endpoints

All endpoints accept only `GET` and return JSON.
&nbsp;·&nbsp; ทุก endpoint รับเฉพาะ method `GET` และคืนค่าเป็น JSON

| Method | Path | Description / คำอธิบาย |
| ------ | ---- | --------------------- |
| GET | `/api/v1/continents` | Continents / ทวีป |
| GET | `/api/v1/countries` | Countries / ประเทศ |
| GET | `/api/v1/geographies` | Geographies, regions / ภูมิภาค |
| GET | `/api/v1/provinces` | Provinces / จังหวัด |
| GET | `/api/v1/districts` | Districts / อำเภอ, เขต |
| GET | `/api/v1/sub-districts` | Sub-districts / ตำบล, แขวง |

### Query parameters (filters) / พารามิเตอร์สำหรับ filter

Every endpoint accepts `id` to fetch a single record. Other filters can be combined
(they are AND-ed together). Filters ending in `*_name_*` are resolved by joining the related table.

> ทุก endpoint รองรับ `id` เพื่อดึงเรคคอร์ดเดียว ส่วน filter อื่น ๆ ส่งพร้อมกันได้
> (ค่าจะถูก AND เข้าด้วยกัน) — filter ที่ลงท้ายด้วย `*_name_*` จะ join ไปยังตารางที่เกี่ยวข้อง

| Endpoint | Query param | Description / คำอธิบาย |
| -------- | ----------- | --------------------- |
| `/api/v1/continents` | `id` | filter by continent id / filter ทวีปตาม id |
| | `code` | filter by continent code (AS, EU, …) / filter ตามโค้ดทวีป |
| `/api/v1/countries` | `id` | filter by country id, ISO numeric / filter ประเทศตาม id (ISO numeric) |
| | `iso_alpha2` | filter by ISO alpha-2 code, e.g. TH / filter ตามโค้ด ISO alpha-2 |
| | `iso_alpha3` | filter by ISO alpha-3 code, e.g. THA / filter ตามโค้ด ISO alpha-3 |
| | `continent_id` | filter by continent id / filter ตาม id ของทวีป |
| | `continent_name_th` | filter by continent name, Thai / filter ตามชื่อทวีป (ไทย) |
| | `continent_name_en` | filter by continent name, English / filter ตามชื่อทวีป (อังกฤษ) |
| `/api/v1/geographies` | `id` | filter by region id / filter ภูมิภาคตาม id |
| | `country_id` | filter by country id / filter ตาม id ของประเทศ |
| | `country_name_th` | filter by country name, Thai / filter ตามชื่อประเทศ (ไทย) |
| | `country_name_en` | filter by country name, English / filter ตามชื่อประเทศ (อังกฤษ) |
| `/api/v1/provinces` | `id` | filter by province id / filter จังหวัดตาม id |
| | `geography_id` | filter by region id / filter ตาม id ของภูมิภาค |
| | `geography_name` | filter by region name / filter ตามชื่อภูมิภาค |
| `/api/v1/districts` | `id` | filter by district id / filter อำเภอตาม id |
| | `province_id` | filter by province id / filter ตาม id ของจังหวัด |
| | `province_name_th` | filter by province name, Thai / filter ตามชื่อจังหวัด (ไทย) |
| | `province_name_en` | filter by province name, English / filter ตามชื่อจังหวัด (อังกฤษ) |
| `/api/v1/sub-districts` | `id` | filter by sub-district id / filter ตำบลตาม id |
| | `district_id` | filter by district id / filter ตาม id ของอำเภอ |
| | `district_name_th` | filter by district name, Thai / filter ตามชื่ออำเภอ (ไทย) |
| | `district_name_en` | filter by district name, English / filter ตามชื่ออำเภอ (อังกฤษ) |
| | `province_id` | filter by province id via district / filter ตาม id ของจังหวัด (ผ่านอำเภอ) |

### Examples / ตัวอย่าง

```bash
# All continents / countries — ทวีป / ประเทศทั้งหมด
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/continents"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/countries"

# Country by ISO code / countries in a continent — ประเทศตามโค้ด ISO / ประเทศในทวีป
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/countries?iso_alpha2=TH"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/countries?continent_name_en=Asia"

# All provinces / จังหวัดทั้งหมด
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/provinces"

# By id — works on every endpoint / filter ตาม id (ใช้ได้ทุก endpoint)
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/provinces?id=1"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/geographies?id=2"

# Provinces in a given region / จังหวัดในภูมิภาคที่ระบุ
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/provinces?geography_id=2"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/provinces?geography_name=ภาคกลาง"

# Districts in a given province / อำเภอในจังหวัดที่ระบุ
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/districts?province_id=1"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/districts?province_name_th=กรุงเทพมหานคร"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/districts?province_name_en=Bangkok"

# Sub-districts in a given district / ตำบลในอำเภอที่ระบุ
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/sub-districts?district_id=1"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/sub-districts?district_name_th=เขตพระนคร"
```

## Running locally & deployment / การรันในเครื่องและการ deploy

The service source is in [`api/`](api/). See [`api/README.md`](api/README.md) for how to run
it locally (Go / Docker) and its full endpoint/parameter reference, and
[`thai-address-api-deployment-guide.md`](thai-address-api-deployment-guide.md) for the
Supabase + Cloud Run deployment walkthrough.

> โค้ดของเซอร์วิสอยู่ใน [`api/`](api/) — ดูวิธีรันในเครื่อง (Go / Docker) และเอกสารอ้างอิง
> endpoint/parameter ฉบับเต็มได้ที่ [`api/README.md`](api/README.md) และดูขั้นตอน deploy บน
> Supabase + Cloud Run ได้ที่ [`thai-address-api-deployment-guide.md`](thai-address-api-deployment-guide.md)

## Support / Tip / สนับสนุนโปรเจกต์

If this API saves you time, buying me a coffee is appreciated 🙏
&nbsp;·&nbsp; ถ้า API นี้ช่วยประหยัดเวลาคุณ เลี้ยงกาแฟกันสักแก้วก็ยินดีมากครับ 🙏

> **This is a voluntary tip, not a purchase.** It is not tied to any API access or
> feature — the API stays completely free for everyone whether you tip or not. Tipping
> buys you no extra rights, quota, or support tier.
>
> **นี่คือการให้ทิปตามความสมัครใจ ไม่ใช่การซื้อสิทธิ์ใด ๆ** — การจ่ายเงินนี้ไม่ผูกกับสิทธิ์
> การใช้งานหรือฟีเจอร์ใด ๆ ทั้งสิ้น API ยังใช้ได้ฟรีเหมือนเดิมสำหรับทุกคน ไม่ว่าจะให้ทิปหรือไม่ก็ตาม
> การให้ทิปไม่ได้ทำให้คุณได้สิทธิ์ โควตา หรือระดับการซัพพอร์ตเพิ่มเติมแต่อย่างใด

[![ko-fi](https://img.shields.io/badge/Support%20me%20on-Ko--fi-FF5E5B?logo=ko-fi&logoColor=white)](https://ko-fi.com/haritpromtarin)
[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20me%20a%20coffee-FFDD00?logo=buymeacoffee&logoColor=black)](https://buymeacoffee.com/haritpromt8)

- Ko-fi: [ko-fi.com/haritpromtarin](https://ko-fi.com/haritpromtarin)
- Buy Me a Coffee: [buymeacoffee.com/haritpromt8](https://buymeacoffee.com/haritpromt8)

## Data source & attribution / แหล่งข้อมูลและ attribution

The Thai administrative division data (province, district, sub-district, geography)
served by this API comes from the
[kongvut/thai-province-data](https://github.com/kongvut/thai-province-data) project by
**Kongvut Sangkla**, used under the MIT License.

> ข้อมูลเขตการปกครองของไทย (จังหวัด อำเภอ ตำบล ภูมิภาค) ที่ให้บริการผ่าน API นี้
> นำมาจากโปรเจกต์ [kongvut/thai-province-data](https://github.com/kongvut/thai-province-data)
> โดย **Kongvut Sangkla** ใช้งานภายใต้สัญญาอนุญาต MIT

> MIT License — Copyright (c) 2025 Kongvut Sangkla

This project does not claim ownership of the underlying dataset. All credit for
collecting and maintaining the data goes to the upstream author. The full upstream
license text is preserved in [`LICENSE`](LICENSE) under the *Third-party data notice*
section, as required by the MIT License. If you use this API, please keep the same
attribution to the original data source.

The continent and country reference data (195 countries, 7 continents) is served from the
same database. Thai names and capitals are curated in-house; ISO alpha-2 / alpha-3 / numeric
codes and English capitals are derived from the public **ISO 3166** standard. Kosovo has no
official ISO numeric code — the user-assigned value `926` / `XKX` is used.

> ข้อมูลอ้างอิงทวีปและประเทศ (195 ประเทศ, 7 ทวีป) ให้บริการจากฐานข้อมูลเดียวกัน — ชื่อไทยและ
> เมืองหลวงจัดทำเอง ส่วนโค้ด ISO alpha-2 / alpha-3 / numeric และชื่อเมืองหลวงภาษาอังกฤษ
> อ้างอิงจากมาตรฐานสากล **ISO 3166** (โคโซโวไม่มีเลข ISO numeric อย่างเป็นทางการ จึงใช้ค่า `926` / `XKX`)

> โปรเจกต์นี้ไม่ได้อ้างความเป็นเจ้าของชุดข้อมูลต้นทาง เครดิตทั้งหมดในการรวบรวมและดูแล
> ข้อมูลเป็นของผู้จัดทำต้นทาง ข้อความสัญญาอนุญาตฉบับเต็มถูกเก็บไว้ในไฟล์ [`LICENSE`](LICENSE)
> ในส่วน *Third-party data notice* ตามที่สัญญาอนุญาต MIT กำหนด หากนำ API นี้ไปใช้
> กรุณาคง attribution ให้แหล่งข้อมูลต้นทางไว้เช่นเดิมด้วย

## License / สัญญาอนุญาต

The code in this repository is released under the [MIT License](LICENSE).
&nbsp;·&nbsp; โค้ดในรีโพนี้เผยแพร่ภายใต้ [สัญญาอนุญาต MIT](LICENSE)

Address data © [Kongvut Sangkla / kongvut/thai-province-data](https://github.com/kongvut/thai-province-data),
also under the MIT License — see [`LICENSE`](LICENSE) for the retained upstream notice.

> ข้อมูลที่อยู่ © [Kongvut Sangkla / kongvut/thai-province-data](https://github.com/kongvut/thai-province-data)
> ภายใต้สัญญาอนุญาต MIT เช่นกัน — ดูข้อความต้นทางที่คงไว้ได้ในไฟล์ [`LICENSE`](LICENSE)

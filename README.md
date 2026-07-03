# Thai Address API

A free, open REST API for Thailand's administrative divisions —
province (จังหวัด), district (อำเภอ), sub-district (ตำบล), and geography/region (ภูมิภาค).

A small [Go](api/) service acts as a proxy in front of [Supabase](https://supabase.com/)
(PostgreSQL + PostgREST), keeping the Supabase `apikey` server-side so it never reaches the client.
Deployed on Google Cloud Run.

Data sourced from [kongvut/thai-province-data](https://github.com/kongvut/thai-province-data) (MIT License).

☕ If this API is useful to you, consider [supporting the project](#donate).

## Endpoints

All endpoints accept only `GET` and return JSON.

| Method | Path | Description |
| ------ | ---- | ----------- |
| GET | `/api/v1/provinces` | Provinces (จังหวัด) |
| GET | `/api/v1/districts` | Districts (อำเภอ/เขต) |
| GET | `/api/v1/sub-districts` | Sub-districts (ตำบล/แขวง) |
| GET | `/api/v1/geographies` | Geographies/regions (ภูมิภาค) |

### Query parameters (filters)

Every endpoint accepts `id` to fetch a single record. Other filters can be combined
(they are AND-ed together). Filters ending in `*_name_*` are resolved by joining the
related table.

| Endpoint | Query param | Description |
| -------- | ----------- | ----------- |
| `/api/v1/provinces` | `id` | filter by province id |
| | `geography_id` | filter by region id |
| | `geography_name` | filter by region name |
| `/api/v1/districts` | `id` | filter by district id |
| | `province_id` | filter by province id |
| | `province_name_th` | filter by province name (Thai) |
| | `province_name_en` | filter by province name (English) |
| `/api/v1/sub-districts` | `id` | filter by sub-district id |
| | `district_id` | filter by district id |
| | `district_name_th` | filter by district name (Thai) |
| | `district_name_en` | filter by district name (English) |
| | `province_id` | filter by province id (via district) |
| `/api/v1/geographies` | `id` | filter by region id |

### Examples

```bash
# All provinces
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/provinces"

# By id (works on every endpoint)
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/provinces?id=1"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/geographies?id=2"

# Provinces in a given region
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/provinces?geography_id=2"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/provinces?geography_name=ภาคกลาง"

# Districts in a given province
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/districts?province_id=1"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/districts?province_name_th=กรุงเทพมหานคร"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/districts?province_name_en=Bangkok"

# Sub-districts in a given district
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/sub-districts?district_id=1"
curl "https://thai-address-api-373901862529.asia-southeast1.run.app/api/v1/sub-districts?district_name_th=เขตพระนคร"
```

## Running locally & deployment

The service source is in [`api/`](api/). See [`api/README.md`](api/README.md) for how to run
it locally (Go / Docker) and its full endpoint/parameter reference, and
[`thai-address-api-deployment-guide.md`](thai-address-api-deployment-guide.md) for the
Supabase + Cloud Run deployment walkthrough.

## Donate

If this API saves you time, buying me a coffee is appreciated 🙏

[![ko-fi](https://img.shields.io/badge/Support%20me%20on-Ko--fi-FF5E5B?logo=ko-fi&logoColor=white)](https://ko-fi.com/haritpromtarin)
[![Buy Me A Coffee](https://img.shields.io/badge/Buy%20me%20a%20coffee-FFDD00?logo=buymeacoffee&logoColor=black)](https://buymeacoffee.com/haritpromt8)

- Ko-fi: [ko-fi.com/haritpromtarin](https://ko-fi.com/haritpromtarin)
- Buy Me a Coffee: [buymeacoffee.com/haritpromt8](https://buymeacoffee.com/haritpromt8)

## License

Released under the [MIT License](LICENSE). Address data © [kongvut/thai-province-data](https://github.com/kongvut/thai-province-data) (MIT).

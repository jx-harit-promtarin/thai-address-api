package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"golang.org/x/time/rate"
	// "github.com/joho/godotenv"
)

var (
	supabaseURL string // e.g. https://onipthuciutummxseytg.supabase.co
	supabaseKey string // anon key, kept server-side only
	allowOrigin string // ค่า Access-Control-Allow-Origin (default "*")
)

// withCORS ครอบ handler เพื่อใส่ CORS headers และตอบ preflight (OPTIONS) ให้ browser
// ปกติ frontend ที่อยู่คนละ origin จะยิงไม่ผ่านถ้าไม่มี header เหล่านี้
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// ตอบ preflight request ทันที ไม่ต้องเข้า handler จริง
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

// proxyQuery ยิง GET ไปยัง Supabase REST API ด้วย query string ที่เตรียมไว้
// (apikey ไม่หลุดไปถึง client) แล้วส่ง response กลับให้ client
func proxyQuery(w http.ResponseWriter, table string, params url.Values) {
	req, err := http.NewRequest("GET", supabaseURL+"/rest/v1/"+table+"?"+params.Encode(), nil)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("upstream error: %v", err)
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// filteredHandler สร้าง handler ที่ proxy ไปยัง table หนึ่ง โดยรองรับ filter 2 แบบ
//   - directFilters : query param -> column บน table ตรง ๆ (เช่น province_id)
//   - relatedTable  : ชื่อ table ที่เชื่อมด้วย FK สำหรับ filter ตามชื่อของ record แม่
//   - nameFilters   : query param -> column ชื่อบน relatedTable (เช่น province_name_th -> name_th)
//
// เมื่อมีการ filter ด้วย nameFilters จะใช้ inner-join (`relatedTable!inner(...)`)
// เพื่อให้ตัด row ที่ไม่ match ออกจากผลลัพธ์ฝั่ง table หลักด้วย
func filteredHandler(table string, directFilters map[string]string, relatedTable string, nameFilters map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query()
		params := url.Values{}
		params.Set("select", "*")

		for param, column := range directFilters {
			if v := q.Get(param); v != "" {
				params.Set(column, "eq."+v)
			}
		}

		// รวบรวม column ชื่อที่ต้อง embed สำหรับ inner-join (เฉพาะที่ถูกส่งมา)
		embedCols := ""
		for param, column := range nameFilters {
			if v := q.Get(param); v != "" {
				if embedCols != "" {
					embedCols += ","
				}
				embedCols += column
				params.Set(relatedTable+"."+column, "eq."+v)
			}
		}
		if embedCols != "" {
			params.Set("select", "*,"+relatedTable+"!inner("+embedCols+")")
		}

		proxyQuery(w, table, params)
	}
}

// envFloat อ่าน env เป็น float64 ถ้าไม่มีหรือ parse ไม่ได้จะใช้ค่า default
func envFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

// envInt อ่าน env เป็น int ถ้าไม่มีหรือ parse ไม่ได้จะใช้ค่า default
func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func main() {
	// โหลดค่าจากไฟล์ .env ถ้ามี (ไม่มีก็ไม่เป็นไร ใช้ env ของระบบแทน)
	// if err := godotenv.Load(); err != nil {
	// 	log.Printf("no .env file loaded: %v", err)
	// }

	supabaseURL = os.Getenv("SUPABASE_URL")
	supabaseKey = os.Getenv("SUPABASE_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_KEY must be set")
	}

	// origin ที่อนุญาตให้เรียกข้าม origin ได้ (default "*" = ทุก origin)
	allowOrigin = os.Getenv("ALLOW_ORIGIN")
	if allowOrigin == "" {
		allowOrigin = "*"
	}

	// rate limit ต่อ IP: ปรับได้ผ่าน env RATE_LIMIT_RPS / RATE_LIMIT_BURST
	// default 5 req/s, burst 10 — เพียงพอสำหรับ address lookup
	rps := envFloat("RATE_LIMIT_RPS", 5)
	burst := envInt("RATE_LIMIT_BURST", 10)
	limiter := newIPRateLimiter(rate.Limit(rps), burst)
	log.Printf("rate limit: %.1f req/s per IP, burst %d", rps, burst)

	// guard ครอบ handler: CORS ชั้นนอกสุด (ตอบ preflight + ให้ 429 มี CORS header
	// เพื่อให้ browser อ่าน error ได้) แล้วค่อยเข้า rate limit
	guard := func(h http.HandlerFunc) http.HandlerFunc {
		return withCORS(limiter.withRateLimit(h))
	}

	http.HandleFunc("/api/v1/provinces", guard(filteredHandler(
		"province",
		map[string]string{"id": "id", "geography_id": "geography_id"},
		"geography",
		map[string]string{"geography_name": "name"},
	)))
	http.HandleFunc("/api/v1/districts", guard(filteredHandler(
		"district",
		map[string]string{"id": "id", "province_id": "province_id"},
		"province",
		map[string]string{"province_name_th": "name_th", "province_name_en": "name_en"},
	)))
	http.HandleFunc("/api/v1/sub-districts", guard(filteredHandler(
		"sub_district",
		map[string]string{"id": "id", "district_id": "district_id"},
		"district",
		map[string]string{
			"district_name_th": "name_th",
			"district_name_en": "name_en",
			"province_id":      "province_id", // filter ผ่าน district.province_id (sub_district ไม่มี column นี้ตรง ๆ)
		},
	)))
	http.HandleFunc("/api/v1/geographies", guard(filteredHandler(
		"geography",
		map[string]string{"id": "id"},
		"",
		nil,
	)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

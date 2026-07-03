package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	// "github.com/joho/godotenv"
)

var (
	supabaseURL string // e.g. https://onipthuciutummxseytg.supabase.co
	supabaseKey string // anon key, kept server-side only
)

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

	http.HandleFunc("/api/v1/provinces", filteredHandler(
		"province",
		map[string]string{"id": "id", "geography_id": "geography_id"},
		"geography",
		map[string]string{"geography_name": "name"},
	))
	http.HandleFunc("/api/v1/districts", filteredHandler(
		"district",
		map[string]string{"id": "id", "province_id": "province_id"},
		"province",
		map[string]string{"province_name_th": "name_th", "province_name_en": "name_en"},
	))
	http.HandleFunc("/api/v1/sub-districts", filteredHandler(
		"sub_district",
		map[string]string{"id": "id", "district_id": "district_id"},
		"district",
		map[string]string{
			"district_name_th": "name_th",
			"district_name_en": "name_en",
			"province_id":      "province_id", // filter ผ่าน district.province_id (sub_district ไม่มี column นี้ตรง ๆ)
		},
	))
	http.HandleFunc("/api/v1/geographies", filteredHandler(
		"geography",
		map[string]string{"id": "id"},
		"",
		nil,
	))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package main

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// visitor เก็บ limiter ของแต่ละ IP พร้อมเวลาที่เห็นล่าสุด (ไว้ล้างตัวเก่าทิ้ง)
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// ipRateLimiter จำกัดจำนวน request ต่อ IP แบบ in-memory
// หมายเหตุ: บน Cloud Run ที่ scale หลาย instance ค่านี้จะเป็น "ต่อ instance"
// ไม่ใช่ global — เพียงพอสำหรับกันการยิงรัวจาก IP เดียว
type ipRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rps      rate.Limit
	burst    int
}

func newIPRateLimiter(rps rate.Limit, burst int) *ipRateLimiter {
	l := &ipRateLimiter{
		visitors: make(map[string]*visitor),
		rps:      rps,
		burst:    burst,
	}
	go l.cleanupLoop()
	return l
}

// getLimiter คืน limiter ของ IP นั้น สร้างใหม่ถ้ายังไม่มี
func (l *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	if v, ok := l.visitors[ip]; ok {
		v.lastSeen = time.Now()
		return v.limiter
	}

	lim := rate.NewLimiter(l.rps, l.burst)
	l.visitors[ip] = &visitor{limiter: lim, lastSeen: time.Now()}
	return lim
}

// cleanupLoop ล้าง IP ที่เงียบไปนานทิ้งเป็นระยะ กัน map โตไม่จำกัด
func (l *ipRateLimiter) cleanupLoop() {
	for {
		time.Sleep(time.Minute)
		l.mu.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

// clientIP ดึง IP จริงของ client บน Cloud Run โดยอ่านจาก X-Forwarded-For ก่อน
// (ตัวแรกใน list คือ client IP จริง) ถ้าไม่มีค่อย fallback ไป RemoteAddr
// หมายเหตุ: header นี้ client ปลอมได้ แต่พอสำหรับกัน abuse ทั่วไป
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip := strings.TrimSpace(strings.Split(xff, ",")[0]); ip != "" {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// withRateLimit ครอบ handler เพื่อจำกัด request ต่อ IP เกินโควตาตอบ 429
func (l *ipRateLimiter) withRateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !l.getLimiter(clientIP(r)).Allow() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

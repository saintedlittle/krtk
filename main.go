package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
)

// Config структура для конфигурации
type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	CORS struct {
		AllowedOrigins []string `yaml:"allowed_origins"`
		AllowedMethods []string `yaml:"allowed_methods"`
		AllowedHeaders []string `yaml:"allowed_headers"`
	} `yaml:"cors"`
	Database struct {
		File string `yaml:"file"`
	} `yaml:"database"`
}

var (
	db     *sql.DB
	config Config
)

func init() {
	// Загрузка конфигурации
	configData, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Инициализация базы данных
	initDB()
}

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Установка заголовков CORS
		if len(config.CORS.AllowedOrigins) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(config.CORS.AllowedOrigins, ", "))
		}
		if len(config.CORS.AllowedMethods) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.CORS.AllowedMethods, ", "))
		}
		if len(config.CORS.AllowedHeaders) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.CORS.AllowedHeaders, ", "))
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", config.Database.File)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Создание таблицы
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_code TEXT NOT NULL UNIQUE,
		original_url TEXT NOT NULL,
		metadata TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP,
		is_permanent BOOLEAN NOT NULL DEFAULT FALSE
	);
	`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized successfully")
}

func captchaHandler(w http.ResponseWriter, r *http.Request) {
	captchaID := captcha.New()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"captcha_id": captchaID,
	})
}

func serveCaptchaImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	captchaID := vars["captchaID"]
	if captchaID == "" {
		http.Error(w, "Invalid CAPTCHA request", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	if err := captcha.WriteImage(&buf, captchaID, captcha.StdWidth, captcha.StdHeight); err != nil {
		http.Error(w, "Failed to generate CAPTCHA", http.StatusInternalServerError)
		return
	}

	img, _, err := image.Decode(&buf)
	if err != nil {
		http.Error(w, "Failed to decode CAPTCHA", http.StatusInternalServerError)
		return
	}

	greenImg := image.NewRGBA(img.Bounds())
	draw.Draw(greenImg, greenImg.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	for x := 0; x < greenImg.Bounds().Dx(); x++ {
		for y := 0; y < greenImg.Bounds().Dy(); y++ {
			original := img.At(x, y)
			r, g, b, a := original.RGBA()
			gray := uint8((r + g + b) / 3 >> 8)
			greenImg.Set(x, y, color.RGBA{0, gray, 0, uint8(a >> 8)})
		}
	}

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, greenImg)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		OriginalUrl  string `json:"originalUrl"`
		CaptchaId    string `json:"captchaId"`
		CaptchaToken string `json:"captchaToken"`
		Metadata     string `json:"metadata"`
		TTL          *int64 `json:"ttl"`
		IsPermanent  bool   `json:"isPermanent"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if !captcha.VerifyString(request.CaptchaId, request.CaptchaToken) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid CAPTCHA"})
		return
	}

	if !isValidURL(request.OriginalUrl) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid URL"})
		return
	}

	if request.TTL != nil && request.IsPermanent {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot specify TTL for permanent links"})
		return
	}

	var expiresAt *time.Time
	if !request.IsPermanent {
		if request.TTL == nil {
			defaultTTL := int64(30 * 24 * 60 * 60) // 30 days in seconds
			request.TTL = &defaultTTL
		}
		expTime := time.Now().Add(time.Duration(*request.TTL) * time.Second)

		// КОНВЕРТАЦИЯ В UTC ДОБАВЛЕНА ЗДЕСЬ
		utcTime := expTime.UTC()
		expiresAt = &utcTime
	}

	shortCode := generateShortCode(6)
	_, err := db.Exec(
		"INSERT INTO links (short_code, original_url, metadata, expires_at, is_permanent) VALUES (?, ?, ?, ?, ?)",
		shortCode,
		request.OriginalUrl,
		request.Metadata,
		expiresAt,
		request.IsPermanent,
	)

	if err != nil {
		log.Printf("Error saving link: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create short URL"})
		return
	}

	fullURL := "https://" + config.Server.Host + "/" + shortCode
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"shortUrl": fullURL,
	})
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		TotalShortened int `json:"totalShortened"`
		ActiveLinks    int `json:"activeLinks"`
		TodayLinks     int `json:"todayLinks"`
	}

	// Get total shortened links
	err := db.QueryRow("SELECT COUNT(*) FROM links").Scan(&stats.TotalShortened)
	if err != nil {
		log.Printf("Error getting total links: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get statistics"})
		return
	}

	// Get active links
	err = db.QueryRow("SELECT COUNT(*) FROM links WHERE is_permanent = 1 OR expires_at > datetime('now')").Scan(&stats.ActiveLinks)
	if err != nil {
		log.Printf("Error getting active links: %v", err)
	}

	// Get today's links
	today := time.Now().Format("2006-01-02")
	err = db.QueryRow("SELECT COUNT(*) FROM links WHERE date(created_at) = ?", today).Scan(&stats.TodayLinks)
	if err != nil {
		log.Printf("Error getting today's links: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func cleanupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowUTC := time.Now().UTC()
	_, err := db.Exec("DELETE FROM links WHERE is_permanent = 0 AND expires_at < ?", nowUTC)
	if err != nil {
		log.Printf("Error cleaning expired links: %v", err)
		http.Error(w, "Failed to clean expired links", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Expired links cleaned successfully"))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	var (
		original    string
		isPermanent bool
		expiresAt   *time.Time
	)
	// Убираем проверку времени из SQL-запроса
	err := db.QueryRow(`
        SELECT original_url, is_permanent, expires_at 
        FROM links 
        WHERE short_code = ?
    `, shortCode).Scan(&original, &isPermanent, &expiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting URL: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Проверяем срок действия в коде
	if !isPermanent {
		if expiresAt == nil {
			http.Error(w, "URL expired", http.StatusGone)
			return
		}
		// Сравниваем в UTC
		if time.Now().UTC().After(*expiresAt) {
			http.Error(w, "URL expired", http.StatusGone)
			return
		}
	}

	http.Redirect(w, r, original, http.StatusTemporaryRedirect)
}

func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func generateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	r := mux.NewRouter()

	// API routes with CORS middleware
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(enableCORS)
	apiRouter.HandleFunc("/captcha", captchaHandler).Methods("GET")
	apiRouter.HandleFunc("/shorten", shortenHandler).Methods("POST")
	apiRouter.HandleFunc("/stats", statsHandler).Methods("GET")
	apiRouter.HandleFunc("/cleanup", cleanupHandler).Methods("POST")

	// Captcha image route
	r.HandleFunc("/captcha/{captchaID}", serveCaptchaImage).Methods("GET")

	// Redirect handler with regex support
	r.HandleFunc("/{shortCode:[a-zA-Z0-9]+}", redirectHandler).Methods("GET")

	// Serve static files
	fs := http.FileServer(http.Dir("./frontend/dist"))
	r.PathPrefix("/").Handler(fs)

	log.Printf("Server starting on %s", config.Server.Port)
	log.Fatal(http.ListenAndServe(config.Server.Port, r))
}

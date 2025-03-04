package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	_ "image/png"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/dchest/captcha"
	_ "github.com/mattn/go-sqlite3"
)

const shortURLLength = 5

var db *sql.DB
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var indexTmpl = template.Must(template.ParseFiles("templates/index.html"))
var redirectTmpl = template.Must(template.ParseFiles("templates/redirect.html"))

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "krtk.db")
	if err != nil {
		log.Fatal(err)
	}

	query := `CREATE TABLE IF NOT EXISTS links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short TEXT UNIQUE,
		original TEXT
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func generateShortURL() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, shortURLLength)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

func saveLink(short, original string) error {
	_, err := db.Exec("INSERT INTO links (short, original) VALUES (?, ?)", short, original)
	return err
}

func getOriginalURL(short string) (string, bool) {
	var original string
	err := db.QueryRow("SELECT original FROM links WHERE short = ?", short).Scan(&original)
	if err != nil {
		return "", false
	}
	return original, true
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	data := struct {
		CaptchaURL string
	}{
		CaptchaURL: "/captcha",
	}
	err := indexTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func captchaHandler(w http.ResponseWriter, r *http.Request) {
	captchaID := captcha.New()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"captcha_id": captchaID,
	})
}

func serveCaptchaImage(w http.ResponseWriter, r *http.Request) {
	captchaID := strings.TrimPrefix(r.URL.Path, "/captcha/")
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
	var req struct {
		Original  string `json:"original"`
		CaptchaID string `json:"captcha_id"`
		Captcha   string `json:"captcha"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !isValidURL(req.Original) {
		http.Error(w, `{"message": "Invalid input URL"}`, http.StatusBadRequest)
		return
	}

	if !captcha.VerifyString(req.CaptchaID, req.Captcha) {
		http.Error(w, `{"message": "Invalid CAPTCHA"}`, http.StatusForbidden)
		return
	}

	short := generateShortURL()
	if err := saveLink(short, req.Original); err != nil {
		http.Error(w, `{"message": "Error saving link"}`, http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"short": fmt.Sprintf("/go/%s", short)}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	short := strings.TrimPrefix(r.URL.Path, "/go/")
	original, exists := getOriginalURL(short)
	if !exists {
		http.NotFound(w, r)
		return
	}

	data := struct {
		EscapedURL string
	}{
		EscapedURL: template.HTMLEscapeString(original),
	}

	err := redirectTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func main() {
	initDB()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/captcha", captchaHandler)
	http.HandleFunc("/captcha/", serveCaptchaImage)
	http.HandleFunc("/go/", redirectHandler)

	fmt.Println("Server running at 0.0.0.0:2137")
	log.Fatal(http.ListenAndServe(":2137", nil))
}

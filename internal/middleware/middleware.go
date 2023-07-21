package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GzipRequest if Content-Encoding set to gzip, than convert body gzip.NewReader
func GzipRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer reader.Close()
		r.Body = reader
		next.ServeHTTP(w, r)
	})

}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write byte
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHandle set Content-Encoding: gzip for Accept-Encoding: gzip
func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

const authorization = "Authorization"

// SetUserCookie set cookie for post method and check setted cookie for req: '/api/user/urls'
func SetUserCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.RequestURI == "/") || (r.Method == http.MethodPost) {
			log.Printf("set cookie")
			cookie := http.Cookie{
				Name:  authorization,
				Value: strconv.Itoa(int(time.Now().Unix())),
			}
			http.SetCookie(w, &cookie)
		}

		if r.RequestURI == "/api/user/urls" && r.Method == http.MethodGet {
			log.Printf("check cookie")
			auth, err := r.Cookie(authorization)
			log.Printf("%v", auth)
			if auth == nil {
				log.Printf("cookie == nil")
				http.Error(w, err.Error(), http.StatusNoContent)
			}
		}

		next.ServeHTTP(w, r)
	})
}

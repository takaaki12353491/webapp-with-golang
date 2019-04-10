package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

// templは１つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() //フラグを解釈します
	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey("S1e2c3u5r3i4t9y1Key")
	gomniauth.WithProviders(
		facebook.New("882590500251-bbl05istqh74u0f4jfddau7f88klri33.apps.googleusercontent.com",
			"sZcmKIPgoMnCYQdeFl9wOUgS",
			"http://localhost:8080/auth/callback/facebook"),
		github.New("882590500251-bbl05istqh74u0f4jfddau7f88klri33.apps.googleusercontent.com",
			"sZcmKIPgoMnCYQdeFl9wOUgS",
			"http://localhost:8080/auth/callback/github"),
		google.New("882590500251-bbl05istqh74u0f4jfddau7f88klri33.apps.googleusercontent.com",
			"sZcmKIPgoMnCYQdeFl9wOUgS",
			"http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// チャットルームを開始します
	go r.run()
	// Webサーバーを起動します
	log.Println("Webサーバーを開始します。ポート：　", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

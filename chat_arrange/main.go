package main

import (
	"../chat_room/trace"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

// templは1つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServerHTTPはHTTPリクエストを処理します
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
	flag.Parse() // フラグを解釈します
	//Gomniauthのセットアップ
	gomniauth.SetSecurityKey("1111")
	gomniauth.WithProviders(
		google.New("827778022927-5rh6jo8mauis03rr7t682ejnkpadqfl7.apps.googleusercontent.com", "PwQKcyBaCGo5biCsAqCfUfK4", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// チャットルームを開始します。
	go r.run()
	// Webサーバを起動します。
	log.Println("Webサーバを開始します。ポート： ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

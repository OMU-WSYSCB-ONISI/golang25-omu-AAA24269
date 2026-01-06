package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key"))
const logFile = "public/logs.json" // データの保存先

const loginHTML = `
<html>
<head>
    <style>body { font-family: sans-serif; padding: 20px; }</style>
</head>
<body>
    <h1>ログイン</h1>
    <form action="/login" method="post">
        ID: <input type="text" name="username"><br>
        PW: <input type="password" name="password"><br>
        <br>
        <input type="submit" value="ログイン">
    </form>
    <p><a href="/bbs">掲示板に戻る</a></p>
</body>
</html>
`

const templateHTML = `
<html>
<head>
    <style>
        p { border-bottom: 1px solid silver; padding: 0.5em; }
        span { font-weight: bold; color: #005500; }
        .del-btn { color: red; font-size: 0.8em; margin-left: 10px; cursor: pointer; }
        form { display: inline; }

        .login-status { text-align: right; background: #eee; padding: 10px; margin-bottom: 20px;}

        /* Markdown用のスタイル */
        blockquote { border-left: 3px solid #ccc; padding-left: 10px; color: #666; }
        code { background-color: #f0f0f0; padding: 2px 4px; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="login-status">
        {{if .IsLoggedIn}}
            <b>ログイン中</b> | <a href="/logout">ログアウト</a>
        {{else}}
            未ログイン | <a href="/login">ログインしてください</a>
        {{end}}
    </div>

    <h1>BBS (Markdown対応)</h1>

    {{if .IsLoggedIn}}
    <div>
        <form action='/write' method='post'>
            名前: <input type='text' name='name' value='管理人'><br>
            本文: <textarea name='body' style='width:30em; height:5em;' placeholder='Markdown記法が使えます'></textarea><br>
            <input type='submit' value='書込'>
        </form>
    </div>
    {{else}}
    <p style="color:red">投稿するにはログインが必要です</p>
    {{end}}

    <hr>

    {{range .Logs}}
    <div style="border-bottom: 1px solid silver; padding: 0.5em;">
        ({{ .ID }}) <span>{{ .Name }}</span>
        <small>{{ formatDate .CTime }}</small>

        {{if $.IsLoggedIn}}
        <form action="/delete" method="post">
            <input type="hidden" name="id" value="{{ .ID }}">
            <input type="submit" value="[削除]" class="del-btn" style="border:none; background:none;">
        </form>
        {{end}}

        <div>
            {{ .Body | parseMarkdown }}
        </div>
    </div>
    {{end}}
</body>
</html>
`

// Log 掲示板に保存するデータを構造体で定義
type Log struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Body  string `json:"body"`
	CTime int64  `json:"ctime"`
}

type PageData struct {
	Logs       []Log
	IsLoggedIn bool // ログイン状態を画面に伝えるフラグ
}

func main() {
	fmt.Printf("Go version: %s\n", runtime.Version())

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/hello", hellohandler)
	http.HandleFunc("/bbs", showHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Println("Launch server...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to launch server: %v", err)
	}
}

func hellohandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "こんにちは ")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.New("login").Parse(loginHTML)
		t.Execute(w, nil)
		return
	}

	// 簡易的な認証
	username := r.FormValue("username")
	password := r.FormValue("password")

	// ID: admin, PASS: password でログイン成功とする
	if username == "admin" && password == "password" {
		// セッションの取得
		session, _ := store.Get(r, "session-name")
		session.Values["authenticated"] = true // 認証済みフラグを立てる
		session.Save(r, w)                     // Cookieに保存

		http.Redirect(w, r, "/bbs", http.StatusFound)
	} else {
		http.Error(w, "認証失敗", http.StatusUnauthorized)
	}
}

//ログアウト処理
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = false // フラグを下ろす
    session.Options.MaxAge = -1 // Cookieを削除
	session.Save(r, w)
	http.Redirect(w, r, "/bbs", http.StatusFound)
}

// ログインチェック
func isLoggedIn(r *http.Request) bool {
	session, _ := store.Get(r, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); ok {
		return auth
	}
	return false
}

// 書き込みログを画面に表示する
func showHandler(w http.ResponseWriter, r *http.Request) {
	// ログを読み出してHTMLを生成

	logs := loadLogs() // データを読み出す
	loggedIn := isLoggedIn(r)//ログイン状態を取得する
	funcMap := template.FuncMap{
		"formatDate": func(t int64) string {
			return time.Unix(t, 0).Format("2006/1/2 15:04")
		},
		"parseMarkdown": func(text string) template.HTML {
			extensions := parser.CommonExtensions | parser.AutoHeadingIDs
			p := parser.NewWithExtensions(extensions)
			htmlFlags := html.CommonFlags | html.HrefTargetBlank
			opts := html.RendererOptions{Flags: htmlFlags}
			renderer := html.NewRenderer(opts)
			byteHTML := markdown.ToHTML([]byte(text), p, renderer)
			return template.HTML(byteHTML)
		},
	}
	// HTML全体を出力
	t, err := template.New("bbs").Funcs(funcMap).Parse(templateHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{
		Logs:       logs,
		IsLoggedIn: loggedIn,
	}
	t.Execute(w, data)
}

// フォームから送信された内容を書き込み
func writeHandler(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm() // フォームを解析
	var log Log
	log.Name = r.FormValue("name")
	log.Body = r.FormValue("body")
	if log.Name == "" {
		log.Name = "名無し"
	}
	if log.Body == "" {
		http.Redirect(w, r, "/bbs", http.StatusFound)
		return
	}
	logs := loadLogs() // 既存のデータを読み出し
	newID := 1
	if len(logs) > 0 {
		newID = logs[len(logs)-1].ID + 1
	}
	log.ID = newID
	log.CTime = time.Now().Unix()
	logs = append(logs, log)
	saveLogs(logs)
	http.Redirect(w, r, "/bbs", http.StatusFound)
}

// 削除機能
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外は許可しない
	if !isLoggedIn(r) {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// フォームからIDを取得して数値に変換
	idStr := r.FormValue("id")
	targetID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// データを読み込む
	logs := loadLogs()
	newLogs := []Log{}

	// 削除対象「以外」を新しいスライスに追加する（フィルタリング）
	for _, log := range logs {
		if log.ID != targetID {
			newLogs = append(newLogs, log)
		}
	}

	// 保存してリダイレクト
	saveLogs(newLogs)
	http.Redirect(w, r, "/bbs", http.StatusFound)
}

// ファイルからログファイルの読み込み
func loadLogs() []Log {
	// ファイルを開く
	text, err := os.ReadFile(logFile)
	if err != nil {
		return make([]Log, 0)
	}
	// JSONをパース
	var logs []Log
	json.Unmarshal([]byte(text), &logs)
	return logs
}

// ログファイルの書き込み
func saveLogs(logs []Log) {
	// JSONにエンコード
	bytes, _ := json.Marshal(logs)
	// ファイルへ書き込む
	os.WriteFile(logFile, bytes, 0644)
}

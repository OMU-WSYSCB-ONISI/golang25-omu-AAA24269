package main

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"runtime"
)

const saveFile = "public/memo.txt" // データファイルの保存先

func main() {
	fmt.Printf("Go version: %s\n", runtime.Version())

	http.HandleFunc("/hello", hellohandler)
	http.HandleFunc("/memo", memo)
	http.HandleFunc("/mwrite", mwrite)

	fmt.Println("Launch server...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to launch server: %v", err)
	}
}

func hellohandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "こんにちは from Codespace !")
}

func memo(w http.ResponseWriter, r *http.Request) {
	// データファイルを開く
	text, err := os.ReadFile(saveFile)
	if err != nil {
		text = []byte("ここにメモを記入してください。")
	}
	// HTMLのフォームを返す
	htmlText := html.EscapeString(string(text))
	s := "<html>" +
	"<style>textarea { width:99%; height:200px; }</style>" +
	"<form method='get' action='/mwrite'>" +
	"<textarea name='text'>" + htmlText + "</textarea>" +
	"<input type='submit' value='保存' /></form></html>"
	w.Write([]byte(s))
}

func mwrite(w http.ResponseWriter, r *http.Request) {
	// 投稿されたフォームを解析
	r.ParseForm()
	if len(r.Form["text"]) == 0 { // 値が書き込まれてない時
		w.Write([]byte("フォームから投稿してください。"))
		return
	}
	text := r.Form["text"][0]
	// データファイルへ書き込む
	os.WriteFile(saveFile, []byte(text), 0644)
	fmt.Println("save: " + text)
	// ルートページへリダイレクトして戻る --- (*4)
	http.Redirect(w, r, "/memo", 301)
}

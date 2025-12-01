package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func main() {
	// Week 04: ここに課題のコードを記述してください
	// 詳細な課題内容はLMSで確認してください

	fmt.Println("Week 04 課題")
    http.HandleFunc("/info", infoHandler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("サーバ起動エラー:", err)
    }
}
	// 以下に実装してください
func infoHandler(w http.ResponseWriter, r *http.Request){
    jst, _ := time.LoadLocation("Asia/Tokyo")
    now := time.Now().In(jst).Format("2006年01月02日 15:04:05")
ua :=r.Header.Get("User-Agent")

fmt.Fprintf(w, "今の時刻は %s で，利用しているブラウザは「%s」ですね", now, ua)
}


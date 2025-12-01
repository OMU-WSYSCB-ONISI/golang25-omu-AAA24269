package main

import (
    "fmt"
    "math/rand"
    "net/http"
    "time"
)

func main() {
	// Week 03: ここに課題のコードを記述してください
	// 詳細な課題内容はLMSで確認してください

	fmt.Println("Week 03 課題")
    http.HandleFunc("/webfortune", fortuneHandler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("サーバ起動エラー:", err)
    }

	// 以下に実装してください

}
func fortuneHandler(w http.ResponseWriter, r *http.Request) {
    fortunes := []string{"大吉", "中吉", "吉", "凶"}
    seed := time.Now().UnixNano()
    rnd := rand.New(rand.NewSource(seed))
    result := fortunes[rnd.Intn(len(fortunes))]

    fmt.Fprintf(w, "今の運勢は %s です", result)
}

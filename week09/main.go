package main
import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
)

func main() {
	fmt.Printf("Go version: %s\n", runtime.Version())

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/hello", hellohandler)
	http.HandleFunc("/fdump", fdump)
	http.HandleFunc("/calx", calculator)

	fmt.Println("Launch server...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to launch server: %v", err)
	}
}

func hellohandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "こんにちは from Glitch !")
}

func fdump(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("errorだよ")
	}
	// フォームはマップとして利用でき以下で内容を確認できる．
	for k, v := range r.Form {
		fmt.Printf("%v : %v\n", k, v)
	}
}

func calculator(w http.ResponseWriter, r *http.Request) {
	// フォームの入力値を取得
	left, _ := strconv.Atoi(r.FormValue("left"))
	right, _ := strconv.Atoi(r.FormValue("right"))
	op := r.FormValue("op") //演算子（ラジオボタンの値）

	// 四則演算の処理
	// 変換エラーがなければ、演算子に従って計算
	var result int
	// 演算子ごとに分岐
	switch op {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		result = left / right
	}

	//HTMLの文字列
	h := `
        <html>
            <head>
                <title>電卓アプリ</title>
            </head>
            <body>
                <form action ="/calx" method="get">
                    左項目：<input type="text" name="left"><br>
                    右項目：<input type="text" name="right"><br>
                    演算子：
                    <input type="radio" name="op" value="+" checked> +
                    <input type="radio" name="op" value="-" > -
                    <input type="radio" name="op" value="*" > ×
                    <input type="radio" name="op" value="/" > ÷
                    <br><input type="submit" name="送信"><hr>

                    [フォームの入力値]<br>
                    左項目：` + strconv.Itoa(left) + `<br>
                    右項目：` + strconv.Itoa(right) + `<br>
                    演算子：` + op + `<br>
                    演算結果：` + strconv.Itoa(result) + `<br>
                </form>
            </body>
        </html>
        `
	// クライアント（ブラウザ）にHTMLを送信
	fmt.Fprint(w, h)
}

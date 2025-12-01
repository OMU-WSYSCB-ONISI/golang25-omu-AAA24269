package main
import (
	"fmt"
	"net/http"
	"runtime"
  "strconv"
  "strings"
)

func main() {
	fmt.Printf("Go version: %s\n", runtime.Version())

	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.HandleFunc("/hello", hellohandler)
        http.HandleFunc("/enq", enqhandler)
	http.HandleFunc("/fdump", fdump)
        http.HandleFunc("/cal00", cal00handler)
        http.HandleFunc("/cal01", calpmhandler)
        http.HandleFunc("/sum", sumhandler)

	fmt.Println("Launch server...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to launch server: %v", err)
	}
}

func hellohandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "こんにちは from Codespace !")
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

func enqhandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("errorだよ")
	}
	// r.FormValue("name")として，フォーム中name欄の値を得る
	fmt.Fprintln(w, r.FormValue("name")+"さん，ご協力ありがとうございます.\n年齢は"+r.FormValue("age")+"で，性別は"+r.FormValue("gend")+"で，出身地は"+r.FormValue("birthplace")+"ですね")
}

func cal00handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("errorだよ")
	}
	price, _ := strconv.Atoi(r.FormValue("price"))
	num, _ := strconv.Atoi(r.FormValue("num"))
	fmt.Fprint(w, "合計金額は ")
	fmt.Fprintln(w, price*num)
}

func calpmhandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("errorだよ")
	}
	x, _ := strconv.Atoi(r.FormValue("x"))
	y, _ := strconv.Atoi(r.FormValue("y"))
	switch r.FormValue("cal0") {
	case "+":
		fmt.Fprintln(w, x+y)
	case "-":
		fmt.Fprintln(w, x-y)
	}
}

func sumhandler(w http.ResponseWriter, r *http.Request) {
	var sum, tt int
	if err := r.ParseForm(); err != nil {
		fmt.Println("errorだよ")
	}
	tokuten := strings.Split(r.FormValue("dd"), ",")
	fmt.Println(tokuten)
	for i := range tokuten {
		tt, _ = strconv.Atoi(tokuten[i])
		sum += tt
	}
	fmt.Fprintln(w, sum)
	fmt.Println(sum)
}

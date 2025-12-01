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
	http.HandleFunc("/bmi", bmihandler)
	http.HandleFunc("/cal02", calallhandler)
	http.HandleFunc("/avgdist", avgdistHandler)

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
func bmihandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("errorだよ")
	}
	weight, _ := strconv.ParseFloat(r.FormValue("weight"), 64)
	heightCm, _ := strconv.ParseFloat(r.FormValue("height"), 64)
	heightM := heightCm / 100.0
	bmi := weight / (heightM * heightM)
	fmt.Fprintf(w, "合計金額は %f\n", bmi)
	fmt.Printf("BMI = %f\n", bmi)
}

func calallhandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, "errorだよ")
		return
	}
	x, _ := strconv.ParseFloat(r.FormValue("x"), 64)
	y, _ := strconv.ParseFloat(r.FormValue("y"), 64)

	switch r.FormValue("cal1") {
	case "+":
		fmt.Fprintf(w, "%.2f\n", x+y)
	case "-":
		fmt.Fprintf(w, "%.2f\n", x-y)
	case "*":
		fmt.Fprintf(w, "%.2f\n", x*y)
	case "/":
		if y != 0 {
			fmt.Fprintf(w, "%.2f\n", x/y)
		} else {
			fmt.Fprintln(w, "0で割ることはできません")
		}
	}
}
func avgdistHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, "errorだよ")
		return
	}
	raw := strings.ReplaceAll(r.FormValue("scores"), " ", "")
	tokens := strings.Split(raw, ",")

	var sum, count int
	dist := make([]int, 11) // 0〜100を10点ごとに分けた11区分

	for _, t := range tokens {
		score, err := strconv.Atoi(strings.TrimSpace(t))
		if err != nil || score < 0 || score > 100 {
			continue
		}
		sum += score
		count++
		dist[score/10]++
	}

	if count == 0 {
		fmt.Fprintln(w, "有効な得点が入力されていません。")
		return
	}

	avg := float64(sum) / float64(count)
	fmt.Fprintf(w, "平均点：%.2f\n", avg)
	fmt.Fprintln(w, "得点分布：")
	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "%2d〜%2d点：%d人\n", i*10, i*10+9, dist[i])
	}
	fmt.Fprintf(w, "100点：%d人\n", dist[10])
}

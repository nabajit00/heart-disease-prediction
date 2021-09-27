package main

import (
	"fmt"
	"heart_disease/service"
	"html/template"
	"log"
	"net/http"
)

const (
	footerPath = "templates/footer.html"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", footerPath)
	if err != nil {
		http.ServeFile(w, r, "./static/internal error.html")
		return
	}
	t.Execute(w, struct {
		Version string
	}{Version: "0.1"})
}

func submit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	output := service.GetPrediction(r)
	if output < 0 {
		output = 0
	}

	t, err := template.ParseFiles("templates/result.html", footerPath)
	if err != nil {
		http.ServeFile(w, r, "./static/internal error.html")
		return
	}
	var o int8 = 0
	if output > 0.5 {
		o = 1
	}
	service.InsertResult(r.RemoteAddr, o)

	t.Execute(w, struct {
		Yes bool
		Per string
	}{Yes: output > 0.5, Per: fmt.Sprintf("%.2f", output*100)})
}

func about(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/about.html", footerPath)
	if err != nil {
		http.ServeFile(w, r, "./static/internal error.html")
		return
	}
	t.Execute(w, nil)
}

func stats(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/stats.html", footerPath)
	if err != nil {
		http.ServeFile(w, r, "./static/internal error.html")
		return
	}

	data, err := service.GetStats()
	if err != nil {
		http.ServeFile(w, r, "./static/internal error.html")
		return
	}

	var yper float32 = float32(data.PredictedCHD) / float32(data.TotalPred)
	var nper float32 = float32(data.PredictedNoCHD) / float32(data.TotalPred)

	t.Execute(w, struct {
		Stat service.Stats_t
		Ylen float32
		Nlen float32
		Yper string
		Nper string
	}{
		Stat: data,
		Ylen: yper * 50,
		Nlen: nper * 50,
		Yper: fmt.Sprintf("%0.2f", yper*100),
		Nper: fmt.Sprintf("%0.2f", nper*100),
	})
}

func main() {
	err := service.LoadParameters()
	if err != nil {
		panic(err)
	}
	service.InitDB()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/about", about)
	http.HandleFunc("/stats", stats)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

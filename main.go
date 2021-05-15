package main

import (
	"embed"
	"fmt"
	"github.com/go-co-op/gocron"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed client/dist
var clientFS embed.FS

var indexToWeekDay = map[string]time.Weekday{
	"0": time.Sunday,
	"1": time.Monday,
	"2": time.Tuesday,
	"3": time.Wednesday,
	"4": time.Thursday,
	"5": time.Friday,
	"6": time.Saturday,
}

func main() {
	distFS, err := fs.Sub(clientFS, "client/dist")
	if err != nil {
		log.Fatal(err)
	}

	location, _ := time.LoadLocation("Europe/Istanbul")
	s := gocron.NewScheduler(location)
	s.StartAsync()

	http.Handle("/", http.FileServer(http.FS(distFS)))

	http.HandleFunc("/api/create-alarm", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Not allowed", http.StatusNotFound)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error while reading imgFile", http.StatusBadRequest)
			return
		}
		_, _ = file, header

		name := r.Form.Get("name")
		if name == "" {
			http.Error(w, "You have to specify name", http.StatusBadRequest)
			return
		}

		gettime := r.Form.Get("time")
		if gettime == "00:00:00" {
			http.Error(w, "You have to specify time", http.StatusBadRequest)
			return
		}

		repeatType := r.Form.Get("repeatType")

		task := func() {
			fmt.Println("tetiklendi")
		}

		s.Every(1)
		s.Day()

		val, ok := indexToWeekDay[repeatType]
		if ok {
			s.Weekday(val)
		} else {
			// special
		}

		s.At(gettime)
		jb, err := s.Do(task)

		fmt.Println(jb)

		//_, err = s.Every(1).Day().Saturday().At(gettime).Do(task)
		if err != nil {
			fmt.Println(err.Error())
		}

		//_, _ = s.Every(1).Hour().StartAt(specificTime).Do(task)
		//fmt.Println(name, gettime, repeatType, file, header)
	})

	log.Println("Starting HTTP server at http://localhost:3000 ...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

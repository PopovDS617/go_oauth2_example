package auth

import (
	"html/template"
	"log"
	"net/http"
)

func HandleGetPersonalAccount(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)

	t := template.New("account.html")

	t, err := t.ParseFiles("templates/account.html")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := struct{ Name string }{"Hello"}

	if err := t.Execute(w, res); err != nil {
		log.Println()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

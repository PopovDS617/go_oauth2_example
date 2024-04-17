package auth

import (
	"golangoauth2example/internal/model"
	"html/template"
	"log"
	"net/http"
)

func HandleGetPersonalAccount(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user := ctx.Value("user")

	user, ok := user.(model.User)

	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	t := template.New("account.html")

	t, err := t.ParseFiles("templates/account.html")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := t.Execute(w, user); err != nil {

		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

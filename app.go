package main

import (
	"database/sql"

	"fyne.io/fyne/app"
	"github.com/axylos/kramer_pager/kramer"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	a := app.New()

	k := kramer.NewKramer(a)

	db, err := sql.Open("sqlite3", "/tmp/fooby.db")
	if err != nil {
		println("got an error")
		println(err)
	}
	resp, _ := db.Query("SELECT * FROM users")

	for resp.Next() {
		var name string
		resp.Scan(&name)
	}
	k.Run()
}

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/kelindar/column"

	bolt "go.etcd.io/bbolt"
)

var boltdb *bolt.DB
var accountdb *column.Collection
var gamesdb *column.Collection

func main() {
	boltdb, _ = bolt.Open("sessions.db", 0600, nil)
	boltdb.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("authkey"))
		tx.CreateBucket([]byte("hashes"))
		return nil
	})

	accountdb = column.NewCollection()
	accountdb.CreateColumn("username", column.ForString())
	accountdb.CreateColumn("password", column.ForString())
	accountdb.CreateColumn("datecreated", column.ForInt64())

	gamesdb = column.NewCollection()
	gamesdb.CreateColumn("creator", column.ForString()) // username
	gamesdb.CreateColumn("id", column.ForString())
	gamesdb.CreateColumn("name", column.ForString())
	gamesdb.CreateColumn("description", column.ForString())
	gamesdb.CreateColumn("datecreated", column.ForInt64())

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"authkey", "id", "name", "desc"},
	}))
	r.Use(middleware.Logger)
	r.Use(putConversion)

	r.Route("/admin", func(r chi.Router) {
		r.Get("/games", listgames)
		r.Get("/users", listusers)
		r.Get("/keyring", listauths)
	})
	r.Route("/api", func(r chi.Router) {
		r.Get("/signin", signin)
		r.Get("/authkey/{hash}/{user}", getauthkey)
		r.Get("/register", register)
		r.Get("/creategame", creategame)
		r.Get("/deletegame", deletegame)
		r.Get("/usergames/{user}", usergames)
		r.Get("/game/{id}", gameid)
		r.Post("/setnamedesc", setnamedesc)
	})

	println("started chi api on :8000")
	http.ListenAndServe(":8000", r)
}

func setnamedesc(w http.ResponseWriter, r *http.Request) {
	var authkey = r.Header.Get("authkey")
	var gid = r.Header.Get("id")
	var name = r.Header.Get("name")
	var desc = r.Header.Get("desc")

	var creator = validauth(authkey)
	if creator == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	gamesdb.Query(func(txn *column.Txn) error {
		txn.WithValue("id", func(v interface{}) bool { return v == gid }).Range(func(i uint32) {
			creatorval, _ := txn.String("creator").Get()

			if creatorval == creator {
				txn.String("name").Set(name)
				txn.String("description").Set(desc)
			}
		})

		return nil
	})
}

func gameid(w http.ResponseWriter, r *http.Request) {
	var id = chi.URLParam(r, "id")
	gamesdb.Query(func(txn *column.Txn) error {
		txn.WithValue("id", func(v interface{}) bool { return v == id }).Range(func(i uint32) {
			creatorval, _ := txn.String("creator").Get()
			nameval, _ := txn.String("name").Get()
			idval, _ := txn.String("id").Get()
			dateval, _ := txn.Int64("datecreated").Get()
			descval, _ := txn.String("description").Get()

			w.Write([]byte(`{
	"creator": "` + creatorval + `",
	"id": "` + idval + `",
	"name": "` + nameval + `",
	"description": "` + descval + `",
	"datecreated": ` + fmt.Sprint(dateval) + `
}`))
		})

		return nil
	})
}

func usergames(w http.ResponseWriter, r *http.Request) {
	var creator = chi.URLParam(r, "user")
	if creator == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var out string = "["
	gamesdb.Query(func(txn *column.Txn) error {
		txn.WithValue("creator", func(v interface{}) bool { return v == creator }).Range(func(i uint32) {
			creatorval, _ := txn.String("creator").Get()
			nameval, _ := txn.String("name").Get()
			idval, _ := txn.String("id").Get()
			dateval, _ := txn.Int64("datecreated").Get()
			descval, _ := txn.String("description").Get()

			if creatorval == creator {
				out += `{
	"creator": "` + creatorval + `",
	"id": "` + idval + `",
	"name": "` + nameval + `",
	"description": "` + descval + `",
	"datecreated": ` + fmt.Sprint(dateval) + `
},`
			}
		})

		return nil
	})
	if out[len(out)-1] == byte(',') {
		out = out[0 : len(out)-1]
	}
	out += "]"
	w.Write([]byte(out))
}

func listusers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Total of " + fmt.Sprint(accountdb.Count()) + " elements\n"))
	accountdb.Query(func(txn *column.Txn) error {
		usercol := txn.String("username")
		passcol := txn.String("password")
		datecol := txn.Int64("datecreated")

		txn.With("username").Range(func(i uint32) {
			userval, _ := usercol.Get()
			passval, _ := passcol.Get()
			dateval, _ := datecol.Get()

			w.Write([]byte(fmt.Sprintf("username: '%s' password:'*' dateval:'%d'\n", userval, passval, dateval)))
		})

		return nil
	})
}

func listauths(w http.ResponseWriter, r *http.Request) {
	boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("authkey"))
		b.ForEach(func(k, v []byte) error {
			w.Write([]byte(fmt.Sprintf("authkey: '%s' user: '%s'\n", k, v)))
			return nil
		})
		return nil
	})
}

func listgames(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Total of " + fmt.Sprint(gamesdb.Count()) + " elements\n"))
	gamesdb.Query(func(txn *column.Txn) error {
		creatorcol := txn.String("creator")
		namecol := txn.String("name")
		idcol := txn.String("id")

		txn.With("creator").Range(func(i uint32) {
			creatorval, _ := creatorcol.Get()
			nameval, _ := namecol.Get()
			idval, _ := idcol.Get()

			w.Write([]byte(fmt.Sprintf("creator: '%s' name:'%s' id:'%s'\n", creatorval, nameval, idval)))
		})

		return nil
	})
}

func deletegame(w http.ResponseWriter, r *http.Request) {
	var authkey = r.Header.Get("authkey")
	var id = r.Header.Get("id")

	var creator = validauth(authkey)
	if creator == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	gamesdb.Query(func(txn *column.Txn) error {
		creatorcol := txn.String("creator")
		idcol := txn.String("id")

		txn.WithValue("creator", func(v interface{}) bool { return v == creator }).Range(func(i uint32) {
			creatorval, _ := creatorcol.Get()
			idval, _ := idcol.Get()

			if creatorval == creator && idval == id {
				fmt.Printf("deleted game index num: %v\n", txn.DeleteAt(i))
			}
		})

		return nil
	})
}

func creategame(w http.ResponseWriter, r *http.Request) {
	var authkey = r.Header.Get("authkey")
	var creator = validauth(authkey)
	var gid = genid()

	println(creator + " created a new game")

	if creator == "" {
		println("cannot create game b/c invalid authkey")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid authkey"}`))
	}

	gamesdb.InsertObject(map[string]interface{}{
		"creator":     creator,
		"id":          gid,
		"name":        "New Game",
		"description": "Newly created game by " + creator,
		"datecreated": time.Now().Unix(),
	})

	w.Write([]byte(`{
	"creator": "` + creator + `",
	"id": "` + gid + `",
	"name": "New Game",
	"description": "Newly created game by ` + creator + `",
	"datecreated": ` + fmt.Sprint(time.Now().Unix()) + `
}`))
}

func getauthkey(w http.ResponseWriter, r *http.Request) {
	var hash string = chi.URLParam(r, "hash")
	var user string = chi.URLParam(r, "user")
	var exists bool = false
	println("generating authkey for " + user)
	boltdb.View(func(tx *bolt.Tx) error {
		res := tx.Bucket([]byte("hashes")).Get([]byte(hash))
		if string(res) == user {
			exists = true
		}
		return nil
	})
	if exists {
		var authkey = createauthkey(user)
		w.Write([]byte(authkey))

		boltdb.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("hashes"))
			b.Delete([]byte(hash))
			return nil
		})
	}
}

func signin(w http.ResponseWriter, r *http.Request) {
	println("signing in user: '" + r.FormValue("user") + "' pass: '" + r.FormValue("pass") + "'")
	var exists bool = false

	accountdb.Query(func(txn *column.Txn) error {
		user := txn.String("username")
		pass := txn.String("password")

		txn.With("username").Range(func(i uint32) {
			userval, _ := user.Get()
			passval, _ := pass.Get()

			if userval == r.FormValue("user") && passval == r.FormValue("pass") {
				exists = true
			}
		})

		return nil
	})

	if exists {
		var hash = createhash(r.FormValue("user"))
		http.Redirect(w, r, "http://localhost:3000/login#"+hash+"#"+r.FormValue("user"), http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, "http://localhost:3000/login#incorrect", http.StatusTemporaryRedirect)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	println("registering user: '" + r.FormValue("user") + "' pass: '" + r.FormValue("pass") + "'")
	var exists = false
	accountdb.Query(func(txn *column.Txn) error {
		txn.With("username").Range(func(i uint32) {
			userval, _ := txn.String("username").Get()

			if userval == r.FormValue("user") {
				exists = true
			}
		})

		return nil
	})

	if !exists {
		accountdb.InsertObject(map[string]interface{}{
			"username":    r.FormValue("user"),
			"password":    r.FormValue("pass"),
			"datecreated": time.Now().Unix(),
		})
		http.Redirect(w, r, "http://localhost:3000/login#registered", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, "http://localhost:3000/login#exists", http.StatusTemporaryRedirect)
	}
}

func createhash(user string) string {
	var hash = randomString(20)
	boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hashes"))
		err := b.Put([]byte(hash), []byte(user))
		if err != nil {
			panic(err)
		}
		return nil
	})
	return hash
}

func createauthkey(user string) string {
	var key = randomString(20)
	boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("authkey"))
		b.Put([]byte(key), []byte(user))
		return nil
	})
	return key
}

func validauth(authkey string) string {
	var user string = ""
	boltdb.View(func(tx *bolt.Tx) error {
		res := tx.Bucket([]byte("authkey")).Get([]byte(authkey))
		user = string(res)
		return nil
	})
	return user
}

func genid() string {
	var val string = randomString(10)
	for idused(val) {
		val = randomString(10)
	}
	return val
}

func idused(id string) bool {
	var exists bool = false
	gamesdb.Query(func(txn *column.Txn) error {
		ids := txn.String("id")

		txn.With("id").Range(func(i uint32) {
			gid, _ := ids.Get()

			if gid == id {
				exists = true
			}
		})

		return nil
	})
	return exists
}

func putConversion(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				r.ParseForm()
				if r.Form["_method"] != nil && r.FormValue("_method") == "PUT" {
					r.Method = "PUT"
				}
			}
			next.ServeHTTP(w, r)
		})
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(len int) string {
	keychars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = keychars[rand.Intn(62)]
	}
	return string(bytes)
}

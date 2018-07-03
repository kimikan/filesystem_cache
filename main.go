package main

import (
	"AmbrWeb/fs"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func main2() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":80") // listen and serve on 0.0.0.0:8080
}

func saveEmail(email string) error {
	db, e := bolt.Open("email.db", 0600, nil)
	if e != nil {
		return e
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte("EMAILS"))
		if e != nil {
			return e
		}

		v, e2 := time.Now().MarshalBinary()
		if e2 != nil {
			return e2
		}

		return b.Put([]byte(email), v)
	})
}

func getEmails() (string, error) {
	db, e := bolt.Open("email.db", 0600, nil)
	if e != nil {
		return "", e
	}
	defer db.Close()

	result := ""
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("EMAILS"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var t time.Time
			ex := t.UnmarshalBinary(v)
			if ex != nil {
				continue
			}

			result += fmt.Sprintf("key=%s, value=%s\n", k, t.String())
		}
		return nil
	})

	return result, nil
}

func main() {
	memfs, err := fs.New("views")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	//router.Static("/views", "./assets")
	//http.Handle("/memfs/", http.StripPrefix("/memfs/", http.FileServer(fs)))

	router.StaticFS("/views", memfs)
	//router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	//router.StaticFS("/views", httpBox)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "views/index.html")
	})

	router.GET("/index_cn.html", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "views/index_CN.html")
	})

	router.POST("/email", func(c *gin.Context) {
		err := saveEmail(c.Request.FormValue("email"))
		if err != nil {
			c.Writer.WriteString(err.Error())
		} else {
			c.Writer.WriteString("{\"msg\":\"email save success!\",\"code\":200}")
		}
	})

	router.GET("/emails", func(c *gin.Context) {
		emails, e := getEmails()
		if e != nil {
			c.Writer.WriteString(err.Error())
		} else {
			c.Writer.WriteString(emails)
		}
	})

	router.Run(":8888")
}

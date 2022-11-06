package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/mastercactapus/quicklink/internal/web"
	"github.com/mastercactapus/quicklink/pkg/store"
)

func main() {
	log.SetFlags(log.Lshortfile)
	addr := flag.String("addr", ":8080", "http service address")
	pg := flag.String("pg", "", "postgres connection string")
	txt := flag.String("txt", "", "text file to use for persistence")
	flag.Parse()

	var s store.Store
	switch {
	case *pg != "":
		var err error
		s, err = store.NewPostgres(*pg)
		if err != nil {
			log.Fatal(err)
		}
	case *txt != "":
		s = store.NewTXTFile(*txt)
	default:
		s = store.NewMemStore()
	}

	srv := &Server{
		Store: s,
		app:   &web.App{Store: s},
	}

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("listening on http://" + l.Addr().String())
	err = http.Serve(l, srv)
	if err != nil {
		log.Fatal(err)
	}
}

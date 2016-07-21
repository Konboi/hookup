package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	PROXY_TIMEUP_SECONDS = 5
)

var (
	configFile = flag.String("c", "config.yml", "config file")
	config     Config
)

func HookHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.String(), "/", "", 1)

	for _, v := range config.Hooks {
		if v.Match(path) {
			req, err := http.NewRequest(r.Method, v.To, r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, err.Error())
				return
			}
			ctx, cancel := context.WithTimeout(req.Context(), PROXY_TIMEUP_SECONDS*time.Second)
			defer cancel()

			req.WithContext(ctx)
			req.Header = r.Header

			log.Printf("localhost/%s to redirect %s \n", path, v.To)

			respCh := make(chan *http.Response)
			errCh := make(chan error)
			//TODO use Context
			go func() {
				client := http.DefaultClient
				resp, err := client.Do(req)
				if err != nil {
					errCh <- err
					return
				}

				respCh <- resp
				return
			}()

			select {
			case resp := <-respCh:
				w.WriteHeader(resp.StatusCode)
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				w.Write(buf.Bytes())
				return
			case err := <-errCh:
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
				return
			case <-ctx.Done():
				err := ctx.Err()
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, fmt.Errorf("request timeout"))
				return
			}
		}
	}

	http.NotFound(w, r)
	return
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	conf, err := NewConfig(*configFile)
	if err != nil {
		log.Println(err)
		return
	}
	config = conf

	http.HandleFunc("/", HookHandler)

	log.Println("start hookup ", fmt.Sprintf(":%d", config.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}

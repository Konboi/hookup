package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHookHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HookHandler))
	defer ts.Close()

	cli := http.DefaultClient

	t.Run("404", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/hogehoge", ts.URL), nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := cli.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusNotFound {
			t.Fatal("redirect error")
		}
	})

	t.Run("proxy github", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/github", ts.URL), nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := cli.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != "fugafuga" {
			t.Fatal("redirect error", string(body))
		}
	})

	t.Run("proxy bitbucket", func(t *testing.T) {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/bitbucket", ts.URL), nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := cli.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != "hogehoge" {
			t.Fatal("redirect error", string(body))
		}
	})

	t.Run("proxy slow", func(t *testing.T) {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/slow", ts.URL), nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := cli.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusInternalServerError {
			t.Fatalf("err not timeout status is $d", res.StatusCode)
		}
	})
}

func TestMain(m *testing.M) {
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "fugafuga")
	}))
	defer ts1.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hogehoge")
	}))
	defer ts2.Close()

	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(8 * time.Second)

		fmt.Fprintf(w, "hogehoge")
	}))
	defer ts3.Close()

	configFile := fmt.Sprintf(`
port: 12345
hooks:
  - from: github
    to: %s
  - from: bitbucket
    to: %s
  - from: slow
    to: %s
`, ts1.URL, ts2.URL, ts3.URL)

	tempfile, err := ioutil.TempFile("", "example.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer os.Remove(tempfile.Name())

	if _, err = tempfile.Write([]byte(configFile)); err != nil {
		fmt.Println(err.Error())
	}

	if err = tempfile.Close(); err != nil {
		fmt.Println(err.Error())
	}

	config, err = NewConfig(tempfile.Name())
	if err != nil {
		fmt.Println(err.Error())
	}

	os.Exit(m.Run())
}

package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := `
port: 1234
hooks:
  - from: fromconfig1
    to: toconfig1
  - from: fromconfig2
    to: toconfig2
`

	tempfile, err := ioutil.TempFile("", "example.yml")
	if err != nil {
		t.Fatal("error create tempfile")
	}
	defer os.Remove(tempfile.Name())

	if _, err = tempfile.Write([]byte(config)); err != nil {
		t.Fatal("error write config yml")
	}

	if err = tempfile.Close(); err != nil {
		t.Fatalf("error close config yml: %s", err.Error())
	}

	conf, err := NewConfig(tempfile.Name())
	if err != nil {
		t.Fatalf("error load config file: %s", err.Error())
	}

	if conf.Port != 1234 {
		t.Fatalf("error load config port is %d", conf.Port)
	}

	if len(conf.Hooks) != 2 {
		t.Fatal("error load config hooks")
	}

	if conf.Hooks[0].From != "fromconfig1" {
		t.Fatal("error load config hooks")
	}

}

func TestMatch(t *testing.T) {
	hook := Hook{
		From: "teststring",
	}

	if !hook.Match("teststring") {
		t.Fatal("fail hook.From match")
	}

	if hook.Match("dummy") {
		t.Fatal("fail match string")
	}

	if hook.Match("test") {
		t.Fatal("fail match string")
	}

	if hook.Match("teststringggg") {
		t.Fatal("fail match string")
	}
}

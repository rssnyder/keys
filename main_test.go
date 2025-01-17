package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	key = time.Now().Unix()
)

func TestNilKey(t *testing.T) {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:3948/%d", key), nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	respData, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 204, resp.StatusCode)
	assert.Equal(t, "", string(respData))
}

func TestNewNamedKey(t *testing.T) {

	value := "foo"

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:3948/%d", key), strings.NewReader(value))
	resp, _ := client.Do(req)
	respData, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, value, string(respData))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://localhost:3948/%d", key), nil)
	resp, _ = client.Do(req)
	respData, _ = io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, value, string(respData))
}

func TestDupKey(t *testing.T) {

	value := "bar"

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:3948/%d", key), strings.NewReader(value))
	resp, _ := client.Do(req)
	respData, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 409, resp.StatusCode)

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://localhost:3948/%d", key), nil)
	resp, _ = client.Do(req)
	respData, _ = io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.NotEqual(t, value, string(respData))
}

func TestNewRandomKey(t *testing.T) {

	value := "foobar"
	key := generateKey(value)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:3948", strings.NewReader(value))
	resp, _ := client.Do(req)
	respData, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, key, string(respData))

	req, _ = http.NewRequest("GET", fmt.Sprintf("http://localhost:3948/%s", string(respData)), nil)
	resp, _ = client.Do(req)
	respData, _ = io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, value, string(respData))
}

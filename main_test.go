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

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:3948/%d", key), nil)
	client := &http.Client{}
	resp, _ := client.Do(req)
	respData, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 204, resp.StatusCode)
	assert.Equal(t, "", string(respData))
}

func TestNewNamedKey(t *testing.T) {

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:3948/%d", key), strings.NewReader("bar"))
	client := &http.Client{}
	resp, _ := client.Do(req)
	respData, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "bar", string(respData))
}

func TestGetNewNamedlKey(t *testing.T) {

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:3948/%d", key), strings.NewReader("bar"))
	client := &http.Client{}
	resp, _ := client.Do(req)
	respData, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "bar", string(respData))
}

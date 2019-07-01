package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	cachedResponses []response
)

type response struct {
	Endpoint string            `json:"endpoint"`
	Regex    string            `json:"regex,omitempty"`
	Method   string            `json:"method"`
	Code     string            `json:"code"`
	Body     string            `json:"body"`
	Headers  map[string]string `json:"headers"`
	Timeout  string            `json:"timeout,omitempty"`
}

type responses struct {
	Responses []response `json:"responses"`
}

func endpointRouter(c *gin.Context) {
	endpoint := c.Param("endpoint")
	switch endpoint {
	case "/canned/upload":
		setCannedResponses(c)
	case "/canned/upload/file":
		setCannedResponseFromFile(c)
	default:
		getResponse(c, endpoint)
	}

}

func setCannedResponseFromFile(c *gin.Context) {
	file, _, err := c.Request.FormFile("responses")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = storeResponses(buf.Bytes())
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
}

func setCannedResponses(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = storeResponses(body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}

func storeResponses(b []byte) error {
	var r responses
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	for _, res := range r.Responses {
		if res.Code == "" {
			return errors.New("invalid empty code specified")
		}
		_, err := strconv.Atoi(res.Code)
		if err != nil {
			return errors.Wrapf(err, "invalid non-numeric status code specified %v", res.Code)
		}
		if res.Endpoint == "" {
			return errors.New("invalid empty endpoint specified")
		}
		if res.Method == "" {
			return errors.New("invalid empty method specified")
		}
		if res.Timeout != "" {
			_, err := strconv.ParseInt(res.Timeout, 10, 64)
			if err != nil {
				return errors.Wrapf(err, "invalid timeout specified %v", res.Timeout)
			}
		}
		updatedEntry := false
		for i, v := range cachedResponses {
			if v.Endpoint == res.Endpoint && v.Method == res.Method && v.Regex == res.Regex {
				cachedResponses[i].Code = res.Code
				cachedResponses[i].Body = res.Body
				cachedResponses[i].Headers = res.Headers
				cachedResponses[i].Timeout = res.Timeout
				updatedEntry = true
			}
		}
		if !updatedEntry {
			cachedResponses = append(cachedResponses, res)
		}
	}
	return nil
}

func getResponse(c *gin.Context, endpoint string) {
	var r response
	for _, v := range cachedResponses {
		if v.Endpoint == endpoint && v.Method == c.Request.Method {
			r = v
			break
		}
		if v.Regex != "" {
			e := regexp.MustCompile(v.Regex).FindString(endpoint)
			if e != "" {
				r = v
				break
			}
		}
	}
	if r.Endpoint == "" {
		err := fmt.Errorf("unable to find a response for %s:%s", c.Request.Method, endpoint)
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	for k, v := range r.Headers {
		c.Writer.Header().Set(k, v)
	}

	code, _ := strconv.Atoi(r.Code)
	if r.Timeout != "" {
		timeout, _ := strconv.ParseInt(r.Timeout, 10, 64)
		log.Printf("waiting %v seconds before responding...", timeout)
		time.Sleep(time.Duration(timeout) * time.Second)
	}
	c.String(code, r.Body)
}

func storeResponsesFromFile(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return errors.Wrapf(err, "failed to read responses file %s", f)
	}
	err = storeResponses(b)
	if err != nil {
		return errors.Wrap(err, "failed to store responses")
	}
	return err
}

func main() {
	log.Println("Canned")

	port := flag.String("port", ":8888", "service port")
	responses := flag.String("responses", "", "JSON list of responses")
	flag.Parse()
	log.Printf("Started with args port%s responses:%s\n", *port, *responses)

	if *responses != "" {
		err := storeResponsesFromFile(*responses)
		if err != nil {
			log.Fatal(err)
		}
	}

	router := gin.Default()
	router.Any("/*endpoint", endpointRouter)
	err := router.Run(*port)
	if err != nil {
		log.Fatalf("failed to start listening on port %s: %v", *port, err)
	}
}

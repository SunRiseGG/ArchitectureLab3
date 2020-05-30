package integration

import (
	. "gopkg.in/check.v1"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const baseAddress = "http://balancer:8090"

var servers = make(map[string]int)
servers["server1:8080"] = 0
servers["server2:8080"] = 0
servers["server3:8080"] = 0

var client1 = http.Client{
	Timeout: 10 * time.Second,
}
var client2 = http.Client{
	Timeout: 10 * time.Second,
}
var client3 = http.Client{
	Timeout: 10 * time.Second,
}


func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s * TestSuite) TestBalancer(c *C) {
	// TODO: Реалізуйте інтеграційний тест для балансувальникка.
	counter  := 0

  for range time.Tick(1 * time.Second) {
	  resp1, err1 := client1.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	  resp2, err2 := client2.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	  resp3, err3 := client3.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		responsesPool = []string{
			resp1.Header.Get("lb-from"),
			resp2.Header.Get("lb-from"),
			resp3.Header.Get("lb-from"),
		}
	  if err1 != nil {
		  c.Error(err1)
	  }
	  if err2 != nil {
		  c.Error(err2)
	  }
	  if err3 != nil {
		  c.Error(err3)
	  }
		for responseHeader := range responsesPool {
			servers[responseHeader]++
			c.Log("response from", responseHeader)
		}
		counter++
		if counter == 20 {
			break
		}
	}
	for _, value := range servers {
		c.Log("Total responses from", key, value)
	}
}


func BenchmarkBalancer(b *testing.B) {
	// TODO: Реалізуйте інтеграційний бенчмарк для балансувальникка.
	for i := 0; i < b.N; i++ {
    _, err1 := client1.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
    _, err2 := client2.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
    _, err3 := client3.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))

		if err1 != nil {
		  b.Error(err1)
	  }
	  if err2 != nil {
		  b.Error(err2)
	  }
	  if err3 != nil {
		  b.Error(err3)
	  }
  }
}

package integration

import (
	. "gopkg.in/check.v1"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const baseAddress = "http://localhost:8090"

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
	count  := 0
  count1 := 0
  count2 := 0
  count3 := 0

  for range time.Tick(1 * time.Second) {
	  resp1, err1 := client1.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	  resp2, err2 := client2.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	  resp3, err3 := client3.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	  if err1 != nil {
		  c.Error(err1)
	  }
	  if err2 != nil {
		  c.Error(err2)
	  }
	  if err3 != nil {
		  c.Error(err3)
	  }
	  switch resp1.Header.Get("lb-from") {
	    case "server1:8080":
			  count1++
		  case "server2:8080":
		  	count2++
		  case "server3:8080":
			  count3++
	  }
		switch resp2.Header.Get("lb-from") {
	    case "server1:8080":
			  count1++
		  case "server2:8080":
		  	count2++
		  case "server3:8080":
			  count3++
	  }
		switch resp3.Header.Get("lb-from") {
	    case "server1:8080":
			  count1++
		  case "server2:8080":
		  	count2++
		  case "server3:8080":
			  count3++
	  }
	  c.Log("response from", resp1.Header.Get("lb-from"))
	  c.Log("response from", resp2.Header.Get("lb-from"))
	  c.Log("response from", resp3.Header.Get("lb-from"))
		count++
		if count == 20 {
			break
		}
	}
	c.Log("Total responses from", count1)
	c.Log("Total responses from", count2)
	c.Log("Total responses from", count3)
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

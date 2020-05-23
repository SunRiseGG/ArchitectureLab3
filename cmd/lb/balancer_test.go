package main

import (
	"testing"
	 . "gopkg.in/check.v1"
 )

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s * TestSuite) TestBalancer(c *C) {
  //serversPool = make(map[int]string)
	//serversPool[0] = "server1:8080"
	//serversPool[1] = "server2:8080"
	//serversPool[2] = "server3:8080"

	for i := 0; i < len(serversPool); i++ {
		res := health(serversPool[i])
		c.Assert(res, Equals, false)
	}

	addr := []string{
      "127.0.0.1:8081",
      "127.0.0.1:8081",
      "127.0.0.1:8082",
      "127.0.0.1:8080",
      "127.0.0.1:8082",
			"127.0.0.1:8081",
			"127.0.0.1:8082",
			"127.0.0.1:8080",
  }

	healthyServers = []int{0, 1, 2}

	for i := 0; i < len(addr); i++ {
		res := chooseServer(addr[i])
		c.Log(serversPool);
		c.Log(healthyServers);
		if i == 0 || i == 1 || i == 5 {
			c.Assert(res, Equals, "server1:8080")
		}
		if i == 2 || i == 4 || i == 6 {
			c.Assert(res, Equals, "server3:8080")
		}
		if i == 3 || i == 7 {
			c.Assert(res, Equals, "server2:8080")
		}
	}

	healthyServers = []int{0, 2}

	for i := 0; i < len(addr); i++ {
		res := chooseServer(addr[i])
		if i == 0 || i == 1 || i == 3 || i == 5 || i == 7{
			c.Assert(res, Equals, "server1:8080")
		}
		if i == 2 || i == 4 || i == 6 {
			c.Assert(res, Equals, "server3:8080")
		}
	}

  healthyServers = []int{0, 1}

	for i := 0; i < len(addr); i++ {
		res := chooseServer(addr[i])
		if i == 3 || i == 7 {
			c.Assert(res, Equals, "server2:8080")
		}
		if i == 0 || i == 1 || i == 5 || i == 2 || i == 4 || i == 6 {
			c.Assert(res, Equals, "server1:8080")
		}
	}

	healthyServers = []int{2}

	for i := 0; i < len(addr); i++ {
		res := chooseServer(addr[i])
		c.Assert(res, Equals, "server3:8080")
	}

  healthyServers = []int{1}

	for i := 0; i < len(addr); i++ {
		res := chooseServer(addr[i])
		c.Assert(res, Equals, "server2:8080")
	}

  healthyServers = []int{0}

	for i := 0; i < len(addr); i++ {
		res := chooseServer(addr[i])
		c.Assert(res, Equals, "server1:8080")
	}

}

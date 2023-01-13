# alkhanif
tools

# Sample usage

```
package main

import (
  "fmt"
  "github.com/azams/alkhanif"
)

var (
  apiKey = "API_KEY"
)

func main {
  response, err := alkhanif.ChatGPT(apiKey, message)
  if err != nil {
    panic(err)
  }
  fmt.Println(response)
}

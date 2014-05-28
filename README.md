# godo

A wrapper for Digitalocean's API in Go

## Usage
```go
package main

import (
	"fmt"

	"github.com/pengux/godo"
)

func main() {
	do := godo.NewClient([CLIENT_ID], [API_KEY])

	fmt.Println(do.GetAllImages())
}
```

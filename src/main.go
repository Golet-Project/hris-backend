package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC)

	fmt.Println(t.Format("15:04"))
}

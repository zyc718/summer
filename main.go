package main

import (
	"fmt"
	"time"
)

func main() {
	current := time.Now()
	fmt.Printf(`Time is %v 
    Beautiful things are happening quietly
                                  --Summer
`, current.Format("2006-01-02 15:04:05"))
	NineTable()

}

func NineTable() {
	start := time.Now()
	tc := time.Since(start)
	fmt.Printf("过程时间是%v\n", tc)
}

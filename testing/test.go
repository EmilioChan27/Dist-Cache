package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	mean := 65
	stddev := 15
	for i := 0; i < 100; i++ {
		fmt.Println(int(math.Abs(rand.NormFloat64()*float64(stddev) + float64(mean))))
	}
}

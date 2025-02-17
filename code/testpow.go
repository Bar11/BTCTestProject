package code

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

//
// hello btc, fuck money 100 years!
// hash
// 2268c832a13bcb433447d25e049bd50f11656cea37c9e5cc584680c3759aa7e9

func POW() {
	start := time.Now()
	for i := 0; i < 1000000000; i++ {
		data := sha256.Sum256([]byte(strconv.Itoa(i)))
		fmt.Printf("%10d,%x\n", i, data)
		fmt.Printf("%s\n", string(data[len(data)-2:]))
		if string(data[len(data)-3:]) == "000" {
			usedTime := time.Since(start)
			fmt.Printf("success used:%d ms", usedTime)
			break
		}
	}
}

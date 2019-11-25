package output

import "fmt"

func ToKafka()  {
	go func(Storage chan string) {
		for {
			fmt.Println(<-Storage)
		}
	}(Storage)
}


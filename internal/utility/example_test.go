package utility

import "fmt"

func ExampleGenerateUniqKey() {
	key := GenerateUniqKey()
	fmt.Println(len(key))

	// Output:
	// 6
}

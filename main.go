package main

import (
	"fmt"
	"github.com/iamGreedy/compiler-assignment-5/cdtgo"
	"golang.org/x/text/encoding/korean"
	"io/ioutil"
	"os"
)

func main() {
	var targets = []string{
		"./CDT/Examples/bubble.mc",
		"./CDT/Examples/comment.mc",
		"./CDT/Examples/ext.mc",
		"./CDT/Examples/factorial.mc",
		"./CDT/Examples/mod.mc",
		//"./test.mc",
		"./CDT/Examples/pal.mc",
		"./CDT/Examples/perfect.mc",
		"./CDT/Examples/prime.mc",
		"./CDT/Examples/retval.mc",
		"./CDT/Examples/test.mc",
	}
	for _, target := range targets {
		fmt.Println("====================================================================================")
		file, err := os.Open(target)
		if err != nil {
			fmt.Println("파일 없음 : ", target)
			continue
		}
		defer file.Close()
		bts, err := ioutil.ReadAll(korean.EUCKR.NewDecoder().Reader(file))
		if err != nil {
			panic(err)
		}
		fmt.Println("테스트 대상 : ", target)
		fmt.Println("------------------------------------------------------------------------------------")
		fmt.Println(string(bts))
		fmt.Println("------------------------------------------------------------------------------------")
		fmt.Println("파서 출력")
		doFile(string(bts))
		fmt.Println()
	}
}

func doFile(src string) {
	tree, err := cdtgo.Parse(src)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	fmt.Println(tree.Format(5))

}

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			//fmt.Println("读取完成")
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}

		result += string((buf[:n]))
	}
	return
}

func getsingle(str, result string) (result2 string) {
	word := regexp.MustCompile(str)
	words := word.FindAllStringSubmatch(result, -1)
	if words != nil{
		result2 = words[0][1]
	}else {
		result2 = "大失败。。。"
	}
	return
}

//爬取单个页面
func SpiderPage(i int) {
	url := "https://hitokoto.cn/?id=" + strconv.Itoa(i)
	result, err := HttpGet(url)
	if err !=nil {
		fmt.Println("error", err)
		return
	}

	str := `<div class="word" id="hitokoto_text">(?s:(.*?))</div>`
	str2 := `<div class="author" id="hitokoto_author">(?s:(.*?))</div>`

	word := getsingle(str, result)
	author := getsingle(str2, result)

	f, err := os.OpenFile("一言.txt", os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(f, word+author)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getyiyan(start, end int)  {
	fmt.Printf("正在整第 %d 到 %d 条一言\n", start, end)

	for i := start; i <= end; i++ {
		SpiderPage(i)
	}

	for i := start; i <= end; i++ {
		fmt.Printf("第 %d 条一言整好了\n", i)
	}
}

func main() {
	var start, end int
	fmt.Println("开始：")
	fmt.Scan(&start)
	fmt.Println("结束：")
	fmt.Scan(&end)

	//爬取
	getyiyan(start, end)
}

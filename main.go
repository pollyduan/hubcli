package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/pollyduan/hubcli/internal"
	flag "github.com/spf13/pflag"
)

var pageSize int
var page int
var help bool
var proxy string

func main() {
	if help || flag.NArg() < 1 {
		ShowHelpMessage()
	} else {
		do()
	}
}

func do() {
	namespace := "library"
	var repo = flag.Arg(0)
	if strings.Contains(repo, "/") {
		kv := strings.Split(os.Args[1], "/")
		namespace = kv[0]
		repo = kv[1]
	}

	endpoint := fmt.Sprintf("https://hub.docker.com/v2/namespaces/%v/repositories/%v/tags?page_size=%v&page=%v", namespace, repo, pageSize, page)
	client := req.C()

	r := client.R()
	resp := new(internal.Tags)
	_, err := r.SetResult(resp).Get(endpoint)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("total tag count: %v, current page: %v/%v\n\n", resp.Count, page, (resp.Count+int64(pageSize)-1)/int64(pageSize))
	for i := 0; i < len(resp.Results); i++ {
		pattern := "%-100v%-30v%-30v\n"
		if i == 0 {
			fmt.Printf(pattern, "Image", "Size", "LastUpdated")
		}
		v := resp.Results[i]
		host := "docker.io"
		if proxy != "" {
			host = "dockerproxy.com"
		}
		image := fmt.Sprintf("%v/%v/%v:%v", host, namespace, repo, v.Name)
		fmt.Printf(pattern, image, v.FullSize, v.LastUpdated)
	}
}

func ShowHelpMessage() {
	fmt.Printf("%s \n\n", "hubcli is a tool which can search docker image from hub.docker.com.")
	fmt.Printf("Usage: %v [namespace/]repository \n\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.IntVarP(&pageSize, "pagesize", "s", 10, "Number of images to get per page. Defaults to 10. Max of 100.")
	flag.IntVarP(&page, "page", "p", 1, "Page number to get. Defaults to 1.")
	flag.BoolVarP(&help, "help", "h", false, "Show help infomation.")
	flag.StringVarP(&proxy, "proxy", "x", "", "Show proxy repo addr. ex: dockerproxy.com")
	flag.Parse()
}

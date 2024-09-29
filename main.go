// Repository : https://github.com/Canop/weblowercaser

package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"golang.org/x/net/html"
	"time"
)

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func GetFixedURL(url string) string {
	if strings.ContainsRune(url, ':') {
		return url
	}
	return strings.ToLower(url)
}

func FixLinks(n *html.Node) (nbLinks int) {
	if n.Type==html.ElementNode {
		for i, attr := range n.Attr {
			if attr.Key == "href" || attr.Key == "src" {
				n.Attr[i].Val = GetFixedURL(attr.Val)
				nbLinks++
				break
			}
		}
	}
	for c:= n.FirstChild; c!=nil; c=c.NextSibling {
		nbLinks += FixLinks(c)
	}
	return
}

func main() {
	from := flag.String("from", "", "source dir")
	to := flag.String("to", "", "dest dir")
	flag.Parse()
	if *from == "" {
		log.Fatal("source path not provided")
	}
	if *to == "" {
		log.Fatal("dest path not provided")
	}
	if *to == *from {
		log.Fatal("source and destination must be distinct")
	}
	nbFiles := 0
	nbHtmlFiles := 0
	nbDirectories := 0
	nbLinks := 0
	start := time.Now()
	filepath.Walk(*from, func(path string, info os.FileInfo, err error) error {
		fi, err := os.Stat(path)
		if err!=nil {
			log.Fatal(err)
		}
		mode := fi.Mode()
		rel, _ := filepath.Rel(*from, path)
		dest := filepath.Join(*to, strings.ToLower(rel))
		if info.IsDir() {
			nbDirectories++
			os.Mkdir(dest, mode)
			if err!=nil {
				log.Fatal(err)
			}
		} else {
			nbFiles++
			switch strings.ToLower(filepath.Ext(rel)) {
			case ".html", ".htm": // note : we don't know how to fix incomplete files
				nbHtmlFiles++
				r, err := os.Open(path)
				if err!=nil {
					log.Fatal(err)
				}
				defer r.Close()
				doc, err := html.Parse(r)
				if err != nil {
					log.Fatal(err)
				} else {
					nbLinks += FixLinks(doc)
					w, err := os.Create(dest)
					if err != nil {
						log.Fatal(err)
					}
					defer w.Close()
					err = html.Render(w, doc)
					if err != nil {
						log.Fatal(err)
					}
				}
			default :
				_, err = CopyFile(dest, path)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		return nil
	})
	log.Println("Finished")
	log.Println("Duration:", time.Since(start))
	log.Println("Directories:", nbDirectories)
	log.Println("Files:", nbFiles)
	log.Println("HTML files:", nbHtmlFiles)
	log.Println("Fixed links:", nbLinks)
}

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
  "strconv"
  "github.com/PuerkitoBio/goquery"
  "golang.org/x/net/html"
  curl "github.com/andelf/go-curl"
  "flag"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
  epPtr := flag.String("episode", "799", "a string")
  flag.Parse()

	//check folder
  if _, err := os.Stat(*epPtr); os.IsNotExist(err) {
  	err := os.Mkdir(*epPtr, 0777)
    check(err)
  }

	//add url to download
  getAllImageUrl("http://mangadoom.co/One-Piece/"+*epPtr+"/", *epPtr)
}

func getAllImageUrl(url string, ep string) {
  //init eas
  easy := curl.EasyInit()
  defer easy.Cleanup()
  //header
  easy = generateHeader(easy)
  //loop page
  for i := 0; i < 100; i++ {
    fmt.Println("page " + strconv.Itoa(i))
    url := url+strconv.Itoa(i)
    ur := getData(easy, url, ep)
    go downloadFromUrl(ur, ep, strconv.Itoa(i+1))

    ts := strings.Split(ur, "/")
    extn := strings.Split(ts[len(ts)-1], ".")
    if (len(extn) == 1) {
      i = 100
    }else{
      if (extn[1] == "png" || extn[1] == "jpg" ){
        //lnjut
      }else{
        i = 100
      }
    }
  }

}

func getData(easy *curl.CURL, url string, ep string) (link string) {
	a := ""
  fmt.Println(url)

  //calback
  fooTest := func (buf []byte, userdata interface{}) bool {

		a = a + string(buf)
    //fmt.Println(a)
		nod, err := html.Parse(strings.NewReader(a))
		check(err)

		doc := goquery.NewDocumentFromNode(nod)

				doc.Find("td").Each(func(i int, s *goquery.Selection){
          s.Eq(0).Each(func(k int, bb *goquery.Selection){
              bb.Find("img").Each(func(l int, cc *goquery.Selection){
                link = cc.AttrOr("src","12345")
              })
          })
        })

		return true
		}

		// page forward to welcome
  	easy.Setopt(curl.OPT_URL, url)
  	easy.Setopt(curl.OPT_HTTPGET, true)
  	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)
  	if err := easy.Perform(); err != nil {
    	println("ERROR: ", err.Error())
  	}
    //fmt.Println("success "+link+" ==\n")

		return link
}


func downloadFromUrl(url string, ep string, fn string){
	tokens := strings.Split(url, "/")
	fileName := ep+"/"+tokens[len(tokens)-1]

  ts := strings.Split(url, "/")
  extn := strings.Split(ts[len(ts)-1], ".")

  if (len(extn) == 1) {
    fmt.Println("\n------------------------------------------------------------\n===== Done Routine =====")
  }else{
    fileName = ep+"/"+fn+"."+extn[1]
    fmt.Println("Downloading", url, "to", fileName)

    // TODO: check file existence first with io.IsExist
  	output, err := os.Create(fileName)
  	if err != nil {
  		fmt.Println("Error while creating", fileName, "-", err)
  		return
  	}
  	defer output.Close()

  	response, err := http.Get(url)
  	if err != nil {
  		fmt.Println("Error while downloading", url, "-", err)
  		return
  	}
  	defer response.Body.Close()

  	_, err = io.Copy(output, response.Body)
  	if err != nil {
  		fmt.Println("Error while downloading", url, "-", err)
  		return
  	}
  	//fmt.Println(n, "bytes downloaded.")
  }
  return
}

func generateHeader(easy *curl.CURL) *curl.CURL {
	easy.Setopt(curl.OPT_VERBOSE, false)
	easy.Setopt(curl.OPT_ENCODING, "Accept-Encoding: gzip, deflate")
	easy.Setopt(curl.OPT_HTTPHEADER, []string{"Accept:text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"})
	easy.Setopt(curl.OPT_HTTPHEADER, []string{"Accept-Language:id-ID,id;q=0.8,en-US;q=0.6,en;q=0.4"})
	easy.Setopt(curl.OPT_USERAGENT, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.132 Safari/537.36")

	return easy
}

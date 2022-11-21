package parsing

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)



func ParseHTML (url string) []string {
	fmt.Println(1)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(2)
		res.Body.Close()
		log.Panic(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println(3)
		res.Body.Close()
		log.Panic("Failed to get html")
	}
	fmt.Println(4)
	z := html.NewTokenizer(res.Body)
	tt := z.Next()

	for tt != html.ErrorToken {
	
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			break
		case tt == html.StartTagToken:
			t := z.Token()
	
			if (t.Data == "p"){
				fmt.Println(11)

				for _, a := range t.Attr{
					fmt.Println(22)

					if a.Key == "class"{
						fmt.Println(a.Val)

						if strings.Contains(a.Val, "exp_title_place"){
							fmt.Println(t)
							fmt.Println(44)

						}
					}
	
					
				}
			}

			
		}
		tt = z.Next()

	}




	// if err != nil {
	// 	res.Body.Close()
	// 	log.Panic(err)
	// }

	// for _, a := range doc.Attr {
	// 	fmt.Println(777)
	// 	fmt.Println(a)
	// }
	res.Body.Close()
	fmt.Printf("doc = %v\n", z)
	fmt.Println(5)
	return nil
}
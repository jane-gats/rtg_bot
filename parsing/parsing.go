package parsing

import (
	//"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Expedition struct {
	Place string
	Name  string
	Link  string
}

func parseTokenizer(tokenizer *html.Tokenizer) []Expedition {
	expeditionList := make([]Expedition, 0, 100)

	tokenType := tokenizer.Next()

	for tokenType != html.ErrorToken {
		if tokenType != html.StartTagToken {
			tokenType = tokenizer.Next()
			continue
		}

		token := tokenizer.Token()
		if token.Data != "p" {
			tokenType = tokenizer.Next()
			continue
		}

		for _, a := range token.Attr {
			if strings.Contains(a.Val, "exp_title_place") {
				expedition := parseExpedition(tokenizer)
				expeditionList = append(expeditionList, expedition)
			}
		}
	}
	return expeditionList
}

func parseExpedition(tokenizer *html.Tokenizer) Expedition {
	var expedition Expedition
	// Разберем параграф с местом и временем экспедиции
	tokenizer.Next()
	token := tokenizer.Token()
	expedition.Place = token.String()

	//Скипаем закрывающий тэг параграфа
	tokenizer.Next()

	//Скипаем пустой текст после закрывающего тэга
	tokenizer.Next()
	//Скипаем открывающий тэг нового параграфа
	tokenizer.Next()

	//Разберем тэг с ссылкой
	tokenizer.Next()
	token = tokenizer.Token()

	for _, a := range token.Attr {
		if a.Key == "href" {
			expedition.Link = a.Val
		}
	}

	//Разберем текст после тэга с ссылкой (название)
	tokenizer.Next()
	token = tokenizer.Token()
	expedition.Name = strings.TrimSpace(token.Data)

	//Перейдем к следующему токену
	tokenizer.Next()
	return expedition
}

func FetchExpeditionsFromUrl(url string) []Expedition {
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		log.Panic(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Panic("Failed to get html")
	}

	tokenizer := html.NewTokenizer(res.Body)

	//tokenType := tokenizer.Next()
	//fmt.Printf("doc = %v\n", tokenizer)
	return parseTokenizer(tokenizer)
}

// func ParseHtmlTag(token html.Token, tokenizer *html.Tokenizer, expedition Expedition) (tokenType html.TokenType) {
//  for _, a := range token.Attr {
//   if strings.Contains(a.Val, "exp_title_place") {
//    tokenType = tokenizer.Next()

//    if tokenType == html.TextToken {
//     token := tokenizer.Token()
//     expedition.Place = token.String()
//    }

//   } else if strings.Contains(a.Val, "exp_title_h") {
//    if tokenType == html.TextToken {
//     token := tokenizer.Token()
//     expedition.Name = token.String()
//    }
//   }
//  }
//  return expedition
// }

// case tokenType == html.TextToken && foundOpenPlaceTag:
//  token := tokenizer.Token()
//  foundOpenPlaceTag = false
//  fmt.Println(token)
//  expedition.Place = token.String()

// case tokenType == html.TextToken && foundOpenHTag:
//  token := tokenizer.Token()
//  foundOpenPlaceTag = false
//  fmt.Println(token)
//  expedition.Name = token.String()
// }

// for tokenType != html.ErrorToken {
//  if tokenType == html.StartTagToken {
//   token := tokenizer.Token()
//   if token.Data == "p" {
//    tokenType = ParseHtmlTag(token, tokenizer, expedition)

//   }
//  }

//  tokenType = tokenizer.Next()
// }

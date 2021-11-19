package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/net/html"
)

type mailContent struct {
	title, body string
}

func ParseAll(w http.ResponseWriter, r *http.Request) {
	err, mailContent1 := crawlApple()
	handleError(err)
	err, mailContent2 := crawlAmazon()
	handleError(err)
	err = sendEmail(mailContent1.title+mailContent2.title, mailContent1.body+mailContent2.body)
	handleError(err)
}

func ParseApple(w http.ResponseWriter, r *http.Request) {
	err, mailContent := crawlApple()
	handleError(err)
	err = sendEmail(mailContent.title, mailContent.body)
	handleError(err)
}

func ParseAmazon(w http.ResponseWriter, r *http.Request) {
	err, mailContent := crawlAmazon()
	handleError(err)
	err = sendEmail(mailContent.title, mailContent.body)
	handleError(err)
}

func crawlApple() (error, mailContent) {
	// targetUrl := "https://www.apple.com/tw/shop/buy-mac/macbook-pro/14-%E5%90%8B-%E5%A4%AA%E7%A9%BA%E7%81%B0%E8%89%B2-8-%E6%A0%B8%E5%BF%83-cpu-14-%E6%A0%B8%E5%BF%83-gpu-512gb#"
	targetUrl := "https://www.amazon.com/Google-WiFi-system-3-Pack-replacement/dp/B01MAW2294?th=1"
	client := &http.Client{}
	resp, err := client.Get(targetUrl)
	handleError(err)
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	file, err := os.Create("./test2.html")
	handleError(err)

	_, err = file.Write(b)
	handleError(err)

	doc, err := html.Parse(resp.Body)
	handleError(err)

	doc = doc.FirstChild.NextSibling.FirstChild

	for doc.Type != html.ElementNode || doc.Data != "body" {
		doc = doc.NextSibling
	}
	doc = doc.FirstChild

	for doc.Type != html.ElementNode || doc.Data != "div" || !checkAttributeExists(doc.Attr, "id", "page") {
		doc = doc.NextSibling
	}

	doc = doc.FirstChild
	for doc.Type != html.ElementNode || doc.Data != "div" || !checkAttributeExists(doc.Attr, "role", "main") {
		doc = doc.NextSibling
	}

	doc = doc.FirstChild
	for doc.Type != html.ElementNode || doc.Data != "div" {
		doc = doc.NextSibling
	}
	doc = doc.FirstChild

	for doc.Type != html.ElementNode || doc.Data != "form" || !checkAttributeExists(doc.Attr, "data-part-number", "MKGP3TA/A") {
		doc = doc.NextSibling
	}

	info := getAttribute(doc.Attr, "data-evar20")
	if strings.Contains(info, "暫未發售") {
		return nil, mailContent{
			"MBP 尚未發售",
			"MBP 尚未發售",
		}

	} else {
		return nil, mailContent{
			"MBP 已發售",
			"MBP 已發售",
		}
	}
	fmt.Println("checked")
	return nil, mailContent{}
}

func crawlAmazon() (error, mailContent) {
	targetUrl := "https://www.amazon.com/-/zh_TW/Google-WiFi-Mesh-%E7%B3%BB%E7%B5%B1-%E8%B7%AF%E7%94%B1%E5%99%A8%E5%8D%87%E7%B4%9A/dp/B08GG9CMLR/ref=sr_1_3?keywords=google+wifi+2020&qid=1637225228&qsid=139-3079787-0386448&sr=8-3&sres=B08GGBZNZ5%2CB08GG9CMLR%2CB07CZPQ9SV%2CB07YMJ57MB%2CB01MDJ0HVG%2CB07YMKD6SM%2CB0842VS9MM%2CB08MVDLVMT%2CB08M38C3LZ%2CB08M478JZF%2CB08P7ZTVLV%2CB07H7RQHF7%2CB08F1YPMS1%2CB081VXXC44%2CB073ZMDMKH%2CB08KTQ1J9H&srpt=NETWORKING_ROUTER"
	c := colly.NewCollector()
	mc := mailContent{}
	c.OnHTML(`div[id=exports_desktop_qualifiedBuybox_tlc_feature_div]`, func(e *colly.HTMLElement) {
		fmt.Println("尚未售完")
		data := strings.Join(e.ChildTexts(`table tr:first-child td`), `:`)
		mc.title = "GOOGLE_WIFI:" + data
		mc.body = "GOOGLE_WIFI:" + data
	})

	c.OnHTML(`div[id=exports_desktop_unqualifiedBuyBox]`, func(e *colly.HTMLElement) {
		fmt.Println("完售")
		mc.title = "GOOGLE_WIFI:" + "完售"
		mc.body = "GOOGLE_WIFI:" + "完售"
	})
	c.Visit(targetUrl)
	return nil, mc
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func checkAttributeExists(attributes []html.Attribute, id, target string) bool {
	for _, a := range attributes {
		if a.Key == id && a.Val == target {
			return true
		}
	}
	return false
}

func getAttribute(attributes []html.Attribute, target string) string {
	for _, a := range attributes {
		if a.Key == target {
			return a.Val
		}
	}
	return ""
}

func sendEmail(subject, body string) error {
	from := mail.NewEmail("gladiator1912", "gladiator1912@gmail.com")
	to := mail.NewEmail("wagaru", "wagaru@hotmail.com")
	message := mail.NewSingleEmail(from, subject, to, subject, "")
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return nil
}

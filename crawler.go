package crawler

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}
type mailContent struct {
	title, body string
}

var mailContents = make([]mailContent, 0)

func ParseAll(ctx context.Context, m PubSubMessage) error {
	fmt.Println("start parsing...")
	defer fmt.Println("finished parsing...")
	err := crawlApple()
	if err != nil {
		return err
	}
	err = crawlAmazon()
	if err != nil {
		return err
	}
	return sendEmail()
}

func ParseApple(ctx context.Context, m PubSubMessage) error {
	fmt.Println("start parsing apple webiste...")
	defer fmt.Println("finished parsing apple website...")
	err := crawlApple()
	if err != nil {
		return err
	}
	return sendEmail()
}

func ParseAmazon(ctx context.Context, m PubSubMessage) error {
	fmt.Println("start parsing amazon webiste...")
	defer fmt.Println("finished parsing amazon website...")
	err := crawlAmazon()
	if err != nil {
		return err
	}
	return sendEmail()
}

func crawlApple() error {
	fmt.Println("start parsing apple...")
	defer fmt.Println("finished parsing apple...")
	targetUrl := "https://www.apple.com/tw/shop/buy-mac/macbook-pro/14-%E5%90%8B-%E5%A4%AA%E7%A9%BA%E7%81%B0%E8%89%B2-8-%E6%A0%B8%E5%BF%83-cpu-14-%E6%A0%B8%E5%BF%83-gpu-512gb#"
	c := colly.NewCollector()
	c.OnHTML(`form[id=configuration-form]`, func(e *colly.HTMLElement) {
		data := e.Attr("data-evar20")
		if strings.Contains(data, "暫未發售") {
			fmt.Println("[APPLE] 暫未發售")
			mailContents = append(mailContents, mailContent{
				"MBP 暫未發售",
				"MBP 暫未發售",
			})
		} else {
			fmt.Println("[APPLE] 已發售")
			mailContents = append(mailContents, mailContent{
				"MBP 已發售",
				"MBP 已發售",
			})
		}
	})
	c.Visit(targetUrl)
	return nil
}

func crawlAmazon() error {
	fmt.Println("start parsing amazon...")
	defer fmt.Println("finished parsing amazon...")
	targetUrl := "https://www.amazon.com/-/zh_TW/Google-WiFi-Mesh-%E7%B3%BB%E7%B5%B1-%E8%B7%AF%E7%94%B1%E5%99%A8%E5%8D%87%E7%B4%9A/dp/B08GG9CMLR/ref=sr_1_3?keywords=google+wifi+2020&qid=1637225228&qsid=139-3079787-0386448&sr=8-3&sres=B08GGBZNZ5%2CB08GG9CMLR%2CB07CZPQ9SV%2CB07YMJ57MB%2CB01MDJ0HVG%2CB07YMKD6SM%2CB0842VS9MM%2CB08MVDLVMT%2CB08M38C3LZ%2CB08M478JZF%2CB08P7ZTVLV%2CB07H7RQHF7%2CB08F1YPMS1%2CB081VXXC44%2CB073ZMDMKH%2CB08KTQ1J9H&srpt=NETWORKING_ROUTER"
	c := colly.NewCollector()
	c.OnHTML(`div[id=exports_desktop_qualifiedBuybox_tlc_feature_div],div[id=qualifiedBuybox_tlc_feature_div]`, func(e *colly.HTMLElement) {
		data := strings.Join(e.ChildTexts(`table tr:first-child td`), `:`)
		fmt.Println("[GOOGLE_WIFI] " + data)
		mailContents = append(mailContents, mailContent{
			"GOOGLE_WIFI:" + data,
			"GOOGLE_WIFI:" + data,
		})
	})

	c.OnHTML(`div[id=exports_desktop_unqualifiedBuyBox],div[id=unqualifiedBuyBox]`, func(e *colly.HTMLElement) {
		fmt.Println("[GOOGLE_WIFI] 完售")
		mailContents = append(mailContents, mailContent{
			"GOOGLE_WIFI:" + "完售",
			"GOOGLE_WIFI:" + "完售",
		})
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	})

	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Printf("%s\n", r.Body)
	// })

	c.Visit(targetUrl)
	return nil
}

func sendEmail() error {
	titles, bodys := make([]string, 0), make([]string, 0)
	for _, v := range mailContents {
		titles = append(titles, v.title)
		bodys = append(bodys, v.body)
	}
	from := mail.NewEmail("gladiator1912", "gladiator1912@gmail.com")
	to := mail.NewEmail("wagaru", "wagaru@hotmail.com")
	message := mail.NewSingleEmail(from, strings.Join(titles, `/`), to, strings.Join(bodys, `/`), "")
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

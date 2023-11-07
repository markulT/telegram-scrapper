package scrapper

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"strconv"
	"strings"
)

type Scrapper interface {
	AuthorizeBot(c context.Context, respch chan struct{})
	ConfirmChannel(c context.Context, channelName string, respchan chan struct{})
	CollectAvgViews(channelName string) error
}

type ChannelAuthorizerRepository interface {
	SaveCode(c Code) error
	SaveChannelAvgViews(ch string, v float64) error
	GetAllChannels() []*Channel
}

type defaultScrapperImpl struct {
	ScrapperContext context.Context
	ChRepo          ChannelAuthorizerRepository
}

//func InitScrapper() {}

func NewDefaultScrapper(chRepo ChannelAuthorizerRepository) Scrapper {
	scrapperContext, _ := chromedp.NewContext(context.Background())
	return &defaultScrapperImpl{ScrapperContext: scrapperContext, ChRepo: chRepo}
}

//var ScrapperContext context.Context

func (s *defaultScrapperImpl) AuthorizeBot(c context.Context, respch chan struct{}) {

	var htmlContent string

	//wg.Add(1)

	// 1 - open browser
	//
	//ctx, cancel := chromedp.NewContext(context.Background())
	//defer cancel()
	err := chromedp.Run(s.ScrapperContext,
		chromedp.Navigate("https://web.telegram.org/a"),
	)
	if err != nil {
		fmt.Println("1")
		fmt.Println(err)
		c.Done()
	}

	err = chromedp.Run(s.ScrapperContext,
		chromedp.WaitVisible(`//button[normalize-space(text())="Log in by phone Number"]`),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("before click")

	fmt.Println("rendered page")

	err = chromedp.Run(s.ScrapperContext,
		chromedp.Click(`//button[normalize-space(text())="Log in by phone Number"]`),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = chromedp.Run(s.ScrapperContext,
		chromedp.WaitVisible(`input[id="sign-in-phone-number"]`),
	)

	fmt.Println("waited for huynia")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("started waiting for clicking on huynia")

	err = chromedp.Run(s.ScrapperContext,
		chromedp.Focus(`input[id="sign-in-phone-code"]`, chromedp.ByQuery),
	)
	fmt.Println("stopped waiting for hyinia being clicked")

	if err != nil {
		fmt.Println("clicking on input code")
		log.Fatal(err)
	}

	fmt.Println("waiting for all countries to be selected")

	//err = chromedp.Run(ctx,
	//	chromedp.Wait(`div.backdrop`),
	//)
	//fmt.Println("stopped waiting")

	var countryNameNodes []*cdp.Node

	if err := chromedp.Run(s.ScrapperContext, chromedp.Nodes(`span.country-name`, &countryNameNodes, chromedp.ByQueryAll)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("selected all countries")

	if err != nil {
		fmt.Println("click on country - error")
		log.Fatal(err)
	}

	fmt.Println("click on country")

	err = chromedp.Run(s.ScrapperContext, chromedp.WaitVisible(`//span[contains(text(), "Ukraine")]`))
	if err != nil {
		log.Fatal(err)
	}

	if err := chromedp.Run(s.ScrapperContext, chromedp.Click(`//span[contains(text(), "Ukraine")]`, chromedp.NodeVisible)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("after click on country")

	if err != nil {
		fmt.Println("2")
		fmt.Println(err)
		c.Done()
	}
	// 2 - input code, that was sent on user's phone number
	fmt.Println("before input")

	// Print the HTML content.

	err = chromedp.Run(s.ScrapperContext,
		chromedp.SendKeys(`input[id="sign-in-phone-number"]`, "987997410"),
	)

	fmt.Println("after input")

	if err != nil {
		fmt.Println("3")
		fmt.Println(err)
		c.Done()
	}

	err = chromedp.Run(s.ScrapperContext,
		chromedp.WaitVisible(`//button[normalize-space(text())="Next"]`),
	)

	if err != nil {
		fmt.Println("fucking button not rendered yet")
		fmt.Println(err)
		c.Done()
	}

	err = chromedp.Run(s.ScrapperContext, chromedp.OuterHTML("html", &htmlContent))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(htmlContent)

	// Use chromedp to locate and click the button with specific attributes and text content.
	err = chromedp.Run(s.ScrapperContext,
		chromedp.Click(`//button[normalize-space(text())="Next"]`),
	)

	if err != nil {
		fmt.Println("4")
		fmt.Println(err)
		c.Done()
	}

	// 3 - wait...

	var activationCode string
	fmt.Println("Input activation code : ")
	fmt.Scanln(&activationCode)

	// 4 - code inputted

	var text string
	if err := chromedp.Run(s.ScrapperContext, chromedp.Text(`h1`, &text, chromedp.NodeVisible)); err != nil {
		fmt.Println("5")
		c.Done()
	}

	// Print the extracted text.
	fmt.Printf("Text from the <h1> element: %s\n", text)

	err = chromedp.Run(s.ScrapperContext,
		chromedp.SendKeys(`input[id="sign-in-code"]`, activationCode),
	)
	if err != nil {
		fmt.Println("6")
		c.Done()
	}
	// 5 - Submit code
	if err = chromedp.Run(s.ScrapperContext, chromedp.Text(`h3`, &text, chromedp.NodeVisible)); err != nil {
		fmt.Println("7")
		c.Done()
	}

	fmt.Printf("Text from the telegram channel title: %s\n", text)

	//wg.Done()

	respch <- struct{}{}

}

func (s *defaultScrapperImpl) ConfirmChannel(c context.Context, channelName string, respch chan struct{}) {

	userAuthorizer := defaultUserAuthorizer{chRepo: s.ChRepo}

	var err error
	err = chromedp.Run(s.ScrapperContext, chromedp.SendKeys(`input[#telegram-search-input]`, channelName))
	if err != nil {
		fmt.Println("Not able to input channel name")
		c.Done()
	}
	err = chromedp.Run(s.ScrapperContext, chromedp.WaitVisible(`.ListItem.search-result`))
	if err != nil {
		fmt.Println("Not able to input channel name")
		c.Done()
	}

	var originalText string
	err = chromedp.Run(s.ScrapperContext, chromedp.Text(`.ChatInfo>.info>.title>.fullName`, &originalText))
	if err != nil {
		fmt.Println("Not able to locate original channel name")
		c.Done()
	}
	err = chromedp.Run(s.ScrapperContext, chromedp.Click(`.ListItem.search-result`))
	if err != nil {
		fmt.Println("Not able to input channel name")
		c.Done()
	}
	err = chromedp.Run(s.ScrapperContext, chromedp.PollFunction(`() => {
            let header = document.querySelectorAll(".ChatInfo>.info>.title>.fullName")[0];
            return header && header.innerText !== '`+originalText+`';
        }`, nil))
	if err != nil {
		fmt.Println("Not able to input channel name")
		c.Done()
	}
	err = chromedp.Run(s.ScrapperContext, chromedp.Click(`.ChatInfo>.info>.title>.fullName`))
	if err != nil {
		fmt.Println("Not able to input channel name")
		c.Done()
	}

	var rawHtmlDescription string

	err = chromedp.Run(s.ScrapperContext, chromedp.OuterHTML(`div.ChatExtra > :first-child`, &rawHtmlDescription))
	adminChannel, err := getAdminChatFromDescription(rawHtmlDescription)
	if err != nil {
		fmt.Println(err)
		c.Done()
	}

	submitCode := userAuthorizer.generateCode(channelName)
	err = userAuthorizer.chRepo.SaveCode(submitCode)
	if err != nil {
		fmt.Println(err)
		c.Done()
	}
	err = userAuthorizer.sendCode(adminChannel, channelName, submitCode)
	if err != nil {
		fmt.Println(err)
		c.Done()
	}

}

func calcAvgViews(n []string) float64 {
	total := 0.0
	for _, item := range n {
		if strings.Contains(item, "K") {
			number, _ := strconv.ParseFloat(strings.Replace(item, "K", "", -1), 64)
			total += number * 1000
		} else {
			number, _ := strconv.ParseFloat(item, 64)
			total += number
		}
	}
	return total / float64(len(n))
}

func (s *defaultScrapperImpl) CollectAvgViews(channelName string) error {
	var err error
	err = chromedp.Run(s.ScrapperContext, chromedp.SendKeys(`input[#telegram-search-input]`, channelName))
	if err != nil {
		fmt.Println("Not able to input channel name")
		return err
	}
	err = chromedp.Run(s.ScrapperContext, chromedp.WaitVisible(`.ListItem.search-result`))
	if err != nil {
		fmt.Println("Not able to input channel name")
		return err
	}

	var originalText string
	err = chromedp.Run(s.ScrapperContext, chromedp.Text(`.ChatInfo>.info>.title>.fullName`, &originalText))
	if err != nil {
		fmt.Println("Not able to locate original channel name")
		return err
	}
	err = chromedp.Run(s.ScrapperContext, chromedp.Click(`.ListItem.search-result`))
	if err != nil {
		fmt.Println("Not able to input channel name")
		return err
	}
	err = chromedp.Run(s.ScrapperContext, chromedp.PollFunction(`() => {
            let header = document.querySelectorAll(".ChatInfo>.info>.title>.fullName")[0];
            return header && header.innerText !== '`+originalText+`';
        }`, nil))
	if err != nil {
		fmt.Println("Not able to input channel name")
		return err
	}

	var viewsList []string

	err = chromedp.Run(s.ScrapperContext, chromedp.Evaluate(`
            let messages = Array.from(document.querySelectorAll('div.Message'));
            messages = messages.slice(-10);
            return messages.map(el => el.innerText);`, &viewsList))
	if err != nil {
		fmt.Println("Not able to input channel name")
		return err
	}

	avgViews := calcAvgViews(viewsList)

	s.ChRepo.SaveChannelAvgViews(channelName, avgViews)
	return nil
}

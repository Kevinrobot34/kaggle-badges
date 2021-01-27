package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

var (
	kaggleURL      = "https://www.kaggle.com"
	sheildsBaseURL = "https://img.shields.io/badge"
	tierColor      = map[string]string{
		"novice":      "4FCB93",
		"contributor": "20BEFF",
		"expert":      "96508E",
		"master":      "F76629",
		"grandmaster": "DDAA17",
	}
	tierLogo = map[string]string{
		"novice":      "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAEV0lEQVR4Ae1YA5A0RxSObdu2rVKc7Z672LZt20m/XsR2ykkxdqZ77nY3tu3l+fqr6Tmud/P/U1Xzqt5h5tXrh++hZ5aIIooooohCRft9KBfgmu5gir42XAZzRd8ZjvOeBxcJtfFdijZjLv1nuMBcMTKFi9wVfzha7h5K441xu0w02Pxf8J0xrER+yrtdwgcbRb/COK5Ev+Ei1/FD9/OS6zm9tJbj0YH+OwMpOKFoYL8P71ssPNHX8k6mxmGzX/reVafKHOiJ5cYzQSX+Pj0YHgdMwVrDCoh81RrxRIwrkbPZ+C00DjALDeC9u5fWryqXlSsiUxZqQ1eOjMwWCgdi74tS4AAwX03OwsiHmhIDoXHAGPS9xX+eu/KganKOpr2Zov9sJ/pjZGRk1pDUwN23MSVshxF5RHqa8R/T4nhv8d9vHLgvRAMstSyzxcm0Xw+OSwyYhzOIPPdnQsk6WbCdKjzEPLkvd0VuQpRzKFgYC9hMfO700BGhm8TdOr6RcaAII2szFRyPtg3fKqHpfbRGpmjYRvpX/G/+HkTBWuMH8d7wR2EpYNtdZBegYdeEQleP3BLP0SaDVon5wGyGuBY5x5WHh8L40z5/cW6uxPge5Irna7TbB9j4PvT3YelH5g9D9C+ZsG0W9+tNLF9N9iCVWgKzIqgFwzfPVOMPyt679Nhq4Bt0XX2H6TzALXCY9yRWmYkTmJ5k2h9gMUV/Ha9S89VfvZ+Zi7n0czDQmKIXZpjBKESu6CSu5RVcCQkYBEXJPTqs4cGnyRkvelHgrjwVOpmSJ6Mdd7xDYWqiTTJgXVOOKxrAVA1aJu6/zR6KLBh9QcstYrmDfsM5w5nunuSanYGJFrtZjA9UGky2r5e7PbFFM5m0Ogcr6vSfl7q8eKztyAf7i22TJlJ3Z2OuSMWU6GXu+OAC41LfQNdaF5eaRnXul46v3c6a7HKk1hYc8Dl5C02t7t/Cxg51G7j8vIwJ3YhOm4kPMAxbKlhg0h405KQTG1SSA1ZxoIXTv46X3Liqzkx8Ncg2o9Pwv3ayN0foNihYe/nI1q4Teiu4kVlsf1mRsR+59E9TOhWVMDuan7CKLsciBiXAZ21nxfWQa4ab0qnF7a1mwB/76u50Hdk3bE/vM5xDlCsxex+bKRWb0ck1dNIFzTvQIzeZ0C2K+2m5RsWhlL53hQkrdM7x4jt39yY2rcRdbmJ76GpGJ2qgpS95qHwTsQ/Hd3wqoENUOOgvrmw/V+JHbKd11ojPIduoTmP891Zn84RPgxP69ZAt0LeNoTcgxehS9iB/MfPie9TTiQw1qBPvy9xNbNXmJI4fBWwHKa3EMIRreX7Du5DZmwJ9tdjRd5/XoetifHOmxDfG0H+4piELlzJSbfjXVqIEXHNF31gdpQBW5tk/RvfPjpI7tG5xFfzi275pfzeZw57kSp6DzyV43ppG1NkrczCX9nRceRZT9DjT4ham5T52Lf//KaKIIoooolElMQKod/Bn2QAAAABJRU5ErkJggg==",
		"contributor": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAE4ElEQVR4AWLAAkYBoD1rAJIlCaJn256qmTnbthk4B862bQfOtm27s5fftm3b/129iKy7WvXs3mxs1KEiMlo1nS+rM19m1hSrsWM+xeVGcG8+wRW2CrveBywXPfBtK2CNoKNJMDsnmGUFi4y48wSzTILethrbRgveVOJYBxIEzWMDESy2CeZZwWlxrrwDF4CdawW93PnLTro7WWIES/3zYiV2iMoAK+jkV94dF9pK7FInJmpQcM/mmBRL3JHShzERTcDS53V1lxQEOzU2L1+BrWmEzptRTLF/FAaQbRiwuvq9submUtRwnk0w3x1viSN4E9wTBO7LJeY+EgT2E3EYIDjXOppUUD1KxEqV4bzUBXmKq+NwIcFe1vm0Z5+CoNhorAi2CJhoeiHF4fHkAEE3AlOAc8g6DcAnmMpcoPMGRgEcwLJW8K7VIFZZom5SawSPqtvM9uDVgNlW8NnpwPIxgZ8SGtBAvPsIJjYwIgLwc4zgJ2ZYd+xJH1eqhDsu4LWT/kxwRvCl5o0sI9oe/J6dseKfzytwiBHcalI8kRfc5uYc6kEyA7eJEW41N2Th5eTGfIrrjeDMXXpg9SzwzR1ZRlCHTXEyadYKbqLeHQWbtDQ5XaNpf5qhCwjmO5mqLjGkHPBZRjhJcu5aGWueut800jO/aMn6acc+WEmTkjJG48IXMuseJlgh+HnZRmSJGtibGJte+RSv03L/A3ccZwU/mt/xjTsfyy+hz1jXf9KKeeQXBZqpV72hXaMvYWYMV4LGMCDrKXrJeHZhx1WBU8oFX6jAAXTH5urlXMZIY6vwo/o428CuGfXMJ1wJVVbTCr3ELy3WK6htzIDpurLTbIqjMwJ8X+V1OJlXZmPC941oFb2+JGbk5ztj7SbpdRBWdvMW6krM31Gwxt8FTxIwgoXl6w2+gDtOc355REYDs7dN4OdOboV+Ynzz9WZ8gZzgV9/TOvk0wxe/ClihAwOuTAMqWqxX0LGhhQlODJLKXFuB5xq8JMUDChy5pHwWUr1H8n3N1WsFM5mdm+Ljn+xfkb5QA+wjm+JtKxgYls1kgnJXn7+3Cd42grH63ky9NsUCK/iuyReyLLC/Y2AuOzOWX9OH4FWXFUzSDI9GRRe2JGkQUCHF3Zpc2CouIUNpAM2jH+bKrCRD8GFdZQUP8h4p1aQOMHU7DFoLPdoiPflqbFUQnGcE91nBXfkUZ28j2KCscjgDvC8Kt6/EpvkKnMDizem83wou4s5fq5bazTWisWIvC3zbjQwjLJuZFL8zMDUgx7tjJdkml+C7CMBnGKFJKdzIDWnXSDTgG/j0Oz6T+8KMgUhK5n13XBhyvZO31L3iGLYCN9ZrM3/IC44qtsdaLNONQOrQb4IHogG/Yy3WC8DPt4L3mkiSTwauNptsE8fqpziYhZgCm5CVnY3203QvUmUUBvBPPOs7KsGPJVrHD7URWhDN9jq3XXz9xB62xNd62+/eMUnFsrF7proQqXJMlgvl1IVYHjDDxhHEgk1IjZ5hTIqnmzD0MdKrZ6ptqpCPaXv9Nt8aapBKUbAf20B2UibFG8Guwmyb4uFowPtk5nvVoPSdTqMsjwpeE9pwzcBxjf1rsar9HZ+GTU8ouuPXnV9lmZgHN524X+R7Bz12LAjOWuafNO4DlttRsAaPy/w/mh5/AG4MRHMtYvLlAAAAAElFTkSuQmCC",
		"expert":      "data:image/png;base64,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",
		"master":      "data:image/png;base64,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",
		"grandmaster": "data:image/png;base64,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",
	}
	badgeURL = "https://img.shields.io/badge/Kaggle-grandmaster-0?color=DDAA17&style=flat-square"
)

func getBadgeURL(userName, tier, rank, styleOptions, logoOptions string) string {
	u, _ := url.Parse(sheildsBaseURL)
	// if err != nil {
	//     return err
	// }
	color := tierColor[tier]

	rightStr := tier + ", " + rank
	u.Path = path.Join(u.Path, "Kaggle-"+rightStr+"-0")

	// set query parameters
	params := url.Values{}
	params.Add("color", color)
	if logoOptions == "On" {
		logo := tierLogo[tier]
		params.Add("logo", logo)
	}
	params.Add("style", styleOptions)

	// encoding
	u.RawQuery = params.Encode()
	fmt.Printf("%v\n", u)
	return u.String()
}

func getCurrentRank(userName string) (string, string) {
	if userName == "" {
		return "novice", ""
	}

	url := fmt.Sprintf("%v/%v", kaggleURL, userName)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "novice", ""
	}
	defer res.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(res.Body)

	kaggleComponentText := doc.Find(".kaggle-component").Text()
	// get rank
	reRank := regexp.MustCompile(`"rankCurrent":([^,]+),`)
	rank := reRank.FindStringSubmatch(kaggleComponentText)[1]
	// get tier
	reTier := regexp.MustCompile(`"tier":"([^,]+)",`)
	tier := reTier.FindStringSubmatch(kaggleComponentText)[1]

	// check
	fmt.Println(tier, rank)
	fmt.Printf("%T %v\n", tier, tier)
	fmt.Printf("%T %v\n", rank, rank)
	return tier, rank
}

func main() {
	engine := gin.Default()

	engine.LoadHTMLGlob("templates/*")

	// localhost:3000/
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	engine.GET("/main", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"badgeUrl": badgeURL, "title": "TODO List"})
	})

	engine.POST("/generate", func(c *gin.Context) {
		userName := c.PostForm("username")
		styleOptions := c.PostForm("style_options")
		logoOptions := c.PostForm("logo_options")
		tier, rank := getCurrentRank(userName)

		badgeURL = getBadgeURL(userName, tier, rank, styleOptions, logoOptions)
		fmt.Println(badgeURL)
		fmt.Println("styleOptions", styleOptions)
		fmt.Println("logoOptions", logoOptions)
		c.Redirect(http.StatusMovedPermanently, "/main")
	})

	// localhost:3000/user/{userName}
	engine.GET("/user/:userName", func(c *gin.Context) {
		userName := c.Param("userName")
		tier, rank := getCurrentRank(userName)
		badgeURL := getBadgeURL(userName, tier, rank, "flat-square", "Off")
		c.JSON(http.StatusOK, gin.H{
			"message":  "user info",
			"userName": userName,
			"tier":     tier,
			"rank":     rank,
			"badgeUrl": badgeURL,
		})
	})
	engine.Run(":3000")
}

# yandex-market

Client for yandex market [API](https://yandex.com/dev/market/partner/doc/dg/reference/all-methods-docpage/).

Currently support only [some](https://pkg.go.dev/github.com/KazanExpress/yandex-market) of the API methods.

## Install

```bash
go get github.com/KazanExpress/yandex-market
```

## Usage

```golang
package main

import (
    "fmt"

    "github.com/KazanExpress/yandex-market/pkg/market/client"
)

func main() {

    c := client.NewYandexMarketClient(
        client.WithOAuth(os.Getenv("OAUTH_TOKEN"), os.Getenv("OAUTH_CLIENT_ID")),
    )

    feeds, err := c.ListFeeds(context.Background(), campaignID)
    if err != nil {
        t.Fatal(err)
    }

    for _, feed := range feeds {
        fmt.Sprintf("there is feed %d with url %s", feed.ID, feed.URL)
    }
}
```

## Yandex Auth

- How to get oauth token [[RU](https://yandex.ru/dev/oauth/doc/dg/tasks/get-oauth-token.html)], [[ENG](https://yandex.com/dev/oauth/doc/dg/tasks/get-oauth-token.html)]

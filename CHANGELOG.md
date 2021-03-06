# Changelog

## v0.4.0

- Translate all godocs to english.

- Update client methods signature to pass `context.Context`.

- Update methods to use option pattern.

- Refactor code.

## v0.3.2

Small fixes.

## v0.3.1

Remove debug logging.

## v0.3.0

- `ExploreOffers(campaignID, options)` - Поиск по маркету внутри магазина.

## v0.2.0

- `client.SetOfferPrices(campaignID, offers)` - Установить цену предложений
- `client.GetOfferPrices(campaignID, ...)` - Получить установленные цены
- `client.DeleteAllOffersPrices(campaignID)` - Удалить все установленные цены
- `client.HideOffers(campaignID, offers)` - Скрыть предложения в маркете
- `client.UnhideOffers(campaignID, offers)` - Возобновить показ скрытых предложений
- `client.GetHiddenOffers(campaignID, ...)` - Список скрытых с помощью API предложений

## v0.1.1

- по умолчанию `APIEndpoint = "https://api.partner.market.yandex.ru"`

## v0.1.0

- `client.ListFeeds(campaignID)` - Список прайс-листов, размещенных на Яндекс.Маркете для магазина
- `client.RefreshFeed(campaignID, feedID)` - Позволяет сообщить, что магазин обновил прайс-лист

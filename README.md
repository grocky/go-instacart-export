# Instacart CSV Export

Export your order history from Instacart into a CSV. I know, exciting stuff!

If you're like me and use text based accounting to track your expenses, having a CSV to get insight into how much
has been spent on various categories make it easy to import.

## CSV Fields

| Field     | Description                             |
|-----------|-----------------------------------------|
| id        | The order ID                            |
| status    | The order status                        |
| total     | The order total                         |
| createdAt | The order creation date (YYYY-MM-DD)    |
| retailers | A pipe delimited list of retailer names |
| numItems  | The number of items in the order        |

## Sample output

![csv header image](./docs/csv-header.png)

## Setup

This tool was created quick and dirty. It's rough, really rough. It works for my purposes, but your mileage may vary.

You'll need to log into the Instacart web app and get your session token. The token is in the `_instacart_session_id` cookie.

![dev tools screenshot](./docs/devtools-token.png)

Install the binary with the following

```shell
go get -u github.com/grocky/go-instacart-export/...
```

## Usage

Now you can invoke this binary and it will fetch all of your orders and place them into a CSV.

```shell
INSTACART_SESSION_TOKEN=<your token> instacart-export
```

By default, `instacart-export` will only export the first 10 pages of orders (Instacart has strict rate limiting on
this endpoint). You can paginate through orders with the following flags.

| Flag  | Description                                |
|-------|--------------------------------------------|
| start | The first page of order results to request |
| end   | The last page of order results to request  |

> **Note**: I wasn't able to figure out if the parameter for the page size for this endpoint; if anyone figures it out
> feel free to open an issue or pull request!

## License

[MIT © Rocky Gray](LICENSE)

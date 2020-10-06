package conn

import (
	"context"
	"time"

	"github.com/olivere/elastic/v7"
)

type EsClient struct {
	Conn *elastic.Client
}

func NewESClient(ESUrl string) (ec *EsClient, err error) {
	ec = &EsClient{}
	ec.Conn, err = elastic.NewClient(
		elastic.SetURL(ESUrl),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(30*time.Second),
		elastic.SetGzip(true),
	)

	if err != nil {
		return
	}

	_, _, err = ec.Conn.Ping(ESUrl).Do(context.Background())
	return
}

package es

import (
	"context"
	"errors"
	"fmt"

	"github.com/hongxincn/promexp/node2es/config"
	"gopkg.in/olivere/elastic.v5"
)

var Client *EsClient

type EsClient struct {
	client *elastic.Client
	bs     *elastic.BulkService
}

func NewEs5Client() error {
	var esclient *elastic.Client
	configuration := config.Config
	if len(configuration.Es.Urls) == 0 {
		return errors.New("no es configuration")
	}
	var err error
	decyptedPassword, err := config.GetDecryptedPassword(configuration.Es.Password)
	if err != nil {
		return err
	}
	if decyptedPassword != "" {
		esclient, err = elastic.NewClient(elastic.SetURL(configuration.Es.Urls...),
			elastic.SetBasicAuth(configuration.Es.Username, decyptedPassword), elastic.SetSniff(false))
		if err != nil {
			return err
		}
	} else {
		esclient, err = elastic.NewClient(elastic.SetURL(configuration.Es.Urls...), elastic.SetSniff(false))
		if err != nil {
			return err
		}
	}
	Client = &EsClient{client: esclient}
	return nil
}

func (ec *EsClient) AddBulkRequest(msg string) {
	if ec.bs == nil {
		ec.bs = elastic.NewBulkService(ec.client)
	}

	index := config.Config.Es.Index

	bir := elastic.NewBulkIndexRequest().Index(index).Type("doc").Doc(
		string(msg))

	ec.bs.Add(bir)
}

func (ec *EsClient) SubmitBulkRequest() {
	ctx := context.Background()
	_, err := ec.bs.Do(ctx)
	if err != nil {
		fmt.Println(err)
		ec.bs = elastic.NewBulkService(ec.client)
	}
}

package mws

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type MWS struct {
	URL string
}

func (m MWS) QueryPage(ctx context.Context, page int, per_page int, query string) (chan int64, error) {
	q := RawMWSQuery{
		From: int64(page * per_page),
		Size: int64(per_page),

		ReturnTotal:  false,
		OutputFormat: "mws-ids",

		Expressions: []MWSExpression{
			{Term: query},
		},
	}
	xml, err := xml.Marshal(q)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to marshal xml")
	}

	// prepare a new request
	req, err := http.NewRequestWithContext(ctx, "POST", m.URL, bytes.NewReader(xml))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to prepare request")
	}
	req.Header.Set("Content-Type", "text/xml")

	// do it!
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to send request")
	}
	defer resp.Body.Close()

	// decode the body
	var result Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Wrapf(err, "unable to decode result")
	}

	// we have no results
	if len(result.MathWebSearchIDs) == 0 {
		return nil, nil
	}

	// and return the channel
	c := make(chan int64, 2*per_page)
	go func() {
		defer close(c)
		for _, id := range result.MathWebSearchIDs {
			c <- id
		}
	}()
	return c, nil
}

func (m MWS) Query(ctx context.Context, query string, per_page int) chan int64 {
	res := make(chan int64)
	go func() {
		defer close(res)
		page := 0
		for {
			if ctx.Err() != nil {
				fmt.Println("context error")
				break
			}

			// request the current page
			c, err := m.QueryPage(ctx, page, per_page, query)
			if err != nil || c == nil {
				break // we're done!
			}

			// pipe over the id!
			for id := range c {
				res <- id
			}

			// and increment for the next page
			page++
		}
	}()
	return res
}

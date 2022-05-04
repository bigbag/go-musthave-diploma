package accrual

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

const RepeatTimeout int64 = 1

type Repository struct {
	ctx context.Context
	l   logrus.FieldLogger
	url string
}

func NewRepository(
	ctx context.Context,
	l logrus.FieldLogger,
	url string,
) *Repository {
	return &Repository{ctx: ctx, l: l, url: url}
}

func (r *Repository) Get(orderID string) (*Response, error) {
	r.l.Info("accrual: get info for id ", orderID)
	resp, err := http.Get(r.url + "/api/orders/" + orderID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r.l.Info("accrual: response status code ", resp.StatusCode)

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		timeout, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			return nil, err
		}

		return &Response{
			OrderID: orderID,
			Amount:  0,
			Status:  PROCESSING,
			IsFinal: false,
			Timeout: int64(timeout),
		}, nil

	case http.StatusOK:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		info := &AccrualInfo{}
		if err := json.Unmarshal(body, &info); err != nil {
			return nil, err
		}

		switch info.Status {
		case REGISTERED:
			return &Response{
				OrderID: orderID,
				Amount:  0,
				Status:  NEW,
				IsFinal: false,
				Timeout: RepeatTimeout,
			}, nil

		case INVALID:
			return &Response{
				OrderID: orderID,
				Amount:  0,
				Status:  INVALID,
				IsFinal: true,
			}, nil

		case PROCESSING:
			return &Response{
				OrderID: orderID,
				Amount:  0,
				Status:  PROCESSING,
				IsFinal: false,
				Timeout: RepeatTimeout,
			}, nil

		case PROCESSED:
			return &Response{
				OrderID: orderID,
				Amount:  info.Amount,
				Status:  PROCESSED,
				IsFinal: true,
			}, nil

		default:
			return nil, nil
		}
	default:
		return &Response{
			OrderID: orderID,
			Amount:  0,
			Status:  PROCESSING,
			IsFinal: false,
			Timeout: RepeatTimeout,
		}, nil
	}
}

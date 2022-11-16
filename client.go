package kocto

import (
	"errors"

	"github.com/imroc/req/v3"
)

func ResponseError(resp *req.Response, err error) error {
	if err != nil {
		return err
	}

	if resp.IsError() {
		body, err := resp.ToString()
		if err != nil {
			return err
		}

		return errors.New(body)
	}

	return nil
}

package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Post[T, R any](ctx context.Context, lgr log.Logger, urlfull string, headers map[string]string, queryParameters url.Values, body T, responseType R) (R, error) {
	client := http.Client{}

	u, err := url.Parse(urlfull)
	if err != nil {
		return responseType, err
	}

	// add query parameters

	q := u.Query()
	for k, v := range queryParameters {
		// this depends on the type of api, you may need to do it for each of v
		q.Set(k, strings.Join(v, ","))
	}
	// set the query to the encoded parameters
	u.RawQuery = q.Encode()

	b, err := toJSON(body)
	if err != nil {
		return responseType, err
	}
	byteReader := bytes.NewReader(b)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), byteReader)
	if err != nil {
		return responseType, fmt.Errorf("could not prepare http request: %w", err)
	}
	// add headers
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// finally, do the request
	res, err := client.Do(req)
	if err != nil {
		return responseType, err
	}

	if res == nil {
		return responseType, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return responseType, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return responseType, fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
	}

	var responseObject R
	err = json.Unmarshal(responseData, &responseObject)

	if err != nil {
		return responseType, fmt.Errorf("error unmarshaling response: %+v", err)
	}

	return responseObject, nil
}

func Get[R any](ctx context.Context, lgr log.Logger, urlfull string, headers map[string]string, queryParameters url.Values, responseType R) (R, error) {

	client := http.Client{}

	u, err := url.Parse(urlfull)
	if err != nil {
		return responseType, err
	}

	// add query parameters

	q := u.Query()
	for k, v := range queryParameters {
		// this depends on the type of api, you may need to do it for each of v
		q.Set(k, strings.Join(v, ","))
	}
	// set the query to the encoded parameters
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return responseType, fmt.Errorf("could not prepare http request: %w", err)
	}
	// add headers
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// finally, do the request
	res, err := client.Do(req)
	if err != nil {
		return responseType, err
	}

	if res == nil {
		return responseType, fmt.Errorf("error: calling %s returned empty response", u.String())
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return responseType, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return responseType, fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", u.String(), res.Status, responseData)
	}

	var responseObject R
	err = json.Unmarshal(responseData, &responseObject)

	if err != nil {
		return responseType, fmt.Errorf("error unmarshaling response: %+v", err)
	}

	return responseObject, nil

}

func toJSON(T any) ([]byte, error) {
	return json.Marshal(T)
}

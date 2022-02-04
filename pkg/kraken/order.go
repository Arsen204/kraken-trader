package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrSignRequest      = errors.New("sign request error")
	ErrCreateNewRequest = errors.New("cannot create new request")
	ErrSendRequest      = errors.New("send request error")
	ErrBadStatusCode    = errors.New("bad status code error")
	ErrResponseBodyRead = errors.New("response body read error")
	ErrPublishResponse  = errors.New("response publish error")
)

type OrderRequest struct {
	Symbol     string
	Side       string
	Size       float64
	LimitPrice float64
}

const clientRequestTimeout = 10 * time.Second

func (k *krakenService) PlaceOrder(order OrderRequest) error {
	req, err := k.createRequest(order)
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: clientRequestTimeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSendRequest, err)
	}
	defer res.Body.Close()

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return ErrBadStatusCode
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrResponseBodyRead, err)
	}

	_, err = k.publisher.Write(b)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPublishResponse, err)
	}

	return nil
}

func (k *krakenService) createRequest(order OrderRequest) (*http.Request, error) {
	endpointPath := k.config.RestEndpoint + "/sendorder"
	u := url.URL{Scheme: "https", Host: k.config.BaseDemoURI, Path: endpointPath}

	// Query parameters
	v := url.Values{}
	v.Add("orderType", "ioc")
	v.Add("symbol", order.Symbol)
	v.Add("side", order.Side)
	v.Add("size", strconv.FormatFloat(order.Size, 'f', -1, 64))
	v.Add("limitPrice", strconv.FormatFloat(order.LimitPrice, 'f', -1, 64))
	queryString := v.Encode()

	// Creating request
	req, err := http.NewRequest(http.MethodPost, u.String()+"?"+queryString, nil)
	if err != nil {
		return nil, ErrCreateNewRequest
	}

	// Nonce
	nonce := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// Authentication
	authent, err := k.signRequest(endpointPath, nonce, queryString)
	if err != nil {
		return nil, ErrSignRequest
	}

	// Headers
	req.Header.Add("APIKey", k.config.APIKey)
	req.Header.Add("Authent", authent)
	req.Header.Add("Nonce", nonce)

	return req, nil
}

func (k *krakenService) signRequest(endpoint string, nonce string, postData string) (string, error) {
	endpoint = strings.TrimPrefix(endpoint, "/derivatives")
	message := []byte(postData + nonce + endpoint)

	sha := sha256.New()
	sha.Write(message)

	secretDecoded, err := base64.StdEncoding.DecodeString(k.config.APISecret)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha512.New, secretDecoded)
	h.Write(sha.Sum(nil))

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

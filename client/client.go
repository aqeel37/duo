package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/aqeel/cq-source-duo/client/admin" // important note
	"github.com/rs/zerolog"
)

type Client struct {
	logger zerolog.Logger
	Spec   Spec
}

func (c *Client) ID() string {
	return "duo"
}

func (c *Client) Logger() *zerolog.Logger {
	return &c.logger
}

const (
	usersURI = "/admin/v1/users"
)

func New(ctx context.Context, logger zerolog.Logger, s *Spec) (Client, error) {
	// TODO: Add your client initialization here
	c := Client{
		logger: logger,
		Spec:   *s,
	}

	return c, nil
}

func generateHmacSha1Signature(key, data string) string {
	hmacSha1 := hmac.New(sha1.New, []byte(key))
	hmacSha1.Write([]byte(data))
	return fmt.Sprintf("%x", hmacSha1.Sum(nil))
}

func sign(method, host, path, skey, ikey string, params map[string]string) map[string]string {

	// create canonical string with the provided date
	// Replace this with your RFC 2822 formatted string
	// rfc2822String := "Tue, 10 Nov 2009 23:00:00 UTC"

	// Parse the RFC 2822 formatted string
	currentTime := time.Now().UTC()

	// Format the current time in the specified layout
	customFormat := "Mon, 02 Jan 2006 15:04:05 -0700"
	date := currentTime.Format(customFormat)

	canon := []string{date, method, strings.ToLower(host), path}
	args := make([]string, 0, len(params))
	for key, val := range params {
		args = append(args, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(val)))
	}
	sort.Strings(args)

	canon = append(canon, strings.Join(args, "&"))
	canonStr := strings.Join(canon, "\n")

	// sign canonical string
	sig := generateHmacSha1Signature(skey, canonStr)
	// Convert the signature to hexadecimal ASCII
	hexSignature := fmt.Sprintf("%x", sig)
	auth := fmt.Sprintf("%s:%s", ikey, hexSignature)

	// return headers
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth))),
		"Date":          date,
	}

	return headers
}

func (c *Client) GetUsers(ctx context.Context, page int, size int) ([]admin.User, error) {
	// calling the actual API directly
	tokenUrl := c.Spec.ApiHost + usersURI

	users := make([]admin.User, 0)

	clientReq := &http.Client{}
	req, err := http.NewRequest("GET", tokenUrl, nil)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"username": c.Spec.ClientId,
		// "offset":   page,
		// "limit":    size,
	}
	headers := sign("GET", c.Spec.ApiHost, usersURI, c.Spec.ClientSecret, c.Spec.ClientId, params)

	authKey, _ := headers["Authorization"]
	date, _ := headers["Date"]
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {authKey},
		"Date":          {date},
	}

	resp, err := clientReq.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userResponse admin.UserResponseDTO
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		c.logger.Error().Msg(fmt.Sprintf("Error decoding JSON: %v", err))
		return nil, err
	}

	users = append(users, userResponse.Response...)
	return users, nil
}

func nowInSeconds() int64 {
	return time.Now().Unix()
}

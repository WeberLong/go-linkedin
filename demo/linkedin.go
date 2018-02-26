package linkedin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"net/http"

	"github.com/WeberLong/go-linkedin/linkedinAPI"
)

const (
	client_id     = "819oy********12"
	client_secret = "nOm********w4fgrO"
	redirect_url  = "https://**.*****.com/*****.html"
	state         = "2b2ff424b0c5c8499540b8c55997c4ac"
)

type responseLinkedinResult struct {
	Id string `json:"id"`
}

func Login(c echo.Context) error {
	req := c.Request()
	res := c.Response()

	client := &linkedinAPI.API{}
	client.SetCredentials(client_id, client_secret)
	client.Auth(res, req, state, redirect_url)
	return nil
}

func Auth(c echo.Context) error {
	httpClient := http.DefaultClient
	client := &linkedinAPI.API{}
	client.SetCredentials(client_id, client_secret)

	if code := c.QueryParam("code"); code != "" {
		if state == c.QueryParam("state") {
			accessToken, err := client.RetrieveAccessToken(httpClient, code, redirect_url)
			if err != nil {
				err := errors.New("领英accessToken失效")
				return err
			}
			client.SetToken(accessToken)

			li := linkedinAPI.Fields{}
			li.Add("id")
			li.Add("email-address")
			li.Add("first-name")
			li.Add("last-name")
			li.Add("maiden-name")
			li.Add("formatted-name")
			li.Add("headline")
			li.Add("specialties")
			li.Add("num-connections")
			li.Add("picture-url")
			li.Add("location")
			li.Add("summary")
			li.Add("positions")
			li.Add("site-standard-profile-request")
			resp, err := client.Profile(httpClient, "~", li)
			if err != nil {
				return err
			}

			byteRespData, err := json.Marshal(resp)
			if err != nil {
				err := errors.New("json编码职业信息出错")
				return err
			}

			var respData responseLinkedinResult
			err = json.Unmarshal(byteRespData, &respData)

			if err != nil {
				err := errors.New("json编码职业信息出错")
				return err
			}

			return nil
		} else {
			err := errors.New("Invalid state!")
			return err
		}
	} else {
		fmt.Printf("%+v\n", c.QueryParam("error"))
		err := errors.New("error")
		return err
	}
}

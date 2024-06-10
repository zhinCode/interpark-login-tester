package interpark

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"interparkTester/internal/httpclient"
	"interparkTester/pkg/custom"

	"golang.org/x/net/html"
)

var (
	userAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	step0LoginURL  = "https://accounts.interpark.com/authorize/inpark-pc?postProc=FULLSCREEN&version=v2&origin=https%3A%2F%2Fwww.interpark.com%2F"
	step1LoginURL  = "https://accounts.interpark.com/login/submit"
	step3MypageURL = "https://incorp.interpark.com/member/memberjoin.do?_method=memberInfoModify"
	step4LogoutURL = "https://accounts.interpark.com/logout?retUrl=https%3A%2F%2Fincorp.interpark.com%2Fmember%2Fmemberjoin.do"
	refStep0       = "https://www.interpark.com/"
	refStep1       = "https://accounts.interpark.com/login/form"
	refStep2       = "https://accounts.interpark.com/"
	refStep3       = "https://accounts.interpark.com/"
	refStep4       = "https://incorp.interpark.com/"
	step1PostProc  = "FULLSCREEN"
	step1LOGIN_TP  = "1500"
	step1fromSVC   = "inpark"
	step1bizId     = "15"
)

func LoginProcess(userID, userPwd string) error {
	client := httpclient.NewHttpClient()

	if err := prepareLogin(client); err != nil {
		return fmt.Errorf("preparing login: %w", err)
	}

	responseStruct, err := login(client, userID, userPwd)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	if err := callback(client, responseStruct.CallbackURL); err != nil {
		return fmt.Errorf("callback: %w", err)
	}

	if err := fetchUserInfo(client); err != nil {
		return fmt.Errorf("fetching user info: %w", err)
	}

	custom.Logout(client, step4LogoutURL, refStep4, userAgent)
	return nil
}

func prepareLogin(client *http.Client) error {
	req, err := http.NewRequest("GET", step0LoginURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Referer", refStep0)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if location := resp.Header.Get("Location"); location != "" {
		step0LoginURL = location
		req, err = http.NewRequest("GET", step0LoginURL, nil)
		if err != nil {
			return fmt.Errorf("creating request: %w", err)
		}
		req.Header.Add("User-Agent", userAgent)
		req.Header.Add("Referer", refStep0)

		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("making request: %w", err)
		}
		defer resp.Body.Close()
	}

	return nil
}

func login(client *http.Client, userID, userPwd string) (*LoginResponse, error) {
	data := url.Values{}
	data.Set("userId", userID)
	data.Set("userPwd", userPwd)
	data.Set("postProc", step1PostProc)
	data.Set("LOGIN_TP", step1LOGIN_TP)
	data.Set("fromSVC", step1fromSVC)
	data.Set("bizId", step1bizId)

	req, err := http.NewRequest("POST", step1LoginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Referer", refStep1)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var responseStruct LoginResponse
	if err := json.Unmarshal(body, &responseStruct); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	if responseStruct.CallbackURL == "" {
		return nil, fmt.Errorf("callback URL not found in response")
	}

	return &responseStruct, nil
}

func callback(client *http.Client, callbackURL string) error {
	req, err := http.NewRequest("GET", callbackURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Referer", refStep2)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func fetchUserInfo(client *http.Client) error {
	req, err := http.NewRequest("GET", step3MypageURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("Referer", refStep3)
	req.Header.Add("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	charset := "utf-8"
	if strings.Contains(contentType, "charset=") {
		idx := strings.Index(contentType, "charset=")
		charset = strings.ToUpper(contentType[idx+len("charset="):])
	}

	utf8Body, err := custom.DecodeBody(body, charset)
	if err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}

	doc, err := html.Parse(strings.NewReader(utf8Body))
	if err != nil {
		return fmt.Errorf("parsing HTML: %w", err)
	}

	arrMyAreaElements := custom.GetElementByClass(doc, "div", "myArea")
	for _, myArea := range arrMyAreaElements {
		arrPtags := custom.GetTextByTag(myArea, "p")
		name := ""
		re := regexp.MustCompile(`(.*)님은`)
		match := re.FindStringSubmatch(arrPtags[0])
		if len(match) > 1 {
			name = match[1]
			fmt.Println(string(custom.Yellow), "Hello! ["+name+"] Sir!", string(custom.Reset))
		} else {
			return fmt.Errorf("login info page failed")
		}
	}

	return nil
}

type LoginResponse struct {
	CallbackURL string `json:"callback_url"`
	RedirectURI string `json:"redirect_uri"`
	ResultCode  string `json:"result_code"`
}

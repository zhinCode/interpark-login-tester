package custom

import (
	"fmt"
	"net/http"
	"time"
)

func Logout(client *http.Client, step4MypageURL string, refStep4 string, userAgent string) {
	time.Sleep(1*time.Second + 500*time.Millisecond)
	req, err := http.NewRequest("GET", step4MypageURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("Referer", refStep4)
	req.Header.Add("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	cookies := client.Jar.Cookies(req.URL)
	for _, cookie := range cookies {
		cookie.MaxAge = -1
		client.Jar.SetCookies(req.URL, []*http.Cookie{cookie})
	}
	fmt.Println("Logout successful")
}

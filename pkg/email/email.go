package email

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nazarnovak/hobee-be/pkg/herrors2"
)

var (
	apiKey = ""
	domain = ""
)

func Init(cnfApiKey, cnfDomain string) error {
	if cnfApiKey == "" {
		return herrors.New("apiKey should not be empty", "service", "email")
	}

	if cnfDomain == "" {
		return herrors.New("domain should not be empty", "service", "email")
	}

	apiKey = cnfApiKey
	domain = cnfDomain

	return nil
}

func Send(subject, text string) error {
	form := url.Values{}

	form.Add("from", "Excited User <mailgun@sandboxc82077688c034509a3f1dd726b9ae90f.mailgun.org>")
	form.Add("to", "mistikal91@gmail.com")
	form.Add("subject", subject)
	form.Add("text", text)

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", domain),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return herrors.Wrap(err)
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	request.SetBasicAuth("api",  apiKey)

	httpClient := http.Client{Timeout: 30 * time.Second}

	resp, err := httpClient.Do(request)
	if err != nil {
		return herrors.Wrap(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return herrors.Wrap(err)
	}

	if resp.StatusCode != http.StatusOK {
		return herrors.New("Invalid response from send", "statuscode", resp.StatusCode,
			"body", string(body))
	}

	return nil
}

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

type (
	// GitHub information.
	GitHub struct {
		Workflow  string
		Workspace string
		Action    string
		EventName string
		EventPath string
	}

	// Repo information
	Repo struct {
		FullName  string
		Namespace string
		Name      string
	}

	// Build information
	Build struct {
		Tag      string
		Event    string
		Number   int64
		Commit   string
		RefSpec  string
		Branch   string
		Author   string
		Avatar   string
		Message  string
		Email    string
		Status   string
		Link     string
		Started  float64
		Finished float64
	}

	// Config for the plugin.
	Config struct {
		URL       string
		Insecure  bool
		TrustedCA string
		UserID    string
		Token     string
		Message   string
		Drone     bool
		GitHub    bool
		EnvFile   string
		Debug     bool
	}

	// Payload struct
	Payload struct {
		Channel         string                 `json:"channel"`
		Text            string                 `json:"text"`
		Avatar          string                 `json:"avatar"`
		CustomMSgFields map[string]interface{} `json:"custommsg"`
	}

	// Plugin values.
	Plugin struct {
		GitHub  GitHub
		Repo    Repo
		Build   Build
		Config  Config
		Payload Payload
	}
)

func clientHTTP(p *Plugin) (client *http.Client) {
	var urlHTTPSScheme = regexp.MustCompile(`^https://.+`)
	// https url specified
	if urlHTTPSScheme.MatchString(p.Config.URL) {

		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatal(err)
		}
		if caCertPool == nil {
			caCertPool = x509.NewCertPool()
		}
		// custom CAs given
		if len(p.Config.TrustedCA) > 0 {
			caCert, err := ioutil.ReadFile(p.Config.TrustedCA)
			if err != nil {
				log.Fatal(err)
			}
			caCertPool.AppendCertsFromPEM(caCert)
		}
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:            caCertPool,
					InsecureSkipVerify: p.Config.Insecure,
				},
			},
		}
		return client
	}
	// plain http
	return &http.Client{}
}

// Exec executes the plugin.
func (p *Plugin) Exec() error {
	var message string
	var txt bytes.Buffer

	if p.Config.URL == "" || p.Config.UserID == "" || p.Config.Token == "" {
		return errors.New("missing Rocket.Chat config")
	}

	if len(p.Config.Message) > 0 {
		message = p.Config.Message
	} else {
		message = p.Message()
	}

	t := template.Must(template.New("message").Parse(message))
	err := t.Execute(&txt, p)
	if err != nil {
		fmt.Println("something went wrong executing template:", err)
		txt.WriteString(p.Message())
	}
	err = p.SendMessage(txt.String())
	if err != nil {
		return err
	}

	return nil
}

// Message returns default formatted message
func (p Plugin) Message() string {
	if p.Config.GitHub {
		return fmt.Sprintf("%s/%s triggered by %s (%s)",
			p.Repo.FullName,
			p.GitHub.Workflow,
			p.Repo.Namespace,
			p.GitHub.EventName,
		)
	}

	return fmt.Sprintf("[%s] <%s> (%s)『%s』by %s",
		p.Build.Status,
		p.Build.Link,
		p.Build.Branch,
		p.Build.Message,
		p.Build.Author,
	)
}

// SendMessage actually sends message to Rocket.Chat
func (p *Plugin) SendMessage(msg string) error {
	URL := fmt.Sprintf("%s/api/v1/chat.postMessage", p.Config.URL)
	b := new(bytes.Buffer)
	var payload = p.Payload
	payload.Text = msg
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return err
	}
	client := clientHTTP(p)
	req, err := http.NewRequest("POST", URL, b)
	if err != nil {
		return err
	}
	// must specify content-type
	req.Header.Add("Content-type", "application/json")
	// add auth headers
	req.Header.Add("X-User-Id", p.Config.UserID)
	req.Header.Add("X-Auth-Token", p.Config.Token)
	// perform request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed to send notification:\nRESPONSE\nStatus code: %d\nHeaders: %s\nBody: %s\n", resp.StatusCode, resp.Header, string(body))
	} else {
		log.Output(2, "Message sent")
	}

	return nil
}

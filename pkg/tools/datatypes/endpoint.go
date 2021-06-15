package datatypes

import (
	"log"
	"net/url"
	"strconv"
	"strings"
)

func ParseEndpoint(text string) (*Endpoint, error) {
	u, err := url.ParseRequestURI(text)
	if err != nil {
		return nil, err
	}

	endpoint := &Endpoint{
		Scheme: u.Scheme,
	}

	query := u.Query()

	if len(query) > 0 {
		endpoint.Extra = u.Query()
	}

	if len(u.Path) > 0 {
		endpoint.Base = strings.Split(u.Path[1:], "/")[0]
	}

	endpoint.Hostname = u.Hostname()

	i, err := strconv.ParseUint(u.Port(), 10, 16)
	if err == nil {
		endpoint.Port = uint16(i)
	}

	if u.User != nil {
		endpoint.Username = u.User.Username()
		endpoint.Password, _ = u.User.Password()
	}

	return endpoint, nil
}

// openapi:strfmt endpoint
type Endpoint struct {
	Scheme   string
	Hostname string
	Port     uint16
	Base     string
	Username string
	Password string
	Extra    url.Values
}

func (e Endpoint) IsZero() bool {
	return e.Hostname == ""
}

func (e Endpoint) SecurityString() string {
	e.Password = "------"
	return e.String()
}

func (e Endpoint) Host() string {
	if e.Port != 0 {
		return e.Hostname + ":" + strconv.FormatUint(uint64(e.Port), 10)
	}
	return e.Hostname
}

func (e Endpoint) String() string {
	u := url.URL{}
	u.Scheme = e.Scheme
	u.Host = e.Host()

	if e.Extra != nil {
		u.RawQuery = e.Extra.Encode()
	}

	if e.Base != "" {
		u.Path = "/" + e.Base
	}

	if e.Username != "" || e.Password != "" {
		u.User = url.UserPassword(e.Username, e.Password)
	}

	s, err := url.QueryUnescape(u.String())
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func (e *Endpoint) IsTLS() bool {
	if e.Scheme == "" {
		return false
	}
	return e.Scheme[len(e.Scheme)-1] == 's'
}

func (e *Endpoint) UnmarshalText(text []byte) error {
	endpoint, err := ParseEndpoint(string(text))
	if err != nil {
		return err
	}
	*e = *endpoint
	return nil
}

func (e Endpoint) MarshalText() (text []byte, err error) {
	return []byte(e.String()), nil
}

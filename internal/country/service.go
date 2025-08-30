package country

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type Option struct {
	Name         string   `json:"name"`
	ISO2         string   `json:"iso2"`
	ISO3         string   `json:"iso3"`
	CallingCodes []string `json:"callingCodes"`
	FlagEmoji    string   `json:"flagEmoji"`
	FlagSVG      string   `json:"flagSvg"`
	FlagPNG      string   `json:"flagPng"`
	Label        string   `json:"label"`
}

type rcCountry struct {
	CCA2 string `json:"cca2"`
	CCA3 string `json:"cca3"`
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	IDD struct {
		Root     string   `json:"root"`
		Suffixes []string `json:"suffixes"`
	} `json:"idd"`
}

type Service struct {
	http *http.Client
	ttl  time.Duration

	mu        sync.RWMutex
	data      []Option
	expiresAt time.Time
}

func NewService() *Service {
	return &Service{
		http: &http.Client{Timeout: 10 * time.Second},
		ttl:  24 * time.Hour,
	}
}

func (s *Service) iso2ToFlag(iso2 string) string {
	if len(iso2) != 2 {
		return ""
	}
	r := []rune(strings.ToUpper(iso2))
	base := rune(0x1F1E6)
	return string(base+(r[0]-'A')) + string(base+(r[1]-'A'))
}

func (s *Service) buildLabel(c Option) string {
	code := ""
	if len(c.CallingCodes) > 0 {
		code = " (" + strings.Join(c.CallingCodes, ", ") + ")"
	}
	return s.iso2ToFlag(c.ISO2) + " " + c.Name + code
}

func (s *Service) Load() ([]Option, error) {
	s.mu.RLock()
	if s.data != nil && time.Now().Before(s.expiresAt) {
		data := s.data
		s.mu.RUnlock()
		return data, nil
	}
	s.mu.RUnlock()

	url := "https://restcountries.com/v3.1/all?fields=cca2,cca3,name,idd"
	resp, err := s.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	var raw []rcCountry
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	out := make([]Option, 0, len(raw))
	for _, r := range raw {
		if r.CCA2 == "" || r.CCA3 == "" {
			continue
		}
		iso2 := strings.ToUpper(r.CCA2)
		iso3 := strings.ToUpper(r.CCA3)
		name := r.Name.Common

		var codes []string
		if r.IDD.Root != "" && len(r.IDD.Suffixes) > 0 {
			for _, sfx := range r.IDD.Suffixes {
				codes = append(codes, r.IDD.Root+sfx)
			}
		} else if r.IDD.Root != "" {
			codes = []string{r.IDD.Root}
		}

		lc := strings.ToLower(iso2)
		item := Option{
			Name:         name,
			ISO2:         iso2,
			ISO3:         iso3,
			CallingCodes: codes,
			FlagEmoji:    s.iso2ToFlag(iso2),
			FlagSVG:      "https://flagcdn.com/" + lc + ".svg",
			FlagPNG:      "https://flagcdn.com/w320/" + lc + ".png",
		}
		item.Label = s.buildLabel(item)
		out = append(out, item)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })

	s.mu.Lock()
	s.data = out
	s.expiresAt = time.Now().Add(s.ttl)
	s.mu.Unlock()

	return out, nil
}

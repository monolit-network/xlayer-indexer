package models

import (
	"regexp"
	"strings"
	"time"
)

type SocialSourceTag string

const (
	SocialSourceTwitter  SocialSourceTag = "twitter"
	SocialSourceTelegram SocialSourceTag = "telegram"
)

type SocialSource struct {
	Source SocialSourceTag `json:"source"`
	Value  string          `json:"value"`
}

type SocialPost struct {
	Link            string          `json:"link"`
	Text            string          `json:"text"`
	Date            time.Time       `json:"date"`
	Source          SocialSourceTag `json:"source"`
	Author          string          `json:"author"`
	Relations       []string        `json:"related_to"` // links to other posts
	LinksMentioned  []string        `json:"links_mentioned"`
	TokensMentioned []string        `json:"tokens_mentioned"`
}

func (p *SocialPost) FillMentions() {
	// Extract http(s) URLs
	urlRe := regexp.MustCompile(`https?://[^\s]+`)
	rawLinks := urlRe.FindAllString(p.Text, -1)

	linkSeen := make(map[string]struct{}, len(rawLinks))
	links := make([]string, 0, len(rawLinks))
	for _, l := range rawLinks {
		if _, ok := linkSeen[l]; ok {
			continue
		}
		linkSeen[l] = struct{}{}
		links = append(links, l)
	}

	// Extract cashtags like $ETH, $SOL (case-insensitive), 2-15 chars
	// Capture token without the leading $, must contain at least one letter
	tokenRe := regexp.MustCompile(`(?i)(?:^|[^A-Za-z0-9_])\$([a-z0-9]*[a-z][a-z0-9]*)\b`)
	matches := tokenRe.FindAllStringSubmatch(p.Text, -1)

	tokenSeen := make(map[string]struct{}, len(matches))
	tokens := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		t := strings.ToUpper(m[1])
		if len(t) < 2 || len(t) > 15 {
			continue
		}
		if _, ok := tokenSeen[t]; ok {
			continue
		}
		tokenSeen[t] = struct{}{}
		tokens = append(tokens, t)
	}

	p.LinksMentioned = links
	p.TokensMentioned = tokens
}

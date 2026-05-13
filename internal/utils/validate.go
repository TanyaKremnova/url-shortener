package utils

import (
    "net/url"
    "strings"
)

func IsValidURL(raw string) bool {
    parsed, err := url.ParseRequestURI(raw)
    if err != nil {
        return false
    }

    // Must have http or https scheme
    if parsed.Scheme != "http" && parsed.Scheme != "https" {
        return false
    }

    // Must have a host
    if strings.TrimSpace(parsed.Host) == "" {
        return false
    }

    return true
}
package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var relativeURLPattern = regexp.MustCompile(`(href|src|action)="([^"]+)"`)

func FixRelativeUrls(html, baseURL string) (string, error) {
    base, err := url.Parse(baseURL)
    if err != nil {
        return "", err
    }

    var innerErr error
    fixed := relativeURLPattern.ReplaceAllStringFunc(html, func(match string) string {
        groups := relativeURLPattern.FindStringSubmatch(match)
        attr, path := groups[1], groups[2]

        if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "//") || strings.HasPrefix(path, "data:") {
            return match
        }

        ref, err := url.Parse(path)
        if err != nil {
            innerErr = err
            return match
        }

        absolute := base.ResolveReference(ref).String()
        return fmt.Sprintf(`%s="%s"`, attr, absolute)
    })

    if innerErr != nil {
        return "", innerErr
    }
    return fixed, nil
}

func HideTopBar(html string) string {
    return strings.Replace(html,
        `class="lineedyn_mappa_intestazione" style="display: table-row; height: min-content;"`,
        `class="lineedyn_mappa_intestazione" style="display: none; height: min-content;"`,
        1)
}

func FixLeafletScript(html string) string {
    oldScript := `<script src="https://www.setaweb.it/js/leafletCustomConfig.js"></script>`
    newScript := `<script>
        var Stadia_OSMBright = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            minZoom: 11,
            attribution: '© OpenStreetMap',
            referrerPolicy : 'strict-origin-when-cross-origin'
        });
    </script>`
    return strings.Replace(html, oldScript, newScript, 1)
}
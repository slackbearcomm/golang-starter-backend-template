package helpers

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcloughlin/geohash"
)

// StandardizeSpaces will cleanup redundant whitespace
func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// LimitName will reveal limited name for privacy, show first and last char, and if over 8 char will reveal 1 more randomly
func LimitName(in string) string {
	var rp []rune
	if len(in) > 2 {
		rp = []rune(strings.Repeat("*", len(in)-2))
	}

	// pick random spot to reveal in the middle
	if len(in) >= 8 {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		n := r1.Intn(len(in) - 2)
		rp[n] = []rune(in)[n+1]
	}

	switch len(in) {
	case 0:
		return ""
	case 1:
		return string(in[0:1])
	case 2:
		return fmt.Sprintf("%s*", string(in[0:1]))
	}
	return fmt.Sprintf("%s%s%s", string(in[0:1]), string(rp), string(in[len(in)-1:]))
}

// LimitCoordinate will reveal limit location for privacy
func LimitCoordinate(str string) string {
	lat, lng := geohash.Decode(str)
	return fmt.Sprintf("%d.XXXXXXXX, %d.XXXXXXXX", int(lat), int(lng))
}

// StripHTML removes all html tags in the string to make it safer
func StripHTML(str string) string {
	sanitize := bluemonday.StrictPolicy()
	str2 := sanitize.Sanitize(str)
	return str2
}

// Used to strip strings of any character that is not a letter or a "-". Also allows for spaces. Use before running through StripHTML
// Handles cases for names with characters like apostrophe which will be escaped as a html element by bluemonday
func StripName(str string) string {
	notAlpha := regexp.MustCompile(`[^a-zA-Z\- ]`)
	return notAlpha.ReplaceAllString(str, "")
}

// IsValidEmailAddress checks if email address containers any illegal characters or not
func IsValidEmailAddress(str string) bool {
	// non-ascii
	reNotASCII := regexp.MustCompile(`[^a-z0-9\@\.\-\+\_]`)

	return !reNotASCII.MatchString(str)
}

// StripEmailAddress removes regular typo or illegal characters from email address
func StripEmailAddress(str string) string {
	// non-ascii
	reNotASCII := regexp.MustCompile(`[^a-z0-9\@\.\-\+\_]`)

	// to lowercase
	str = strings.ToLower(str)
	// no spaces
	str = strings.Replace(str, " ", "", -1)
	// no non-ascii
	str = reNotASCII.ReplaceAllString(str, "_")

	return str
}

// FindStringSubmatchMap construct hash of regular expression named capture group
// http://blog.kamilkisiel.net/blog/2012/07/05/using-the-go-regexp-package/
func FindStringSubmatchMap(r *regexp.Regexp, s string) map[string]*string {
	captures := make(map[string]*string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = &match[i]
	}
	return captures
}

package utils

import (
	"errors"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var regxPhone *regexp.Regexp

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ExtractBearerToken extracts token from the 'Authorization: Bearer <token>' header
func ExtractBearerToken(r *http.Request) (string, error) {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return "", errors.New("no Authorization header provided")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("unsupported Authorization header format")
	}

	return parts[1], nil
}

// GenerateUUID creates new id string
func GenerateUUID() string {
	return uuid.New().String()
}

// ValidateUUID validates id
func ValidateUUID(id string) bool {
	if id == "" {
		return false
	}
	_, err := uuid.Parse(id)
	return err == nil
}

// FilterEmail value
func FilterEmail(email string) string {
	email = strings.TrimSpace(email)
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ""
	}
	return addr.Address
}

// ValidateEmail validates email value
func ValidateEmail(email string) bool {
	email = FilterEmail(email)
	return email != ""
}

// ValidatePhone validates phone number
func ValidatePhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false
	}
	if regxPhone == nil {
		regxPhone = regexp.MustCompile(`^[+0-9\s\-()]{5,}$`)
	}
	return regxPhone.MatchString(phone)
}

// NormalizePhone removes all characters from the phone except plus and numbers
func NormalizePhone(phone string) string {
	rx := regexp.MustCompile(`[^+0-9]`)
	b := rx.ReplaceAll([]byte(phone), []byte(""))
	phone = string(b)

	// 89004567890 -> +79004567890
	if len(phone) == 11 && phone[0:1] == "8" {
		phone = "+7" + phone[1:]
	}

	// 9004567890 -> +79004567890
	if len(phone) == 10 && phone[0:1] == "9" {
		phone = "+7" + phone
	}

	return phone
}

// JSONConvert converts anything to anything through JSON marshalling
func JSONConvert(source interface{}, target interface{}) error {
	json, err := jsoniter.Marshal(source)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(json, target)
}

// ValidateURL checks if string is a valid URL value
func ValidateURL(urlAddr string) bool {
	_, err := url.ParseRequestURI(urlAddr)
	return err == nil
}

// FilterStrings returns string slice with only non-empty values
func FilterStrings(items []string) []string {
	filtered := make([]string, 0)
	if items == nil {
		return filtered
	}

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// GenerateRandomString generates random string
func GenerateRandomString(symbols string, length uint) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

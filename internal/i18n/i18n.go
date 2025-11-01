// Package i18n
// Description: Internationalization support for sokru
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package i18n

import (
	"fmt"
	"sync"
)

// Language represents a supported language
type Language string

const (
	English Language = "en"
	Spanish Language = "es"
)

// MessageKey represents a message identifier
type MessageKey string

// Message format prefixes
const (
	PrefixSuccess = "✓"
	PrefixError   = "✗"
	PrefixInfo    = "→"
	PrefixWarning = "⚠"
)

// Global i18n instance
var (
	instance *I18n
	once     sync.Once
)

// I18n manages internationalization
type I18n struct {
	currentLang Language
	messages    map[Language]map[MessageKey]string
	mu          sync.RWMutex
}

// GetInstance returns the singleton i18n instance
func GetInstance() *I18n {
	once.Do(func() {
		instance = &I18n{
			currentLang: English,
			messages:    make(map[Language]map[MessageKey]string),
		}
		instance.loadMessages()
	})
	return instance
}

// SetLanguage sets the current language
func (i *I18n) SetLanguage(lang Language) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.currentLang = lang
}

// GetLanguage returns the current language
func (i *I18n) GetLanguage() Language {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.currentLang
}

// T translates a message key to the current language
func (i *I18n) T(key MessageKey, args ...interface{}) string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	langMessages, ok := i.messages[i.currentLang]
	if !ok {
		// Fallback to English
		langMessages = i.messages[English]
	}

	message, ok := langMessages[key]
	if !ok {
		return string(key) // Return key if translation not found
	}

	if len(args) > 0 {
		return fmt.Sprintf(message, args...)
	}
	return message
}

// Success formats a success message
func (i *I18n) Success(key MessageKey, args ...interface{}) string {
	return fmt.Sprintf("%s %s", PrefixSuccess, i.T(key, args...))
}

// Error formats an error message
func (i *I18n) Error(key MessageKey, args ...interface{}) string {
	return fmt.Sprintf("%s %s", PrefixError, i.T(key, args...))
}

// Info formats an info message
func (i *I18n) Info(key MessageKey, args ...interface{}) string {
	return fmt.Sprintf("%s %s", PrefixInfo, i.T(key, args...))
}

// Warning formats a warning message
func (i *I18n) Warning(key MessageKey, args ...interface{}) string {
	return fmt.Sprintf("%s %s", PrefixWarning, i.T(key, args...))
}

// loadMessages loads all translations
func (i *I18n) loadMessages() {
	i.messages[English] = getEnglishMessages()
	i.messages[Spanish] = getSpanishMessages()
}

// Helper functions for global access
func T(key MessageKey, args ...interface{}) string {
	return GetInstance().T(key, args...)
}

func Success(key MessageKey, args ...interface{}) string {
	return GetInstance().Success(key, args...)
}

func Error(key MessageKey, args ...interface{}) string {
	return GetInstance().Error(key, args...)
}

func Info(key MessageKey, args ...interface{}) string {
	return GetInstance().Info(key, args...)
}

func Warning(key MessageKey, args ...interface{}) string {
	return GetInstance().Warning(key, args...)
}

func SetLanguage(lang Language) {
	GetInstance().SetLanguage(lang)
}

func GetLanguage() Language {
	return GetInstance().GetLanguage()
}

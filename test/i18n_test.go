// Package test
// Description: Unit tests for internationalization
// (c) 2024 Alejandro Lopez Monzon <alejandro@kreaker.dev>

package test

import (
	"strings"
	"testing"

	"github.com/alexlm78/sokru/internal/i18n"
)

func TestGetInstance(t *testing.T) {
	i1 := i18n.GetInstance()
	i2 := i18n.GetInstance()

	if i1 != i2 {
		t.Error("GetInstance should return the same instance (singleton)")
	}

	if i1 == nil {
		t.Fatal("GetInstance returned nil")
	}
}

func TestSetAndGetLanguage(t *testing.T) {
	i18nInstance := i18n.GetInstance()

	// Test setting English
	i18nInstance.SetLanguage(i18n.English)
	if i18nInstance.GetLanguage() != i18n.English {
		t.Errorf("Expected language %s, got %s", i18n.English, i18nInstance.GetLanguage())
	}

	// Test setting Spanish
	i18nInstance.SetLanguage(i18n.Spanish)
	if i18nInstance.GetLanguage() != i18n.Spanish {
		t.Errorf("Expected language %s, got %s", i18n.Spanish, i18nInstance.GetLanguage())
	}

	// Reset to English for other tests
	i18nInstance.SetLanguage(i18n.English)
}

func TestTranslation(t *testing.T) {
	i18nInstance := i18n.GetInstance()

	tests := []struct {
		name     string
		lang     i18n.Language
		key      i18n.MessageKey
		args     []interface{}
		contains string
	}{
		{
			name:     "English - simple message",
			lang:     i18n.English,
			key:      i18n.MsgSymlinksStatus,
			args:     nil,
			contains: "Symlinks Status",
		},
		{
			name:     "Spanish - simple message",
			lang:     i18n.Spanish,
			key:      i18n.MsgSymlinksStatus,
			args:     nil,
			contains: "Enlaces Simbólicos",
		},
		{
			name:     "English - message with parameter",
			lang:     i18n.English,
			key:      i18n.MsgFoundConfigurations,
			args:     []interface{}{5},
			contains: "5",
		},
		{
			name:     "Spanish - message with parameter",
			lang:     i18n.Spanish,
			key:      i18n.MsgFoundConfigurations,
			args:     []interface{}{3},
			contains: "3",
		},
		{
			name:     "English - error message",
			lang:     i18n.English,
			key:      i18n.MsgErrorLoadingConfig,
			args:     []interface{}{"test error"},
			contains: "Error loading configuration",
		},
		{
			name:     "Spanish - error message",
			lang:     i18n.Spanish,
			key:      i18n.MsgErrorLoadingConfig,
			args:     []interface{}{"error de prueba"},
			contains: "Error al cargar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i18nInstance.SetLanguage(tt.lang)
			result := i18nInstance.T(tt.key, tt.args...)

			if !strings.Contains(result, tt.contains) {
				t.Errorf("Translation for %s in %s should contain '%s', got '%s'",
					tt.key, tt.lang, tt.contains, result)
			}
		})
	}

	// Reset to English
	i18nInstance.SetLanguage(i18n.English)
}

func TestFormattedMessages(t *testing.T) {
	i18nInstance := i18n.GetInstance()
	i18nInstance.SetLanguage(i18n.English)

	tests := []struct {
		name   string
		fn     func(i18n.MessageKey, ...interface{}) string
		key    i18n.MessageKey
		prefix string
	}{
		{
			name:   "Success message",
			fn:     i18nInstance.Success,
			key:    i18n.MsgSymlinkCreated,
			prefix: i18n.PrefixSuccess,
		},
		{
			name:   "Error message",
			fn:     i18nInstance.Error,
			key:    i18n.MsgErrorLoadingConfig,
			prefix: i18n.PrefixError,
		},
		{
			name:   "Info message",
			fn:     i18nInstance.Info,
			key:    i18n.MsgReadingSymlinksFrom,
			prefix: i18n.PrefixInfo,
		},
		{
			name:   "Warning message",
			fn:     i18nInstance.Warning,
			key:    i18n.MsgNotSymlink,
			prefix: i18n.PrefixWarning,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.key, "test")

			if !strings.HasPrefix(result, tt.prefix) {
				t.Errorf("Message should start with prefix '%s', got '%s'", tt.prefix, result)
			}
		})
	}
}

func TestGlobalHelperFunctions(t *testing.T) {
	// Test global T function
	i18n.SetLanguage(i18n.English)
	result := i18n.T(i18n.MsgSymlinksStatus)
	if !strings.Contains(result, "Symlinks") {
		t.Errorf("Global T() should work, got '%s'", result)
	}

	// Test global Success function
	result = i18n.Success(i18n.MsgSymlinkCreated, "test", "test")
	if !strings.HasPrefix(result, i18n.PrefixSuccess) {
		t.Errorf("Global Success() should have prefix '%s', got '%s'", i18n.PrefixSuccess, result)
	}

	// Test global Error function
	result = i18n.Error(i18n.MsgErrorLoadingConfig, "test")
	if !strings.HasPrefix(result, i18n.PrefixError) {
		t.Errorf("Global Error() should have prefix '%s', got '%s'", i18n.PrefixError, result)
	}

	// Test global Info function
	result = i18n.Info(i18n.MsgReadingSymlinksFrom, "test")
	if !strings.HasPrefix(result, i18n.PrefixInfo) {
		t.Errorf("Global Info() should have prefix '%s', got '%s'", i18n.PrefixInfo, result)
	}

	// Test global Warning function
	result = i18n.Warning(i18n.MsgNotSymlink, "test")
	if !strings.HasPrefix(result, i18n.PrefixWarning) {
		t.Errorf("Global Warning() should have prefix '%s', got '%s'", i18n.PrefixWarning, result)
	}

	// Test global GetLanguage function
	lang := i18n.GetLanguage()
	if lang != i18n.English {
		t.Errorf("Global GetLanguage() should return %s, got %s", i18n.English, lang)
	}
}

func TestMissingTranslation(t *testing.T) {
	i18nInstance := i18n.GetInstance()
	i18nInstance.SetLanguage(i18n.English)

	// Test with a non-existent key
	nonExistentKey := i18n.MessageKey("non_existent_key")
	result := i18nInstance.T(nonExistentKey)

	// Should return the key itself when translation is not found
	if result != string(nonExistentKey) {
		t.Errorf("Missing translation should return key, got '%s'", result)
	}
}

func TestLanguageFallback(t *testing.T) {
	i18nInstance := i18n.GetInstance()

	// Set to an invalid language (should fallback to English)
	i18nInstance.SetLanguage(i18n.Language("invalid"))

	// Try to get a translation - should fallback to English
	result := i18nInstance.T(i18n.MsgSymlinksStatus)
	if !strings.Contains(result, "Symlinks") {
		t.Errorf("Should fallback to English, got '%s'", result)
	}

	// Reset to English
	i18nInstance.SetLanguage(i18n.English)
}

func TestMessagePrefixes(t *testing.T) {
	prefixes := []struct {
		name   string
		prefix string
	}{
		{"Success", i18n.PrefixSuccess},
		{"Error", i18n.PrefixError},
		{"Info", i18n.PrefixInfo},
		{"Warning", i18n.PrefixWarning},
	}

	for _, p := range prefixes {
		t.Run(p.name, func(t *testing.T) {
			if p.prefix == "" {
				t.Errorf("Prefix %s should not be empty", p.name)
			}
		})
	}
}

func TestAllEnglishMessagesExist(t *testing.T) {
	// Test a sample of important message keys
	importantKeys := []i18n.MessageKey{
		i18n.MsgErrorLoadingConfig,
		i18n.MsgSymlinkCreated,
		i18n.MsgSymlinksStatus,
		i18n.MsgApplyingChanges,
		i18n.MsgInitializing,
		i18n.MsgConfigFile,
		i18n.MsgVersion,
	}

	i18nInstance := i18n.GetInstance()
	i18nInstance.SetLanguage(i18n.English)

	for _, key := range importantKeys {
		result := i18nInstance.T(key)
		if result == string(key) {
			t.Errorf("English message for key '%s' does not exist", key)
		}
	}
}

func TestAllSpanishMessagesExist(t *testing.T) {
	// Test a sample of important message keys
	importantKeys := []i18n.MessageKey{
		i18n.MsgErrorLoadingConfig,
		i18n.MsgSymlinkCreated,
		i18n.MsgSymlinksStatus,
		i18n.MsgApplyingChanges,
		i18n.MsgInitializing,
		i18n.MsgConfigFile,
		i18n.MsgVersion,
	}

	i18nInstance := i18n.GetInstance()
	i18nInstance.SetLanguage(i18n.Spanish)

	for _, key := range importantKeys {
		result := i18nInstance.T(key)
		if result == string(key) {
			t.Errorf("Spanish message for key '%s' does not exist", key)
		}
	}

	// Reset to English
	i18nInstance.SetLanguage(i18n.English)
}

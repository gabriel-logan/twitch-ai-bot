package config

import (
	"os"
	"testing"
	"time"
)

func TestMustExistBool(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		want      bool
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid true",
			key:     "TEST_BOOL",
			value:   "true",
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid false",
			key:     "TEST_BOOL",
			value:   "false",
			want:    false,
			wantErr: false,
		},
		{
			name:      "invalid bool",
			key:       "TEST_BOOL",
			value:     "notabool",
			want:      false,
			wantErr:   true,
			errSubstr: "must be a valid boolean",
		},
		{
			name:      "empty string",
			key:       "TEST_BOOL_EMPTY",
			value:     "",
			want:      false,
			wantErr:   true,
			errSubstr: "is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got, err := mustExistBool(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("mustExistBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" {
				if err != nil && !contains(err.Error(), tt.errSubstr) {
					t.Errorf("mustExistBool() error = %v, should contain %v", err, tt.errSubstr)
				}
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("mustExistBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustExistGoEnv(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		want      string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid debug",
			key:     "TEST_GIN_MODE",
			value:   "debug",
			want:    "debug",
			wantErr: false,
		},
		{
			name:    "valid release",
			key:     "TEST_GIN_MODE",
			value:   "release",
			want:    "release",
			wantErr: false,
		},
		{
			name:      "invalid value",
			key:       "TEST_GIN_MODE",
			value:     "invalid",
			want:      "",
			wantErr:   true,
			errSubstr: "must be debug or release",
		},
		{
			name:      "empty string",
			key:       "TEST_GIN_MODE_EMPTY",
			value:     "",
			want:      "",
			wantErr:   true,
			errSubstr: "is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got, err := mustExistGoEnv(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("mustExistGoEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" {
				if err != nil && !contains(err.Error(), tt.errSubstr) {
					t.Errorf("mustExistGoEnv() error = %v, should contain %v", err, tt.errSubstr)
				}
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("mustExistGoEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustExistString(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		want      string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid string",
			key:     "TEST_STRING",
			value:   "hello world",
			want:    "hello world",
			wantErr: false,
		},
		{
			name:      "empty string",
			key:       "TEST_STRING_EMPTY",
			value:     "",
			want:      "",
			wantErr:   true,
			errSubstr: "is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got, err := mustExistString(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("mustExistString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" {
				if err != nil && !contains(err.Error(), tt.errSubstr) {
					t.Errorf("mustExistString() error = %v, should contain %v", err, tt.errSubstr)
				}
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("mustExistString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustExistInt(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		want      int
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid positive int",
			key:     "TEST_INT",
			value:   "100",
			want:    100,
			wantErr: false,
		},
		{
			name:    "valid zero int",
			key:     "TEST_INT",
			value:   "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "valid negative int",
			key:     "TEST_INT",
			value:   "-5",
			want:    -5,
			wantErr: false,
		},
		{
			name:      "invalid int",
			key:       "TEST_INT",
			value:     "notanint",
			want:      0,
			wantErr:   true,
			errSubstr: "must be a valid integer",
		},
		{
			name:      "empty string",
			key:       "TEST_INT_EMPTY",
			value:     "",
			want:      0,
			wantErr:   true,
			errSubstr: "is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got, err := mustExistInt(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("mustExistInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" {
				if err != nil && !contains(err.Error(), tt.errSubstr) {
					t.Errorf("mustExistInt() error = %v, should contain %v", err, tt.errSubstr)
				}
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("mustExistInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustExistStringSlice(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		want      []string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid single element",
			key:     "TEST_SLICE",
			value:   "one",
			want:    []string{"one"},
			wantErr: false,
		},
		{
			name:    "valid multiple elements",
			key:     "TEST_SLICE",
			value:   "one,two,three",
			want:    []string{"one", "two", "three"},
			wantErr: false,
		},
		{
			name:    "valid empty element",
			key:     "TEST_SLICE",
			value:   "one,,two",
			want:    []string{"one", "", "two"},
			wantErr: false,
		},
		{
			name:      "empty string",
			key:       "TEST_SLICE_EMPTY",
			value:     "",
			want:      nil,
			wantErr:   true,
			errSubstr: "is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got, err := mustExistStringSlice(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("mustExistStringSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" {
				if err != nil && !contains(err.Error(), tt.errSubstr) {
					t.Errorf("mustExistStringSlice() error = %v, should contain %v", err, tt.errSubstr)
				}
			}
			if !tt.wantErr && !sliceEqual(got, tt.want) {
				t.Errorf("mustExistStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustExistDuration(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		want      time.Duration
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid seconds",
			key:     "TEST_DURATION",
			value:   "30s",
			want:    30 * time.Second,
			wantErr: false,
		},
		{
			name:    "valid minutes",
			key:     "TEST_DURATION",
			value:   "5m",
			want:    5 * time.Minute,
			wantErr: false,
		},
		{
			name:    "valid hours",
			key:     "TEST_DURATION",
			value:   "1h",
			want:    time.Hour,
			wantErr: false,
		},
		{
			name:      "invalid duration",
			key:       "TEST_DURATION",
			value:     "notaduration",
			want:      0,
			wantErr:   true,
			errSubstr: "must be a valid duration",
		},
		{
			name:      "empty string",
			key:       "TEST_DURATION_EMPTY",
			value:     "",
			want:      0,
			wantErr:   true,
			errSubstr: "is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got, err := mustExistDuration(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("mustExistDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errSubstr != "" {
				if err != nil && !contains(err.Error(), tt.errSubstr) {
					t.Errorf("mustExistDuration() error = %v, should contain %v", err, tt.errSubstr)
				}
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("mustExistDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadEnv(t *testing.T) {
	t.Run("all env vars set", func(t *testing.T) {
		setEnvVars()
		defer unsetAllEnvVars()

		e := loadEnv()

		if e == nil {
			t.Fatal("loadEnv() returned nil")
		}
		if e.GinMode != "debug" {
			t.Errorf("GinMode = %v, want debug", e.GinMode)
		}
		if e.AppName != "testapp" {
			t.Errorf("AppName = %v, want testapp", e.AppName)
		}
		if e.ServerPort != "8080" {
			t.Errorf("ServerPort = %v, want 8080", e.ServerPort)
		}
		if len(e.ServerTrustedProxies) != 2 {
			t.Errorf("ServerTrustedProxies length = %v, want 2", len(e.ServerTrustedProxies))
		}
		if e.TwitchClientID != "test_client_id" {
			t.Errorf("TwitchClientID = %v, want test_client_id", e.TwitchClientID)
		}
		if e.TwitchClientSecret != "test_client_secret" {
			t.Errorf("TwitchClientSecret = %v, want test_client_secret", e.TwitchClientSecret)
		}
		if e.TwitchClientRedirectURI != "http://localhost/callback" {
			t.Errorf("TwitchClientRedirectURI = %v, want http://localhost/callback", e.TwitchClientRedirectURI)
		}
		if e.TwitchBroadcasterID != "test_broadcaster_id" {
			t.Errorf("TwitchBroadcasterID = %v, want test_broadcaster_id", e.TwitchBroadcasterID)
		}
		if e.TwitchBotUserID != "test_bot_user_id" {
			t.Errorf("TwitchBotUserID = %v, want test_bot_user_id", e.TwitchBotUserID)
		}
		if e.TwitchBotUserName != "testbot" {
			t.Errorf("TwitchBotUserName = %v, want testbot", e.TwitchBotUserName)
		}
		if e.TwitchKeyWordToCallBot != "!bot" {
			t.Errorf("TwitchKeyWordToCallBot = %v, want !bot", e.TwitchKeyWordToCallBot)
		}
		if e.GroqAPIKey != "test_groq_api_key" {
			t.Errorf("GroqAPIKey = %v, want test_groq_api_key", e.GroqAPIKey)
		}
		if e.GroqModel != "mixtral-8x7b-32768" {
			t.Errorf("GroqModel = %v, want mixtral-8x7b-32768", e.GroqModel)
		}
		if e.GroqModelFallback != "llama-3.3-70b-versatile" {
			t.Errorf("GroqModelFallback = %v, want llama-3.3-70b-versatile", e.GroqModelFallback)
		}
		if e.GroqMaxContextInput != 1000 {
			t.Errorf("GroqMaxContextInput = %v, want 1000", e.GroqMaxContextInput)
		}
		if e.ContextRequestDuration != 30*time.Second {
			t.Errorf("ContextRequestDuration = %v, want 30s", e.ContextRequestDuration)
		}
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("env initialized", func(t *testing.T) {
		setEnvVars()
		defer unsetAllEnvVars()

		env = loadEnv()

		e := GetEnv()
		if e == nil {
			t.Error("GetEnv() returned nil")
		}
	})
}

func TestReloadEnv(t *testing.T) {
	setEnvVars()
	defer unsetAllEnvVars()

	loadEnv()

	os.Setenv("APP_NAME", "new_app_name")
	defer os.Unsetenv("APP_NAME")

	ReloadEnv()

	e := GetEnv()
	if e.AppName != "new_app_name" {
		t.Errorf("AppName = %v, want new_app_name", e.AppName)
	}
}

func setEnvVars() {
	os.Setenv("GIN_MODE", "debug")
	os.Setenv("APP_NAME", "testapp")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_TRUSTED_PROXIES", "127.0.0.1,::1")
	os.Setenv("TWITCH_CLIENT_ID", "test_client_id")
	os.Setenv("TWITCH_CLIENT_SECRET", "test_client_secret")
	os.Setenv("TWITCH_CLIENT_REDIRECT_URI", "http://localhost/callback")
	os.Setenv("TWITCH_BROADCASTER_ID", "test_broadcaster_id")
	os.Setenv("TWITCH_BOT_USER_ID", "test_bot_user_id")
	os.Setenv("TWITCH_BOT_USER_NAME", "testbot")
	os.Setenv("TWITCH_KEY_WORD_TO_CALL_BOT", "!bot")
	os.Setenv("GROQ_API_KEY", "test_groq_api_key")
	os.Setenv("GROQ_MODEL", "mixtral-8x7b-32768")
	os.Setenv("GROQ_MODEL_FALLBACK", "llama-3.3-70b-versatile")
	os.Setenv("GROQ_MAX_CONTEXT_INPUT", "1000")
	os.Setenv("CONTEXT_REQUEST_DURATION", "30s")
}

func unsetAllEnvVars() {
	os.Unsetenv("GIN_MODE")
	os.Unsetenv("APP_NAME")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("SERVER_TRUSTED_PROXIES")
	os.Unsetenv("TWITCH_CLIENT_ID")
	os.Unsetenv("TWITCH_CLIENT_SECRET")
	os.Unsetenv("TWITCH_CLIENT_REDIRECT_URI")
	os.Unsetenv("TWITCH_BROADCASTER_ID")
	os.Unsetenv("TWITCH_BOT_USER_ID")
	os.Unsetenv("TWITCH_BOT_USER_NAME")
	os.Unsetenv("TWITCH_KEY_WORD_TO_CALL_BOT")
	os.Unsetenv("GROQ_API_KEY")
	os.Unsetenv("GROQ_MODEL")
	os.Unsetenv("GROQ_MODEL_FALLBACK")
	os.Unsetenv("GROQ_MAX_CONTEXT_INPUT")
	os.Unsetenv("CONTEXT_REQUEST_DURATION")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

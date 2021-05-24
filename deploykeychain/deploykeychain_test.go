package deploykeychain

import (
	"errors"
	"testing"
)

func TestIdenitfyRepositoryWithUnquotedArgs(t *testing.T) {
	account, repository, err := IdentifyRepository(
		[]string{
			"deploy-keychain",
			"git@github.com",
			"git-upload-pack nint8835/deploy-keychain.git",
		},
	)

	if err != nil {
		t.Errorf("identifying repo returned unexpected error: %s", err)
	}

	if account != "nint8835" {
		t.Errorf("identify repo returned unexpected account: %s", account)
	}

	if repository != "deploy-keychain" {
		t.Errorf("identify repo returned unexpected repository: %s", repository)
	}
}

func TestIdenitfyRepositoryWithQuotedArgs(t *testing.T) {
	account, repository, err := IdentifyRepository(
		[]string{
			"deploy-keychain",
			"git@github.com",
			"git-upload-pack 'nint8835/deploy-keychain.git'",
		},
	)

	if err != nil {
		t.Errorf("identifying repo returned unexpected error: %s", err)
	}

	if account != "nint8835" {
		t.Errorf("identify repo returned unexpected account: %s", account)
	}

	if repository != "deploy-keychain" {
		t.Errorf("identify repo returned unexpected repository: %s", repository)
	}
}

func TestIdenitfyRepositoryWithNoRepoArg(t *testing.T) {
	_, _, err := IdentifyRepository(
		[]string{
			"deploy-keychain",
			"git@github.com",
			"hello world",
		},
	)

	if !errors.Is(err, ErrNoRepositoryFound) {
		t.Errorf("identify repo returned unexpected error: %s", err)
	}
}

// Installing an NPM package over Git provides a leading slash on the repo name - not sure what other cases this occurs in
func TestIdenitfyRepositoryWithLeadingSlash(t *testing.T) {
	account, repository, err := IdentifyRepository(
		[]string{
			"deploy-keychain",
			"git@github.com",
			"git-upload-pack '/nint8835/deploy-keychain.git'",
		},
	)

	if err != nil {
		t.Errorf("identifying repo returned unexpected error: %s", err)
	}

	if account != "nint8835" {
		t.Errorf("identify repo returned unexpected account: %s", account)
	}

	if repository != "deploy-keychain" {
		t.Errorf("identify repo returned unexpected repository: %s", repository)
	}
}

func TestDetermineKeyFileWithCustomKey(t *testing.T) {
	keyPath, err := DetermineKeyFile(
		Config{
			Keys: map[string]string{
				"nint8835/deploy-keychain": "test.pem",
			},
		},
		"nint8835",
		"deploy-keychain",
	)

	if err != nil {
		t.Errorf("determining key file returned unexpected error: %s", err)
	}

	if keyPath != "test.pem" {
		t.Errorf("determine key file returned unexpected key path: %s", keyPath)
	}
}

func TestDetermineKeyFileWithKeyFromKeyPath(t *testing.T) {
	keyPath, err := DetermineKeyFile(
		Config{
			KeyPath:       "test_keys",
			KeyNameFormat: "{{.account}}-{{.repository}}.pem",
		},
		"nint8835",
		"deploy-keychain",
	)

	if err != nil {
		t.Errorf("determining key file returned unexpected error: %s", err)
	}

	if keyPath != "test_keys/nint8835-deploy-keychain.pem" {
		t.Errorf("determine key file returned unexpected key path: %s", keyPath)
	}
}

func TestDetermineKeyFileWithFallbackKey(t *testing.T) {
	keyPath, err := DetermineKeyFile(
		Config{
			FallbackKey:   "test.pem",
			KeyNameFormat: "{{.account}}-{{.repository}}.pem",
		},
		"nint8835",
		"deploy-keychain",
	)

	if err != nil {
		t.Errorf("determining key file returned unexpected error: %s", err)
	}

	if keyPath != "test.pem" {
		t.Errorf("determine key file returned unexpected key path: %s", keyPath)
	}
}

func TestDetermineKeyFileWithNoMatchingKeys(t *testing.T) {
	_, err := DetermineKeyFile(
		Config{
			KeyNameFormat: "{{.account}}-{{.repository}}.pem",
		},
		"nint8835",
		"deploy-keychain",
	)

	if !errors.Is(err, ErrNoKeyAvailable) {
		t.Errorf("determining key file returned unexpected error: %s", err)
	}
}

func TestDetermineKeyFileHasCorrectOrderOfPrecendence(t *testing.T) {
	config := Config{
		KeyPath:       "test_keys",
		KeyNameFormat: "{{.account}}-{{.repository}}.pem",
		FallbackKey:   "fallback.pem",
		Keys: map[string]string{
			"nint8835/doesnt-exist": "custom.pem",
		},
	}

	keyPath, err := DetermineKeyFile(
		config,
		"nint8835",
		"doesnt-exist",
	)

	if err != nil {
		t.Errorf("determining key file returned unexpected error: %s", err)
	}

	if keyPath != "custom.pem" {
		t.Errorf("determining key file did not use custom keys map as first source of keys (got key %s)", keyPath)
	}

	keyPath, err = DetermineKeyFile(
		config,
		"nint8835",
		"deploy-keychain",
	)

	if err != nil {
		t.Errorf("determining key file returned unexpected error: %s", err)
	}

	if keyPath != "test_keys/nint8835-deploy-keychain.pem" {
		t.Errorf("determining key file did not use key directory as second source of keys (got key %s)", keyPath)
	}

	keyPath, err = DetermineKeyFile(
		config,
		"nint8835",
		"no-key",
	)

	if err != nil {
		t.Errorf("determining key file returned unexpected error: %s", err)
	}

	if keyPath != "fallback.pem" {
		t.Errorf("determining key file did not use fallback key as third source of keys (got key %s)", keyPath)
	}
}

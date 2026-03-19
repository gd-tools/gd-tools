package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"hash"
	"math/big"
	"os"
	"sort"

	"github.com/tv42/zbase32"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/pbkdf2"
)

var (
	SecretsFile   = "secrets.json"
	MailUserScope = "mailuser"
)

type Secret struct {
	Scope  string `json:"scope"`
	Name   string `json:"name"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

type SecretList struct {
	Secrets []Secret `json:"entries"`
}

// GenerateSecret generates a derived secret based on the given mode.
// Supported modes are "bcrypt" (default) and "pbkdf2".
func GenerateSecret(key, mode string) (string, error) {
	switch mode {
	case "bcrypt", "":
		return GenerateBcrypt(key)
	case "pbkdf2":
		return GeneratePBKDF2(key)
	}

	return "", fmt.Errorf("unknown secret mode %s", mode)
}

// GenerateBcrypt creates a bcrypt hash from the given password.
func GenerateBcrypt(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("empty passwords in GenerateBcrypt")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// GeneratePBKDF2 derives a key using PBKDF2 with a random salt and Blake2b as hash function.
// The result is encoded using zbase32 and formatted as "$1$salt$key".
// Used for systems that require PBKDF2-style password hashing (e.g. Borg Backup).
func GeneratePBKDF2(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("empty passwords in GeneratePBKDF2")
	}

	salt := make([]byte, 20)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	h := func() hash.Hash { x, _ := blake2b.New512(nil); return x }
	key := pbkdf2.Key([]byte(password), salt, 16000, 32, h)

	return fmt.Sprintf("$1$%s$%s",
		zbase32.EncodeToString(salt),
		zbase32.EncodeToString(key),
	), nil
}

// LoadSecrets loads secrets from SecretsFile.
// Returns an empty list if the file does not exist.
func LoadSecrets() (*SecretList, error) {
	var list SecretList

	content, err := os.ReadFile(SecretsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &list, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(content, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

// Save writes the secrets to SecretsFile, sorted by scope and name.
func (list *SecretList) Save() error {
	sort.Slice(list.Secrets, func(i, j int) bool {
		if list.Secrets[i].Scope == list.Secrets[j].Scope {
			return list.Secrets[i].Name < list.Secrets[j].Name
		}
		return list.Secrets[i].Scope < list.Secrets[j].Scope
	})

	return SaveJSON(SecretsFile, list)
}

// Get returns the secret entry for the given scope and name, or nil if not found.
func (list *SecretList) Get(scope, name string) *Secret {
	for i := range list.Secrets {
		if list.Secrets[i].Scope == scope && list.Secrets[i].Name == name {
			return &list.Secrets[i]
		}
	}

	return nil
}

// SetMailUser ensures a mail user entry using bcrypt hashing.
// If password is empty, a random password is generated.
func (list *SecretList) SetMailUser(address, password string) (string, string, error) {
	if password == "" {
		newPswd, err := CreatePassword(20)
		if err != nil {
			return "", "", err
		}
		password = newPswd
	}

	output, err := GenerateSecret(password, "bcrypt")
	if err != nil {
		return "", "", err
	}

	if err := list.Set(MailUserScope, address, password, output); err != nil {
		return "", "", err
	}

	return password, output, nil
}

// Set creates or updates a secret entry and persists it.
// Existing entries are only updated if the input value changes.
func (list *SecretList) Set(scope, name, input, output string) error {
	if entry := list.Get(scope, name); entry != nil {
		if entry.Input != input {
			entry.Input = input
			entry.Output = output
		}
	} else {
		entry := Secret{
			Scope:  scope,
			Name:   name,
			Input:  input,
			Output: output,
		}
		list.Secrets = append(list.Secrets, entry)
	}

	return list.Save()
}

// CreatePassword creates a random password without visually ambiguous characters.
func CreatePassword(length int) (string, error) {
	charset := "abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ0123456789"

	pass := make([]byte, length)
	for i := range pass {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		pass[i] = charset[n.Int64()]
	}
	return string(pass), nil
}

// EnsurePassword returns a stored password or creates and stores a new one.
// The returned value is always the Output field of the entry.
// Used e.g. for common Postfix, Dovecot and Roundcube database access.
func EnsurePassword(length int, scope, name string) (string, error) {
	list, err := LoadSecrets()
	if err != nil {
		return "", err
	}

	if entry := list.Get(scope, name); entry != nil {
		return entry.Output, nil
	}

	password, err := CreatePassword(length)
	if err != nil {
		return "", err
	}

	if err := list.Set(scope, name, "", password); err != nil {
		return "", err
	}

	return password, nil
}

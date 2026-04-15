package rotate

import (
	"fmt"
	"time"

	"github.com/vaultpull/internal/audit"
	"github.com/vaultpull/internal/env"
	"github.com/vaultpull/internal/vault"
)

// Result holds the outcome of a rotation operation.
type Result struct {
	Profile   string
	Written   int
	Skipped   int
	Timestamp time.Time
}

// Rotator fetches secrets from Vault and writes them to an env file.
type Rotator struct {
	client *vault.Client
	logger *audit.Logger
}

// New creates a new Rotator.
func New(client *vault.Client, logger *audit.Logger) *Rotator {
	return &Rotator{client: client, logger: logger}
}

// Rotate fetches secrets at vaultPath and writes them to outputFile.
// It returns a Result summarising what changed.
func (r *Rotator) Rotate(profile, vaultPath, outputFile string) (*Result, error) {
	secrets, err := r.client.GetSecrets(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("rotate: fetch secrets for profile %q: %w", profile, err)
	}

	w, err := env.NewWriter(outputFile)
	if err != nil {
		return nil, fmt.Errorf("rotate: open writer for %q: %w", outputFile, err)
	}

	if err := w.Write(secrets); err != nil {
		return nil, fmt.Errorf("rotate: write env file %q: %w", outputFile, err)
	}

	for k, v := range secrets {
		if err := r.logger.Record(profile, k, v); err != nil {
			// non-fatal: log the warning but continue
			fmt.Printf("rotate: audit log warning: %v\n", err)
		}
	}

	return &Result{
		Profile:   profile,
		Written:   len(secrets),
		Timestamp: time.Now().UTC(),
	}, nil
}

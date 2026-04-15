// Package prompt provides lightweight interactive terminal utilities used by
// vaultpull commands to request user confirmation before performing
// potentially destructive operations such as overwriting an existing .env file.
//
// # Usage
//
//	term := prompt.New()
//	ok, err := term.Ask("Overwrite .env?")
//	if err != nil || !ok {
//	    return
//	}
//
// When the --force flag is active, pass a prompt.SkipConfirmer instead so
// the question is never shown to the user.
package prompt

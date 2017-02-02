// Copyright (c) 2015-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package prompt

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/decred/dcrutil/hdkeychain"
	"github.com/decred/translate_seed/walletseed"
)

// ProvideSeed is used to prompt for the wallet seed which maybe required during
// upgrades.
func ProvideSeed() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter existing wallet seed: ")
		seedStr, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		seedStr = strings.TrimSpace(strings.ToLower(seedStr))

		seed, err := hex.DecodeString(seedStr)
		if err != nil || len(seed) < hdkeychain.MinSeedBytes ||
			len(seed) > hdkeychain.MaxSeedBytes {

			fmt.Printf("Invalid seed specified.  Must be a "+
				"hexadecimal value that is at least %d bits and "+
				"at most %d bits\n", hdkeychain.MinSeedBytes*8,
				hdkeychain.MaxSeedBytes*8)
			continue
		}

		return seed, nil
	}
}

// promptList prompts the user with the given prefix, list of valid responses,
// and default list entry to use.  The function will repeat the prompt to the
// user until they enter a valid response.
func promptList(reader *bufio.Reader, prefix string, validResponses []string, defaultEntry string) (string, error) {
	// Setup the prompt according to the parameters.
	validStrings := strings.Join(validResponses, "/")
	var prompt string
	if defaultEntry != "" {
		prompt = fmt.Sprintf("%s (%s) [%s]: ", prefix, validStrings,
			defaultEntry)
	} else {
		prompt = fmt.Sprintf("%s (%s): ", prefix, validStrings)
	}

	// Prompt the user until one of the valid responses is given.
	for {
		fmt.Print(prompt)
		reply, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		reply = strings.TrimSpace(strings.ToLower(reply))
		if reply == "" {
			reply = defaultEntry
		}

		for _, validResponse := range validResponses {
			if reply == validResponse {
				return reply, nil
			}
		}
	}
}

// promptListBool prompts the user for a boolean (yes/no) with the given prefix.
// The function will repeat the prompt to the user until they enter a valid
// reponse.
func promptListBool(reader *bufio.Reader, prefix string, defaultEntry string) (bool, error) {
	// Setup the valid responses.
	valid := []string{"n", "no", "y", "yes"}
	response, err := promptList(reader, prefix, valid, defaultEntry)
	if err != nil {
		return false, err
	}
	return response == "yes" || response == "y", nil
}

// Seed prompts the user whether they want to use an existing wallet generation
// seed.  When the user answers no, a seed will be generated and displayed to
// the user along with prompting them for confirmation.  When the user answers
// yes, a the user is prompted for it.  All prompts are repeated until the user
// enters a valid response. The bool returned indicates if the wallet was
// restored from a given seed or not.
func Seed(reader *bufio.Reader) ([]byte, error) {
	var err error
	for {
		fmt.Print("Enter existing wallet seed " +
			"(followed by a blank line): ")

		// Use scanner instead of buffio.Reader so we can choose choose
		// more complicated ending condition rather than just a single
		// newline.
		var seedStr string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				break
			}
			seedStr += " " + line
		}
		seedStrTrimmed := strings.TrimSpace(seedStr)
		seedStrTrimmed = collapseSpace(seedStrTrimmed)
		wordCount := strings.Count(seedStrTrimmed, " ") + 1

		var seed []byte
		if wordCount == 1 {
			if len(seedStrTrimmed)%2 != 0 {
				seedStrTrimmed = "0" + seedStrTrimmed
			}
			seed, err = hex.DecodeString(seedStrTrimmed)
			if err != nil {
				fmt.Printf("Input error: %v\n", err.Error())
			}
		} else {
			seed, err = walletseed.DecodeFrenchUserInput(seedStrTrimmed)
			if err != nil {
				fmt.Printf("Input error: %v\n", err.Error())
			}
		}
		if err != nil || len(seed) < hdkeychain.MinSeedBytes ||
			len(seed) > hdkeychain.MaxSeedBytes {
			fmt.Printf("Invalid seed specified.  Must be a "+
				"word seed (usually 33 words) using the PGP wordlist or "+
				"hexadecimal value that is at least %d bits and "+
				"at most %d bits\n", hdkeychain.MinSeedBytes*8,
				hdkeychain.MaxSeedBytes*8)
			continue
		}

		fmt.Printf("\nSeed input successful. \nHex: %x\n", seed)

		return seed, nil
	}
}

// Setup prompts for, from a buffered reader, the private and/or public
// encryption passphrases to secure a wallet and a previously derived wallet
// seed to use, if any.  privPass and pubPass will always be non-nil values
// (private encryption is required and choosing to not use public data
// encryption will still encrypt the data with an insecure default), and a
// randomly generated seed of the recommended length will be generated and
// returned after the user has confirmed the seed has been backed up to a secure
// location.
//
// The configPubPass parameter is optional (nil should be used to represent the
// lack of a value).  When non-nil, this value represents a public passphrase
// previously specified in a configuration file.  The user will be given the
// option of using this passphrase if public data encryption is enabled,
// otherwise a user-specified passphrase will be prompted for.
func Setup(r *bufio.Reader) (seed []byte, err error) {

	// Ascertain the wallet generation seed.  This will either be an
	// automatically generated value the user has already confirmed or a
	// value the user has entered which has already been validated.
	seed, err = Seed(r)

	return
}

// collapseSpace takes a string and replaces any repeated areas of whitespace
// with a single space character.
func collapseSpace(in string) string {
	whiteSpace := false
	out := ""
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !whiteSpace {
				out = out + " "
			}
			whiteSpace = true
		} else {
			out = out + string(c)
			whiteSpace = false
		}
	}
	return out
}

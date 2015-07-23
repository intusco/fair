package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

const diceURL = "https://intus.co/dice/%s/fair/%d"

type Roll struct {
	ServerHash string
	ServerRand string
	ClientRand string
	BetValue   int64
	WinValue   int64
	RollValue  int64
}

// VerifyRoll verifies an Intus.co roll is fair.
func VerifyRoll(r *Roll) error {
	serverHash, err := hex.DecodeString(r.ServerHash)
	if err != nil {
		return err
	}

	serverRand, err := hex.DecodeString(r.ServerRand)
	if err != nil {
		return err
	}

	// Does serverRand hash match?
	calcServerHash := sha512.Sum512(serverRand)
	if !bytes.Equal(calcServerHash[:], serverHash) {
		return errors.New("serverRand does not match serverHash")
	}
	log.Println("server hashes match")

	// Calculate roll value matches.
	clientRand, err := hex.DecodeString(r.ClientRand)
	if err != nil {
		return err
	}

	// Concatenate serverRand and clientRand.
	combRand := append(serverRand, clientRand...)

	// Hash the result.
	combHash := sha512.Sum512(combRand)

	// Determine the integer value of combHash.
	combValue := new(big.Int).SetBytes(combHash[:])

	// Now mod the winValue and value to get the provably fair rollValue that is
	// of range of [0, r.WinValue).
	rollValue := new(big.Int)
	new(big.Int).DivMod(combValue, big.NewInt(r.WinValue), rollValue)

	// Check roll values match.
	if rollValue.Int64() != r.RollValue {
		return errors.New("roll values do not match")
	}

	log.Printf("roll value calculated at %v", rollValue)

	// Now see if the roll was a winning one.
	if rollValue.Int64() < r.BetValue {
		log.Printf("provably fair win of %v Satoshi", r.WinValue)
	} else {
		log.Printf("provably fair loss of %v Satoshi", r.BetValue)
	}
	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatal("address and request ID required")
	}

	address := flag.Arg(0)

	requestID, err := strconv.ParseInt(flag.Arg(1), 10, 64)
	if err != nil {
		log.Fatalf("unknown request ID %v", flag.Arg(1))
	}

	log.Printf("processing request ID %v", requestID)
	url := fmt.Sprintf(diceURL, address, requestID)
	log.Printf("connecting to %v", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("server response %v", resp.Status)
	}

	roll := &Roll{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(roll); err != nil {
		log.Fatal(err)
	}

	log.Println("Server Hash:", roll.ServerHash)
	log.Println("Client Rand:", roll.ClientRand)
	if err := VerifyRoll(roll); err != nil {
		log.Fatal(err)
	}
}

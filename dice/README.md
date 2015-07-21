# Provably Fair Dice

[dice.go](https://github.com/intusco/fair/blob/master/dice/dice.go) shows the provably fair dice role implementation for the [Intus Dice Game](https://intus.co/dice/). It is written in the Go Programming Language allowing you to inspect it, compile and run it yourself.

## How To Use
Instructions on how to install Go and compile dice.go can be found [here](https://golang.org/doc/).

Once compiled you can run the dice program as follows:

```
> dice [debit address] [request ID]
```

[debit address] - the wallet debit address used for the roll

[request ID] - the unique request ID for each roll

## Provably Fair Algorithm
1. The server generates a random number unique to each roll, hashes it using SHA512 and presents the hash to the player.
2. The player generates their own unique random number or allows the browser to do it for them.
3. The player's random number is sent to the server when they 'roll the dice' and is concatenated to the servers random number in byte form.
4. The concatenated server and player bytes are then SHA512 hashed again and converted to a number using big-endian format.
5. The number is modded so that if falls within the range [0, winValue]. Where winValue is the payout a player would get if they won.
6. As the dice game has zero edge a win is classed as any modded number below betValue. Where betValue is amount in Satoshi of the initial bet.

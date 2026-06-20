# eridol-core-go

Module eridol-core-go provides a basic interface to build an Eridol implementation. It provides oscillators and FFT suited for Eridol, so you don't have to worry about what "hearing" or "playing" a note means exactly. Take a look at the [Eridol Specification](https://docs.google.com/document/d/1AFxaWJt9EhBtvm7D-6FOtsFJlVOMukvj6GpUrZov-HY/edit?usp=sharing) to learn more.

## Getting Started
You should read through the [docs](https://pkg.go.dev/github.com/averagestardust/eridol-core-go).

The following example program keeps ownership of octave 2 by playing the counter-claim tone.
```go
	core.Run(func(heard core.Heard) (stop bool) {
			// check if someone is playing the claim tone
			isSomeoneClaiming := heard.Octaves[2].Claim

			// if someone is claiming play the counter claim tone
			core.Synth(2).SetNoteImmediately(core.CounterClaim, isSomeoneClaiming)

			// never stop running
			return false
		})
```

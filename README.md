# RFC 20 — ASCII interchange lab

**Covered RFC 20. Next RFC 21 — Network meeting.** This lab follows
[RFC 19's published scheduler](https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0019-locality-lab/commit/92cc9db1d46d19130565aa5bed220693c4d87f0c).
It reuses its Go module layout, CLI conventions, MIT license and CI workflow.

## What RFC 20 defines

On 16 October 1969, Vint Cerf proposed using the USAS X3.4-1968 ASCII
code on ARPANET host-to-host primary connections. Each eight-bit byte
carries a seven-bit ASCII value with the high bit cleared. The receiving
host defines its break character. The RFC gives SRI's period and UCLA's
carriage return as examples.

The RFC Editor lists RFC 20 as **Internet Standard (STD 80)** in the
Legacy stream. Its publication status field is Unknown and its current
status is Internet Standard. The catalog lists no updates, updated-by,
obsoletes or obsoleted-by relationships as checked on 12 November 2026.

- [Official RFC Editor text](https://www.rfc-editor.org/rfc/rfc20.html)
- [Verified errata, including carriage return hex 0D](https://www.rfc-editor.org/rfc/inline-errata/rfc20.html)
- [Official metadata](https://www.rfc-editor.org/rfc/rfc20.json)
- [RFC Editor index](https://www.rfc-editor.org/rfc-index.txt)

ASCII remains a subset of UTF-8. The byte repertoire and the boundaries
between records are separate parts of an interchange contract.
[UTF-8 specification](https://www.rfc-editor.org/info/rfc3629/)

## Applications

1. **Banking batch feed diagnostics, selected.** Check raw interchange
   bytes and demonstrate how the receiver's delimiter changes record
   boundaries before applying business validation.
2. **Gateway contract tests.** Catch accidental UTF-8 input on an
   explicitly ASCII-only interface and report the exact byte offset.
3. **Parser interoperability tests.** Replay the same payload against
   different receiver framing rules and inspect incomplete records.

## Run

Use Go 1.22 or later from the repository root.

```sh
go run ./cmd/asciilab < examples/interchange.dat
go run ./cmd/asciilab -html < examples/interchange.dat > docs/index.html
go test -race ./...
go vet ./...
go build ./...
```

The example contains the 15 bytes represented as `"PAY 12.50\rACK.\r"`.
SRI's example framing returns `"PAY 12"` and `"50\rACK"`, with `"\r"`
left over. UCLA's returns `"PAY 12.50"` and `"ACK."` with no tail.
The `.dat` fixture deliberately contains raw carriage returns.

## Behavior and scope

The CLI reads up to 64 KiB of raw bytes. Any byte in 0x80–0xFF
causes a nonzero exit and reports the first invalid byte's value and
zero-based offset on standard error. Accepted bytes pass unchanged
to the educational framing model. All 128 ASCII values are accepted,
including NUL, DEL and the other control codes.

The two receivers treat every occurrence of their delimiter as a
record boundary. Adjacent delimiters produce empty records. Bytes
after the last delimiter appear in the tail. JSON and HTML escape
control characters for inspection. The output presents historical
delimiter examples as a simple record-splitting exercise. Application
validation, escaping conventions, transport and host behavior are
outside this exercise. GitHub Pages serves a saved Go-generated report.

Tests cover all 256 possible byte values, invalid offsets, both
receiver examples, preservation of control bytes, empty records,
incomplete tails, size limits, and safe HTML escaping. CI regenerates
the published report and checks it against the committed HTML.

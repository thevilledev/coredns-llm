## CoreDNS LLM plugin

CoreDNS plugin that lets you prompt OpenAI-compatible LLM endpoints using DNS TXT queries.

Inspired by [https://llm.pieter.com](llm.pieter.com).

### Features

- OpenAI-compatible API: works with endpoints like `https://api.openai.com/v1/chat/completions` or `https://openrouter.ai/api/v1/chat/completions`.
- API key via environment: `COREDNS_LLM_API_KEY`.
- TXT response chunking to respect DNS limits (default 200 bytes per chunk).

### Build CoreDNS with this external plugin

```bash
./hack/build-coredns-with-llm.sh
```

This script clones CoreDNS, wires in the `llm` plugin, and builds a local `coredns` binary under `_build/coredns/coredns`.

### Corefile example

```
.:53 {
    llm {
        endpoint https://openrouter.ai/api/v1/chat/completions
        model openai/gpt-4o-mini
        # optional
        # chunk_size 200
        # timeout 15
    }
}
```

Set the API key via an environment variable:

```bash
export COREDNS_LLM_API_KEY=your_api_key_here
```

Run CoreDNS:

```bash
_build/coredns/coredns -conf Corefile
```

### Query examples

Using `dig`, the prompt is taken from the QNAME and decoded from standard DNS escapes. Quotes are for the shell, spaces are encoded as `\032` by `dig` and decoded by the plugin.

```bash
dig @127.0.0.1 "what is the meaning of life?" TXT +short
```

For long answers you will see multiple TXT records like `[1/N] ...`, `[2/N] ...`.

### Configuration reference

- `endpoint <url>`: Chat completions URL (OpenAI-compatible).
- `model <name>`: Model name as expected by the endpoint.
- `chunk_size <n>`: Optional per-TXT-string byte limit (default 200).
- `timeout <seconds>`: Optional HTTP timeout (default 15).

### Notes

- No rate limiting, auth scoping, or logging.
- Consider running with careful consideration. :-)



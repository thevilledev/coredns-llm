## CoreDNS LLM plugin

_It's the middle of the night. You only have a pod with `dig` installed on it._

CoreDNS plugin that lets you prompt OpenAI-compatible LLM endpoints using DNS TXT queries.

Inspired by [https://llm.pieter.com](llm.pieter.com).

### Features

- OpenAI-compatible API: works with endpoints like `https://api.openai.com/v1/chat/completions` or `https://openrouter.ai/api/v1/chat/completions`.
- API key via environment: `COREDNS_LLM_API_KEY`.
- TXT response chunking to respect DNS limits (default 200 bytes per chunk).

### Getting started

1. Setup a `Corefile`

```
.:1053 {
    llm {
        endpoint https://openrouter.ai/api/v1/chat/completions
        model openai/gpt-4o-mini
        # optional
        # chunk_size 200
        # timeout 15
    }
}
```

2. Pull the Docker image:

```bash
docker pull ghcr.io/thevilledev/coredns-llm:latest
```

3. Run it and set required configurations:

```bash
docker run \
    -w / \
    -v $(pwd)/Corefile.llm:/Corefile \
    -p 127.0.0.1:1053:1053/udp \
    -p 127.0.0.1:1053:1053 \
    -e COREDNS_LLM_API_KEY=<set-your-key-here> \
    ghcr.io/thevilledev/coredns-llm:v0.1.0
```

4. Prompt it:

```bash
$ dig -p 1053 @127.0.0.1 "how many r's in strawberry" TXT +short
"There are three \"r's\" in the word \"strawberry.\""
```

5. Prompt it longer - and notice the timeouts:

```bash
$ $dig -p 1053 @127.0.0.1 "how many r's in strawberry - a story in 1000 characters" TXT +short
;; communications error to 127.0.0.1#1053: timed out
;; communications error to 127.0.0.1#1053: timed out
"[1/6] Once upon a time in a small, whimsical village, there lived a curious girl named Lily. Unlike everyone else, she had an obsession with strawberries. She often daydreamed about strawberry fields, where"
"[2/6]  they grew abundantly under the sun. One day, while exploring a hidden grove, she stumbled upon an enchanted strawberry bush. It was said that this bush bore the sweetest strawberries, each one contai"
"[3/6] ning a wish.\010\010Excitedly, Lily picked one berry and closely examined it. \"I wonder how many r\226\128\153s are in strawberry?\" she mused. Fascinated, she counted. \"One, two, three!\" she exclaimed, \"Three r's!\" "
"[4/6] Just then, the berry glowed and whispered, \"Your curiosity brings magic.\" With each wish she made, vines of strawberries spiraled around her, creating a path to endless adventure.\010\010Soon, the village k"
"[5/6] new of her magical berries. They learned to cherish curiosity and joy. Lily, now the guardian of the enchanted bush, inspired dreams as sweet as strawberries, reminding everyone that sometimes, the si"
"[6/6] mplest questions lead to the greatest discoveries."
```
### Building

```bash
./hack/build-coredns-with-llm.sh
```

This script clones CoreDNS, wires in the `llm` plugin, and builds a local `coredns` binary under `_build/coredns/coredns`.

Set the API key via an environment variable:

```bash
export COREDNS_LLM_API_KEY=your_api_key_here
```

Run CoreDNS:

```bash
_build/coredns/coredns -conf Corefile
```

### Configuration reference

- `endpoint <url>`: Chat completions URL (OpenAI-compatible).
- `model <name>`: Model name as expected by the endpoint.
- `chunk_size <n>`: Optional per-TXT-string byte limit (default 200).
- `timeout <seconds>`: Optional HTTP timeout (default 15).

### Notes

- No rate limiting, auth scoping, or logging.
- Consider running with careful consideration. :-)



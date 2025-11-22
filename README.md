## CoreDNS LLM plugin

_It's the middle of the night. You only have a pod with `dig` installed on it._

CoreDNS plugin that lets you prompt OpenAI-compatible LLM endpoints using DNS TXT queries.

Inspired by [https://llm.pieter.com](llm.pieter.com)

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
    ghcr.io/thevilledev/coredns-llm:latest
```

4. Prompt it:

```bash
$ dig -p 1053 @127.0.0.1 "how many r's in strawberry" TXT +noall +answer
how\032many\032r's\032in\032strawberry. 0 IN TXT "There are three \"r's\" in the word \"strawberry.\""
```

5. Prompt it longer - and notice the timeouts:

```bash
$ $dig -p 1053 @127.0.0.1 "how many r's in strawberry - a story in 1000 characters" TXT +short
how\032many\032r's\032in\032strawberry\032-\032a\032story\032in\0321000\032characters. 0
IN TXT
"[1/6] In a quaint little town, a group of children gathered every summer to pick strawberries in the vibrant fields nearby. One sunny day, as they filled their baskets, they began a spirited debate: how man"
how\032many\032r's\032in\032strawberry\032-\032a\032story\032in\0321000\032characters. 0
IN TXT
"[2/6] y \"r's\" were in the word \"strawberry\"? \010\010Timmy confidently shouted, \"I believe there are two!\" His friend Emily disagreed, insisting it was three. The kids spent the afternoon counting the letters, dr"
how\032many\032r's\032in\032strawberry\032-\032a\032story\032in\0321000\032characters. 0
IN TXT
"[3/6] awing them in the dirt, and shouting out other words that included the letter \"r\". Laughter echoed through the fields as they searched for other \"r\" words, stumbling upon \"rose,\" \"rocket,\" and \"run.\" "
how\032many\032r's\032in\032strawberry\032-\032a\032story\032in\0321000\032characters. 0
IN TXT
"[4/6] \010\010As the sun began to set, they realized the strawberry patch was blooming with ripe, red fruit. They decided to celebrate their day by making strawberry shortcake. Each kid took turns mashing strawbe"
how\032many\032r's\032in\032strawberry\032-\032a\032story\032in\0321000\032characters. 0
IN TXT
"[5/6] rries and whipping cream, their friendly debate forgotten in the sweetness of their creation. \010\010In the end, it didn\226\128\153t matter how many \"r's\" were in \"strawberry\" \226\128\148 what mattered was the joy of frie"
how\032many\032r's\032in\032strawberry\032-\032a\032story\032in\0321000\032characters. 0
IN TXT
"[6/6] ndship and summer adventures."
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



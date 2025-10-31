# OTUI MCP Registry

> MCP (Model Context Protocol) server plugins registry for [OTUI](https://github.com/hkdb/otui)

## 📊 Statistics

![Plugin Count](https://img.shields.io/badge/dynamic/json?url=https://raw.githubusercontent.com/hkdb/otui-registry/main/plugins.json&query=$.length&label=plugins&color=blue)

## 📋 What is this?

This repository contains a pre-processed, enriched registry of a growing number of **MCP server plugins** for [OTUI](https://github.com/hkdb/otui). More plugins will be added as time allows.

This allows for [OTUI](https://github.com/hkdb/otui) users to easily discover, install, and use MCP servers.

## 🚀 Usage

### In OTUI

OTUI automatically fetches this registry when you press **'r'** (refresh) in the Plugin Manager.

No configuration needed - it just works!

### Manual Access

Fetch the latest registry programmatically:

```bash
curl https://raw.githubusercontent.com/hkdb/otui-registry/main/plugins.json
```

### JSON Structure

```json
[
  {
    "id": "ihor-sokoliuk-mcp-searxng",
    "name": "ihor-sokoliuk/mcp-searxng",
    "description": "MCP server for SearXNG web search integration",
    "category": "utility",
    "license": "MIT",
    "repository": "https://github.com/ihor-sokoliuk/mcp-searxng",
    "author": "ihor-sokoliuk",
    "stars": 304,
    "updated_at": "2025-10-29T05:39:20Z",
    "language": "TypeScript",
    "install_type": "npm",
    "package": "ihor-sokoliuk/mcp-searxng",
    "environment": "SEARXNG_URL",
    "args": "",
    "verified": true,
    "official": false
  },
]
```

## 🤝 Contributing

Found an incorrect entry? Want to suggest improvements?

1. Open an [issue](../../issues)
2. Submit a [pull request](../../pulls)
3. Contribute to [awesome-mcp-servers](https://github.com/punkpeye/awesome-mcp-servers) (upstream source)

## 📜 License

This registry data is compiled from public GitHub repositories and is provided as-is for convenience.

Individual plugins have their own licenses - check the `license` field in `registry.json`.

## 🔗 Related Projects

- [OTUI](https://github.com/hkdb/otui) - Terminal UI for Ollama with MCP plugin support
- [awesome-mcp-servers](https://github.com/punkpeye/awesome-mcp-servers) - Curated list of MCP servers
- [Model Context Protocol](https://modelcontextprotocol.io/) - Official MCP documentation

---

## 🚨 Disclaimer

The list of MCP servers in this registry are not all tested and verified to be safe to run. All users should do their own research and use their best judgement to stay safe when using MCP servers provided by this registry. The Author and the contributors are not held responsible by any damages caused by using the registry or the MCP servers it lists.

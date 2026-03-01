#!/usr/bin/env python3
"""
gen_docs.py - 從 translations.json + commands_meta.json 生成 docs/index.html
用法: ./gen_docs.py
"""
import json
import os
from datetime import datetime, timezone, timedelta

ROOT = os.path.join(os.path.dirname(__file__), "..")
TRANSLATIONS = os.path.join(ROOT, "conf", "translations.json")
COMMANDS_META = os.path.join(ROOT, "conf", "commands_meta.json")
OUTPUT = os.path.join(ROOT, "docs", "index.html")

LEVEL_LABEL = {
    0: "",
    1: '<span class="badge user">User</span>',
    2: '<span class="badge moderator">Moderator</span>',
    3: '<span class="badge admin">Admin</span>',
}

FUNCTION_LABEL = {
    "help":       "通用",
    "basic":      "基本",
    "ai":         "AI",
    "gamecenter": "遊戲中心",
    "gamble":     "賭場",
    "test":       "測試",
}

def load_json(path):
    with open(path, encoding="utf-8") as f:
        return json.load(f)

def build_commands(translations, meta):
    content = translations["zh"]["discord"]["api"]["help"]["content"]
    commands = []
    for name, m in meta.items():
        c = content.get(name, {})
        commands.append({
            "name":     name,
            "command":  c.get("command", f"!{name}"),
            "alias":    c.get("alias", []),
            "desc":     c.get("desc", ""),
            "example":  c.get("example", ""),
            "function": m.get("function", ""),
            "level":    m.get("level", 0),
        })
    return commands

def group_commands(commands):
    groups = {}
    for cmd in commands:
        fn = cmd["function"]
        groups.setdefault(fn, []).append(cmd)
    return groups

def render_command(cmd):
    alias_html = ""
    if cmd["alias"]:
        alias_html = f'<span class="alias">[ {" / ".join(cmd["alias"])} ]</span>'
    badge = LEVEL_LABEL.get(cmd["level"], "")
    example_html = ""
    if cmd["example"]:
        example_html = f'<div class="cmd-example">範例：<code>{cmd["example"]}</code></div>'
    return f"""
      <div class="cmd">
        <div class="cmd-header">
          <code>{cmd["command"]}</code>{alias_html}{badge}
        </div>
        <div class="cmd-desc">{cmd["desc"]}</div>
        {example_html}
      </div>"""

def render_group(fn, cmds):
    label = FUNCTION_LABEL.get(fn, fn)
    rows = "".join(render_command(c) for c in cmds)
    return f"""
    <section>
      <h2>{label}</h2>
      {rows}
    </section>"""

def render_html(groups):
    now = datetime.now(timezone(timedelta(hours=8))).strftime("%Y-%m-%d %H:%M UTC+8")
    sections = "".join(render_group(fn, cmds) for fn, cmds in groups.items())
    return f"""<!DOCTYPE html>
<html lang="zh-TW">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Bot 指令說明</title>
  <style>
    body {{ font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
           max-width: 800px; margin: 40px auto; padding: 0 20px;
           background: #1e1f22; color: #dbdee1; }}
    h1 {{ color: #ffffff; border-bottom: 2px solid #5865f2; padding-bottom: 8px; }}
    h2 {{ color: #5865f2; margin-top: 32px; }}
    .cmd {{ background: #2b2d31; border-radius: 8px; padding: 12px 16px; margin: 8px 0; }}
    .cmd-header {{ display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }}
    code {{ background: #111214; color: #00b0f4; padding: 2px 8px;
            border-radius: 4px; font-size: 0.95em; }}
    .alias {{ color: #949ba4; font-size: 0.85em; }}
    .cmd-desc {{ color: #b5bac1; margin-top: 4px; font-size: 0.9em; }}
    .cmd-example {{ color: #949ba4; margin-top: 4px; font-size: 0.85em; }}
    .badge {{ padding: 2px 6px; border-radius: 4px; font-size: 0.75em; font-weight: bold; }}
    .badge.admin {{ background: #ed4245; color: white; }}
    .badge.moderator {{ background: #f0b232; color: black; }}
    .badge.user {{ background: #23a559; color: white; }}
    footer {{ margin-top: 48px; color: #4e5058; font-size: 0.8em; text-align: center; }}
  </style>
</head>
<body>
  <h1>📖 Bot 指令說明</h1>
  {sections}
  <footer>自動生成於 {now}</footer>
</body>
</html>
"""

def main():
    translations = load_json(TRANSLATIONS)
    meta = load_json(COMMANDS_META)
    commands = build_commands(translations, meta)
    groups = group_commands(commands)

    os.makedirs(os.path.join(ROOT, "docs"), exist_ok=True)
    html = render_html(groups)
    with open(OUTPUT, "w", encoding="utf-8") as f:
        f.write(html)
    print(f"Finish：{OUTPUT}")

if __name__ == "__main__":
    main()

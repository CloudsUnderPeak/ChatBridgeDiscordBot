# go-discordbot

使用 Go 打造的多功能 Discord 機器人，具備 AI 聊天、小遊戲、賭博經濟系統等功能。

## 功能特色

- **AI 聊天** — 與 AI 助手對話，支援 OpenAI 相容的供應商（OpenAI、DeepSeek、Gemini、Grok 等），透過訊息佇列保留對話上下文。
- **小遊戲** — 猜數字（`!guess`）與 1A2B（`!1a2b`），每位使用者獨立遊戲狀態，支援多人同時遊玩。
- **賭博與經濟系統** — 比大小（`!gamble`）、拉霸機（`!slot`）、籌碼查詢（`!chips`）、排行榜（`!rank`）、籌碼轉帳（`!give`）、還款（`!repay`）。
- **基本指令** — 打招呼（`!hi`）、幫助（`!help`）。
- **多機器人支援** — 單一實例同時運行多個 Discord 機器人。
- **頻道層級功能控制** — 可在機器人層級或個別頻道層級各自指定啟用的功能模組。
- **角色權限控制** — 使用者權限等級（訪客 / 一般 / 版主 / 管理員 / 封鎖），基於中介層的身份驗證。
- **Discord 錯誤日誌** — 將 Error / Fatal / Panic 等級日誌自動推送至指定 Discord 頻道。
- **自動重連監控** — 每分鐘偵測機器人連線狀態，斷線時自動重新連線。
- **多語系** — 透過 `conf/translations.json` 支援翻譯。
- **跨平台建置** — 支援 Linux（amd64/arm/arm64）、macOS（arm64）、Windows（amd64）。

## 專案結構

```
.
├── main.go                  # 程式進入點
├── conf/
│   ├── app.yaml             # 應用程式設定（日誌、金鑰、系統語言）
│   ├── discord.yaml         # 機器人設定（Token、功能、AI 代理、賭博）
│   └── translations.json    # 多語系翻譯字串
├── discord/
│   ├── api/
│   │   ├── ai/              # AI 聊天處理器
│   │   ├── basic/           # 基本指令
│   │   ├── help/            # 幫助指令
│   │   ├── gamecenter/      # 小遊戲（猜數字、1A2B）
│   │   └── gamble/          # 賭博與經濟系統
│   ├── middleware/auth/      # 身份驗證中介層
│   └── pkg/
│       ├── discordbot/      # Discord 機器人框架
│       └── discordlogger/   # Discord 錯誤日誌 Hook
├── pkg/
│   ├── aiAgent/             # AI 客戶端抽象層（多供應商支援）
│   ├── config/              # 設定檔載入（Viper）
│   ├── logger/              # 結構化日誌（Logrus）
│   ├── signal/              # 跨平台信號處理
│   ├── sql/                 # SQLite 資料庫（GORM）
│   ├── translate/           # 翻譯管理
│   └── util/                # 工具函式
├── routers/                 # 機器人生命週期與指令路由
├── data/                    # 執行時資料（資料庫、日誌）
├── Dockerfile               # 多階段 Docker 建置
└── Makefile                 # 建置目標
```

## 指令列表

### 通用

| 指令 | 別名 | 說明 | 權限 |
|------|------|------|------|
| `!help` | `!h`, `!幫助` | 列出所有可用指令 | guest |
| `!hi` | `!嗨`, `!安` | 打招呼 | guest |

### AI 聊天（`ai` 模組）

| 指令 | 別名 | 說明 | 權限 |
|------|------|------|------|
| `!ai <message>` | `!mi` | 與 AI 聊天 | guest |

### 小遊戲（`gamecenter` 模組）

| 指令 | 別名 | 說明 | 權限 |
|------|------|------|------|
| `!guess <number>` | `!猜` | 猜數字遊戲 | guest |
| `!resetguess` | `!重猜` | 重置猜數字遊戲 | guest |
| `!1a2b <number>` | — | 1A2B 遊戲 | guest |
| `!reset1a2b` | `!重1a2b` | 重置 1A2B 遊戲 | guest |
| `!peek1a2b` | — | 偷看 1A2B 答案 | **admin** |

### 賭博與經濟（`gamble` 模組）

| 指令 | 別名 | 說明 | 權限 |
|------|------|------|------|
| `!chips` | `!代幣` | 查詢目前籌碼數 | guest |
| `!rank` | `!排行` | 前三名排行榜 | guest |
| `!gamble <amount>` | `!賭` | 比大小遊戲 | guest |
| `!slot <amount>` | `!拉霸` | 拉霸機遊戲 | guest |
| `!repay` | `!還錢` | 破產後恢復預設籌碼 | guest |
| `!give <name> <amount>` | `!給` | 轉帳籌碼給其他玩家 | **admin** |

## 前置需求

- Go 1.24+
- Discord 機器人 Token
- AI 供應商 API 金鑰（選用，僅 AI 聊天功能需要）

## 設定

### `conf/app.yaml`

```yaml
log:
  path: "./data/chatbot.log"
  level: "debug"

token:
  aiToken: ""  # 或透過環境變數設定

system:
  language: "zh"
```

### `conf/discord.yaml`

```yaml
bots:
  - name: BotName
    token: "YOUR_DISCORD_BOT_TOKEN"
    enabled: true
    helpUrl: "https://<USERNAME>.github.io/<REPO>/"  # 選用：!help 改為導向此 URL
    functions:                # 機器人層級啟用的功能（全頻道生效）
      - "basic"
      - "ai"
    logChannels:              # 選用：將 Error/Fatal/Panic 日誌推送至此頻道
      - id: "CHANNEL_ID"
    channels:                 # 選用：頻道層級功能設定（可覆蓋機器人層級）
      - id: "CHANNEL_ID_1"
        functions:
          - "basic"
          - "ai"
      - id: "CHANNEL_ID_2"
        functions:
          - "basic"
          - "gamecenter"
          - "gamble"
    aiAgent:
      provider: "openai"      # openai / deepseek / gemini / grok
      model: "gpt-4o-mini"
      queueLength: 10         # 每位使用者保留的對話歷史長度
      prompt:
        - "You are a helpful assistant."

users:                        # 選用：指定使用者的權限等級
  - id: "USER_ID"
    level: 3                  # 3 = admin

gameCenter:
  guessNumber:
    range: 100                # 猜數字範圍（預設 100，即 1~100）

gamble:
  principal: 1000             # 初始籌碼數
  biggerNumber:
    odds: 0.5                 # 莊家勝率
    minAnte: 100              # 最低下注籌碼
  slotMachine:
    minAnte: 100              # 最低下注籌碼
```

## 頻道層級功能路由

同一個機器人可針對不同頻道啟用不同功能。認證中介層會依據訊息來源頻道，判斷該功能是否在該頻道啟用：

- **機器人層級** (`functions`)：所有頻道皆可使用。
- **頻道層級** (`channels[].functions`)：僅在指定頻道啟用，與機器人層級合併計算。

## 使用者權限等級

| 等級 | 名稱 | 說明 |
|------|------|------|
| `-1` | block | 封鎖，所有指令皆無法使用 |
| `0` | guest | 一般訪客（預設） |
| `1` | user | 一般使用者 |
| `2` | moderator | 版主 |
| `3` | admin | 管理員（可使用 `!peek1a2b`、`!give`） |

未在 `users` 清單中的使用者預設為 `guest（0）`。

## Discord 錯誤日誌

設定 `logChannels` 後，機器人會將 Error / Fatal / Panic 等級的 Logrus 日誌以程式碼區塊格式推送至指定 Discord 頻道：

```
[2025-01-01 12:00:00] [ERROR] [package] file.go:42 FuncName
error message here
```

## 建置與執行

```bash
# 直接執行
make run

# 建置當前平台
make build

# 建置所有平台
make all

# 建置 Docker 映像
make docker

# 建置 Docker 映像並匯出為 tar 檔
make docker-tar
```

## Docker

```bash
# 建置
docker build -t go-discordbot:1.0.0 .

# 執行（透過環境變數傳入金鑰）
docker run -d \
  -e DOLPHINBOT_BOT_TOKEN=your_discord_bot_token \
  -e OPENAI_API_KEY=your_openai_api_key \
  go-discordbot:1.0.0
```

## 環境變數

| 變數 | 說明 |
|------|------|
| `{BOTNAME}_BOT_TOKEN` | 各機器人 Discord Token（`BOTNAME` 為 `discord.yaml` 中 `name` 欄位的大寫，例如 `DOLPHINBOT_BOT_TOKEN`） |
| `OPENAI_API_KEY` | AI 供應商 API 金鑰（覆蓋 `conf/app.yaml` 中的值） |

## 注意事項

### 拉霸機自訂表情符號（Slot Machine Symbols）

`conf/translations.json` 中的 `discord.api.gamble.slot_machine.symbols` 支援自訂 Discord 伺服器 Emoji。

**重要**：Bot 透過 API 發送訊息時，Discord **不會**自動轉換 `:emoji_name:` 格式，必須使用完整格式：

```
<:emoji_name:emoji_id>       # 靜態 Emoji
<a:emoji_name:emoji_id>      # 動態 Emoji（GIF）
```

**取得 Emoji ID 的方法**：在 Discord 聊天框輸入 `\:emoji_name:`（加反斜線），Discord 會顯示原始格式，例如 `<:pepeez:1476858900832718880>`。

**範例設定**：

```json
"symbols": [
  "<:pepeez:1476858900832718880>",
  "<:pepeok:1476859097675468810>",
  "<:pepeclown:1476859010182414347>"
]
```

## 技術棧

- [discordgo](https://github.com/bwmarrin/discordgo) — Discord API
- [go-openai](https://github.com/sashabaranov/go-openai) — OpenAI 相容客戶端
- [Viper](https://github.com/spf13/viper) — 設定檔管理
- [Logrus](https://github.com/sirupsen/logrus) — 日誌
- [GORM](https://gorm.io/) + [go-sqlite3](https://github.com/ncruces/go-sqlite3) — 資料庫

## 授權

本專案僅供私人使用。

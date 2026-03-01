ARG IMAGE_BASE=alpine:latest

# Stage 1: Build stage
FROM golang:1.24 AS builder
WORKDIR /app

COPY . .

# 根據 GOARCH 決定使用哪個 make target 編譯
# RUN if [ "$PLATFORM" = "linux/arm" ]; then \
#       make build; \
#     elif [ "$PLATFORM" = "linux/arm64" ]; then \
#       make build-arm64; \
#     elif [ "$PLATFORM" = "linux/amd64" ]; then \
#       make build-linux; \
#     else \
#       make build; \
#     fi


ARG BUILD_PLATFORM="linux/arm64"
ARG BUILD_DATE
# RUN if [ "$BUILD_PLATFORM" = "linux/arm" ]; then \
#       echo 1 > ./platform; \
#     elif [ "$BUILD_PLATFORM" = "linux/arm64" ]; then \
#       echo 2 > ./platform; \
#     elif [ "$BUILD_PLATFORM" = "linux/amd64" ]; then \
#       echo 3 > ./platform; \
#     else \
#       echo 4 > ./platform; \
#     fi

RUN make build BUILD_PLATFORM=${BUILD_PLATFORM} DATE=${BUILD_DATE}
# RUN make build BUILD_PLATFORM=${BUILD_PLATFORM}
# RUN make all

# Stage 2: Lightweight runtime environment
FROM ${IMAGE_BASE}
WORKDIR /app

RUN apk add --no-cache bash

# COPY --from=builder /app/platform .
COPY --from=builder /app/go-discordbot.app .
# COPY --from=builder /app/build/app/go-discordbot-arm.app .
# COPY --from=builder /app/build/app/go-discordbot-arm64.app .
# COPY --from=builder /app/build/app/go-discordbot-linux.app .
# COPY --from=builder /app/build/app/go-discordbot-mac.app .
COPY --from=builder /app/conf/app.yaml ./conf/app.yaml
COPY --from=builder /app/conf/discord.yaml ./conf/discord.yaml
COPY --from=builder /app/conf/translations.json ./conf/translations.json

# Remove sensitive configuration; pass credentials as environment variables during container runtime
# Note: 
# OPENAI_API_KEY=your_openai_api_key

CMD ["./go-discordbot.app"]
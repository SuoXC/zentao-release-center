#!/bin/bash
# 一键拉起：构建前端 -> 嵌入到后端二进制 -> 后端运行（同时服务 API 与前端静态资源）
# 用法：
#   ./scripts/dev.sh           # 若没构建过则构建前端，启动后端（端口 8080）
#   ./scripts/dev.sh --rebuild # 强制重新构建前端
#   ./scripts/dev.sh --no-fe   # 跳过前端构建（仅 API 模式 + 占位页提示）

set -euo pipefail

CURDIR=$(cd "$(dirname "$0")" && pwd)
ROOT=$(cd "$CURDIR/.." && pwd)
FE_DIR="$ROOT/frontend"
DIST_DIR="$FE_DIR/dist"
LOG_FILE="${LOG_FILE:-/tmp/zrc.log}"
PORT="${ZENTAO_RELEASE_SERVER_PORT:-8080}"

REBUILD=0
SKIP_FE=0
for arg in "$@"; do
  case "$arg" in
    --rebuild) REBUILD=1 ;;
    --no-fe)   SKIP_FE=1 ;;
    -h|--help)
      sed -n '2,7p' "$0" | sed 's/^# //; s/^#//'
      exit 0
      ;;
  esac
done

# 兜底：保证 dist 至少存在一个 index.html，让 go:embed 在任何分支下都能编译。
# 真正的 npm run build 产物会直接覆盖这个占位。
ensure_placeholder() {
  mkdir -p "$DIST_DIR"
  if [[ ! -f "$DIST_DIR/index.html" ]]; then
    cat > "$DIST_DIR/index.html" <<'HTML'
<!doctype html><html lang="zh-CN"><head><meta charset="UTF-8"><title>发布中心</title></head><body><div id="app">前端未构建。请运行 npm run build。</div></body></html>
<!-- placeholder:true -->
HTML
  fi
}

ensure_placeholder

if [[ "$SKIP_FE" != "1" ]]; then
  if [[ -f "$DIST_DIR/index.html" ]] && grep -q 'placeholder:true' "$DIST_DIR/index.html" && [[ "$REBUILD" != "1" ]]; then
    echo "[dev.sh] 未检测到构建产物，开始构建前端"
  fi
  if [[ ! -d "$FE_DIR/node_modules" ]]; then
    echo "[dev.sh] 安装前端依赖 ..."
    (cd "$FE_DIR" && npm install)
  fi
  echo "[dev.sh] 构建前端 (npm run build) ..."
  (cd "$FE_DIR" && npm run build)
else
  echo "[dev.sh] --no-fe，跳过前端构建（将使用占位 index.html）"
fi

pkill -f 'go-build.*zentao-release-center' 2>/dev/null || true
pkill -f 'bin/release-center'             2>/dev/null || true
pkill -f 'go run \.'                       2>/dev/null || true
sleep 1

echo "[dev.sh] 启动后端（日志: $LOG_FILE）..."
cd "$ROOT"
setsid nohup env "ZENTAO_RELEASE_SERVER_PORT=$PORT" go run . </dev/null >"$LOG_FILE" 2>&1 &

sleep 4

if ! curl -sf -m 2 "http://127.0.0.1:$PORT/api/projects/repos" -o /dev/null; then
  echo "[dev.sh] 后端启动失败，请查看日志: $LOG_FILE"
  tail -60 "$LOG_FILE" || true
  exit 1
fi

echo "[dev.sh] ✅ 已启动"
echo "  API:    http://127.0.0.1:$PORT/api/..."
echo "  前端:   http://127.0.0.1:$PORT/  （SPA fallback 到 index.html）"
echo "  日志:   $LOG_FILE"
echo
echo "停止: pkill -f zentao-release-center"

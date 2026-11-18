#!/bin/bash
# 从已下载的 Recovery ROM 中提取 bootloader 相关分区（abl/xbl），用于逆向分析。
# 用法：bash scripts/extract-abl.sh
set -euo pipefail

# ---- 顶部配置 ----
ROM_DIR="${ROM_DIR:-$HOME/Downloads/KSUForge/rom/dada/OS4.0.4.0.XOCCNXM}"   # ROM zip 所在目录
OUT_DIR="${OUT_DIR:-$HOME/Downloads/KSUForge/abl}"                            # 提取输出目录
PARTS="${PARTS:-abl,xbl,xbl_config,abl_config}"                              # 要提取的分区（逗号分隔）
# ------------------

REPO="$(cd "$(dirname "$0")/.." && pwd)"
DUMPER="$REPO/bin/payload-dumper"
[ -x "$DUMPER" ] || { echo "找不到 $DUMPER"; exit 1; }

ZIP=$(ls -1 "$ROM_DIR"/*.zip 2>/dev/null | head -1)
[ -n "$ZIP" ] || { echo "在 $ROM_DIR 下没找到 .zip，请改 ROM_DIR"; exit 1; }

mkdir -p "$OUT_DIR"
echo "ROM: $ZIP"
echo "=== payload 内分区列表 ==="
"$DUMPER" list "$ZIP"

echo "=== 提取分区 ==="
for p in ${PARTS//,/ }; do
  if "$DUMPER" extract "$ZIP" -p "$p" -o "$OUT_DIR"; then
    echo "已提取 $p"
  else
    echo "跳过 $p（提取失败或不存在）"
  fi
done

echo "=== 结果 ==="
ls -la "$OUT_DIR"

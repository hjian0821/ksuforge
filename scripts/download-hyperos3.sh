#!/bin/bash
# 下载 HyperOS 3 线刷包（fastboot ROM）并提取 bootloader 相关 ELF，供逆向对比。
# 线刷包内含 images/abl.elf，比从 recovery 卡刷包 payload 里抠更直接。
set -euo pipefail

# ---- 顶部配置 ----
URL="https://bn.d.miui.com/OS3.0.302.0.WOCCNXM/dada_images_OS3.0.302.0.WOCCNXM_20260407.0000.00_16.0_cn_03b5f57e46.tgz"
TGZ="$HOME/Downloads/KSUForge/hyperos3/images.tgz"     # 下载位置（约 12.9GB）
OUT_DIR="$HOME/Downloads/KSUForge/abl-hyperos3"        # abl/xbl 提取输出目录
FILES="images/abl.elf images/xbl.elf images/xbl_config.elf images/abl_config.elf images/aop.elf images/devcfg.elf"
# ------------------

mkdir -p "$(dirname "$TGZ")" "$OUT_DIR"

echo "== 下载线刷包（断点续传，约 12.9GB）=="
curl -L -C - --retry 3 -o "$TGZ" "$URL"

echo "== 包内 bootloader 相关文件 =="
tar -tzf "$TGZ" 2>/dev/null | grep -iE 'images/(abl|xbl|aop|devcfg)[^/]*\.(elf|img|mbn)$' || true

echo "== 解压 =="
for f in $FILES; do
  if tar -xzf "$TGZ" -C "$OUT_DIR" "$f" 2>/dev/null; then
    echo "已提取 $f"
  else
    echo "跳过 $f（不存在）"
  fi
done

echo "== 结果 =="
find "$OUT_DIR" -type f -exec ls -la {} \;

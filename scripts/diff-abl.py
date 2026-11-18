#!/usr/bin/env python3
"""对比 HyperOS 3 vs 4 的 abl.elf：字节级差异区间 + 字符串级增删，定位 USB 初始化改动。

用法: python3 scripts/diff-abl.py <os3.elf> <os4.elf>
"""
import re
import sys

USB_KW = [b"usb", b"Usb", b"USB", b"dwc", b"Dwc", b"DWC", b"gadget",
          b"fastboot", b"Fastboot", b"qusb", b"phy", b"super", b"high",
          b"full", b"speed", b"ep0", b"config", b"charger"]


def read(path):
    with open(path, "rb") as f:
        return f.read()


def strset(data):
    return set(re.findall(rb"[ -~]{4,}", data))


def is_usb(s):
    return any(k in s for k in USB_KW)


def changed_ranges(a, b):
    """返回 (start, end) 字节差异区间，间隔 <= 32 字节的合并。"""
    n = min(len(a), len(b))
    ranges = []
    i = 0
    while i < n:
        if a[i] != b[i]:
            j = i
            while j < n and a[j] != b[j]:
                j += 1
            if ranges and i - ranges[-1][1] <= 32:
                ranges[-1] = (ranges[-1][0], j)
            else:
                ranges.append((i, j))
            i = j
        else:
            i += 1
    return ranges


def main():
    if len(sys.argv) != 3:
        print("用法: python3 scripts/diff-abl.py <os3.elf> <os4.elf>")
        sys.exit(1)
    a = read(sys.argv[1])
    b = read(sys.argv[2])
    print(f"os3 size = {len(a)}  os4 size = {len(b)}  ({len(b)-len(a):+d})")

    # 字符串级增删（地址无关，最能反映逻辑改动）
    sa, sb = strset(a), strset(b)
    added = [s for s in (sb - sa) if is_usb(s)]
    removed = [s for s in (sa - sb) if is_usb(s)]
    print("\n== os4 新增的 USB/充电相关字符串 ==")
    for s in sorted(added)[:60]:
        print("  +", s.decode("latin1"))
    print("== os4 删除的 USB/充电相关字符串 ==")
    for s in sorted(removed)[:60]:
        print("  -", s.decode("latin1"))

    # 字节级差异区间（定位改动位置，供 IDA 跳转）
    print("\n== 字节差异区间（>8 字节，含附近字符串）==")
    for s, e in changed_ranges(a, b):
        if e - s <= 8:
            continue
        region = b[max(0, s - 64):min(len(b), e + 64)]
        strs = re.findall(rb"[ -~]{4,}", region)
        strs = [x for x in strs if is_usb(x)]
        mark = "  <-- USB" if strs else ""
        print(f"0x{s:08x}-0x{e:08x} ({e-s} bytes){mark}")
        for x in strs[:4]:
            print(f"      {x.decode('latin1')}")


if __name__ == "__main__":
    main()

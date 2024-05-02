#!/usr/bin/env bash

set -euo pipefail
shopt -s inherit_errexit

# Where the rootfs starts in MB
readonly rootfs_start=1

in=$1
out=$2
tmdir=$(mktemp -d)
trap 'rm -rf $tmdir' EXIT
rootfs=$tmdir/01_rootfs
hashtree=$tmdir/02_verity_hashtree
dm_verity_file=$out/dm_verity.txt
roothash=$out/roothash
raw=$out/raw.img

uuid=c1b9d5a2-f162-11cf-9ece-0020afc76f16
salt=0102030405060708090a0b0c0d0e0f

if [ -z "${SOURCE_DATE_EPOCH}" ]; then
    echo "SOURCE_DATE_EPOCH is not set" >&2
    exit 1
fi

mkdir -p "$out"

export E2FSPROGS_FAKE_TIME=$SOURCE_DATE_EPOCH
export MKE2FS_DEVICE_SECTSIZE=512

# guess the size of the rootfs
size=128

truncate -s "${size}M" "$rootfs"

until mkfs.ext4 \
    -E "hash_seed=$uuid" \
    -L root \
    -U $uuid -I 256 \
    -m 0 \
    -b 4096 \
    -T default \
    -d "$in" \
    "$rootfs"
do
    size=$((size + 128))
    rm -f "$rootfs"
    truncate -s "${size}M" "$rootfs"
done

verity_out=$(veritysetup format "$rootfs" "$hashtree" --no-superblock --uuid "$uuid" --salt "$salt" | tee "$dm_verity_file")
root_hash=$(echo "$verity_out" | grep -oP 'Root hash:\s+\K\w{64}' | tr -d "[:space:]")
echo -n "$root_hash" > "$roothash"

# full image size is dos header + rootfs + hashtree
hashtree_size_bytes=$(stat -c %s "$hashtree")
img_size=$((512 + "$size" * 1024 * 1024 + "$hashtree_size_bytes"))

# hash_start is the start of the hashtree in MB
hash_start=$((rootfs_start + size ))
hash_end=$((hash_start + hashtree_size_bytes / 1048576 ))

# create the raw image
truncate -s "$img_size" "$raw"

# create the partition table
parted -s -a optimal "${raw}" -- \
	mklabel msdos \
    mkpart primary ext4 ${rootfs_start}M ${hash_start}M \
    mkpart primary ext4 ${hash_start}M ${hash_end}M
sfdisk --disk-id "${raw}" 0

# write the rootfs and hashtree to the raw image
dd if="$rootfs" of="${raw}" bs=1M seek=${rootfs_start} conv=notrunc,fsync
dd if="$hashtree" of="${raw}" bs=1M seek=${hash_start} conv=notrunc,fsync

sfdisk -J "${raw}" > "${out}/sfdisk.json"

# TODO: remove
cp "$rootfs" "$out/rootfs"
cp "$hashtree" "$out/hashtree"

# TODO: igvm file https://github.com/microsoft/igvm-tooling

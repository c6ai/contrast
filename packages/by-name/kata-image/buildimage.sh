#!/usr/bin/env bash

set -euo pipefail
shopt -s inherit_errexit

# Image layout:
#
#   +---------------------------------+-------------------+-------------------------+
#   | 512B DOS MBR (padded to 1 MiB)  | p0    rootfs      | p1      hashtree        |
#   +---------------------------------+-------------------+-------------------------+
#   |                                 |                   |                         |
#   0                                 1MiB                1MiB + rootfs_size        1MiB + rootfs_size + hashtree_size

# rootfs: erofs filesystem mounted at / (read-only)
# hashtree: dm-verity hashtree without superblock

readonly MIB=1048576

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
# size=128

# truncate -s "${size}M" "$rootfs"

# # TODO: remove ext4
# until mkfs.ext4 \
#     -E "hash_seed=$uuid" \
#     -L root \
#     -U $uuid -I 256 \
#     -m 0 \
#     -b 4096 \
#     -T default \
#     -d "$in" \
#     "$rootfs"
# do
#     size=$((size + 128))
#     rm -f "$rootfs"
#     truncate -s "${size}M" "$rootfs"
# done

# create the rootfs and pad it to 1MiB
mkfs.erofs \
    -z lz4 \
    -b 4096 \
    -T "$SOURCE_DATE_EPOCH" \
    -U "$uuid" \
    --tar=f \
    "$rootfs" \
    "$in"
truncate -s '%1MiB' "$rootfs"

verity_out=$(veritysetup format "$rootfs" "$hashtree" --data-block-size 4096 --hash-block-size 4096 --no-superblock --uuid "$uuid" --salt "$salt" | tee "$dm_verity_file")
truncate -s '%1MiB' "$hashtree"
sed -i 1d "$dm_verity_file"
root_hash=$(echo "$verity_out" | grep -oP 'Root hash:\s+\K\w+' | tr -d "[:space:]")
echo -n "$root_hash" > "$roothash"
hash_type=$(echo "$verity_out" | grep -oP 'Hash type:\s+\K\w+' | tr -d "[:space:]")
echo -n "$hash_type" > "$out/hash_type"
data_blocks=$(echo "$verity_out" | grep -oP 'Data blocks:\s+\K\w+' | tr -d "[:space:]")
echo -n "$data_blocks" > "$out/data_blocks"
data_block_size=$(echo "$verity_out" | grep -oP 'Data block size:\s+\K\w+' | tr -d "[:space:]")
echo -n "$data_block_size" > "$out/data_block_size"
hash_blocks=$(echo "$verity_out" | grep -oP 'Hash blocks:\s+\K\w+' | tr -d "[:space:]")
echo -n "$hash_blocks" > "$out/hash_blocks"
hash_block_size=$(echo "$verity_out" | grep -oP 'Hash block size:\s+\K\w+' | tr -d "[:space:]")
echo -n "$hash_block_size" > "$out/hash_block_size"
hash_algorithm=$(echo "$verity_out" | grep -oP 'Hash algorithm:\s+\K\w+' | tr -d "[:space:]")
echo -n "$hash_algorithm" > "$out/hash_algorithm"
echo -n "$salt" > "$out/salt"

# TODO: calculate correct size based on what dm-verity expects
rootfs_size_mib=$(($(stat -c %s "$rootfs") / "$MIB"))
# full image size is dos header + rootfs + hashtree
hashtree_size_bytes=$(stat -c %s "$hashtree")
hashtree_size_mib=$(($(stat -c %s "$hashtree") / "$MIB"))
# img_size is the size of the full image in bytes
# DOS MBR (padded to 1MiB) + rootfs + hashtree
img_size_bytes=$(("$MIB" + "$rootfs_size_mib" * "$MIB" + "$hashtree_size_bytes" + "$MIB"))

# Where the rootfs starts in MiB
readonly rootfs_start=1
# hash_start is the start of the hashtree in MiB
hash_start=$(( rootfs_start + rootfs_size_mib ))
hash_end=$(( hash_start + hashtree_size_mib ))

rs=$(printf "%4dMiB" "$rootfs_start")
hs=$(printf "%4dMiB" "$hash_start")
he=$(printf "%4dMiB" "$hash_end")
cat << EOF
Image layout:

  +---------------------------------+-------------------+-------------------------+
  | 512B DOS MBR (padded to 1 MiB)  | p0    rootfs      | p1      hashtree        |
  +---------------------------------+-------------------+-------------------------+
  |                                 |                   |                         |
  0                             $rs             $hs                   $he
EOF

# create the raw image
truncate -s "$img_size_bytes" "$raw"

# create the partition table
parted -s -a optimal "${raw}" -- \
	mklabel msdos \
    mkpart primary ext4 ${rootfs_start}MiB ${hash_start}MiB \
    mkpart primary ext4 ${hash_start}MiB '100%'
sfdisk --disk-id "${raw}" 0

# write the rootfs and hashtree to the raw image
dd if="$rootfs" of="${raw}" bs=1MiB seek=${rootfs_start} conv=notrunc,fsync
dd if="$hashtree" of="${raw}" bs=1MiB seek=${hash_start} conv=notrunc,fsync

# TODO: remove
sfdisk -J "${raw}" > "${out}/sfdisk.json"
cp "$rootfs" "$out/rootfs"
cp "$hashtree" "$out/hashtree"

# TODO: igvm file https://github.com/microsoft/igvm-tooling

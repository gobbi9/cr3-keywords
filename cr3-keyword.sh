cr3_to_jpgs() {
  mkdir -p jpgs

  for f in "$@"; do
    [[ "$f" == *.CR3 ]] || continue

    local base="${f:t:r}"
    local out="jpgs/${base}.jpg"

    if [[ -f "$out" ]]; then
      echo "Skipping $f (JPG exists)"
      continue
    fi

    echo "Processing $f..."

    local tmp="jpgs/${base}_tmp.jpg"

    exiftool -b -PreviewImage "$f" > "$tmp"

    local orientation
    orientation=$(exiftool -Orientation -n -s -s -s "$f")

    case "$orientation" in
      3) magick "$tmp" -rotate 180 -resize 2000x2000 -quality 85 "$out" ;;
      6) magick "$tmp" -rotate 90  -resize 2000x2000 -quality 85 "$out" ;;
      8) magick "$tmp" -rotate 270 -resize 2000x2000 -quality 85 "$out" ;;
      *) magick "$tmp" -resize 2000x2000 -quality 85 "$out" ;;
    esac

    rm "$tmp"
  done
}

lm_caption_single() {
  local img="$1"
  local model="${2:-gemma-4-e4b}"
  local prompt_file="${3:-prompt.md}"

  [[ -f "$img" ]] || {
    echo "❌ Image not found: $img"
    return 1
  }

  [[ -f "$prompt_file" ]] || {
    echo "❌ Prompt file not found: $prompt_file"
    return 1
  }

  mkdir -p outputs tmp

  local base="${img:t:r}"
  local out="outputs/${base}.txt"

  if [[ -f "$out" ]]; then
    echo "Skipping $img (TXT exists)"
    return 0
  fi

  echo "Processing $img..."

  local b64_file="tmp/${base}.b64"
  local json_file="tmp/${base}.json"

  # Base64 encode image
  base64 -i "$img" | tr -d '\n' > "$b64_file"

  # Build JSON payload
  cat > "$json_file" <<EOF
{
  "model": "$model",
  "messages": [
    {
      "role": "user",
      "content": [
        {
          "type": "text",
          "text": $(jq -Rs . < "$prompt_file")
        },
        {
          "type": "image_url",
          "image_url": {
            "url": "data:image/jpeg;base64,$(cat "$b64_file")"
          }
        }
      ]
    }
  ]
}
EOF

  # Send request
  curl -s http://localhost:1234/v1/chat/completions \
    -H "Content-Type: application/json" \
    -d @"$json_file" | \
    jq -r '.choices[0].message.content' > "$out"

  rm -f "$b64_file" "$json_file"
}

lm_batch_caption() {
  local model="${1:-gemma-4-e4b}"
  local prompt_file="${2:-prompt.md}"
  local jobs="${3:-2}"

  find jpgs -name "*.jpg" | while read -r img; do
    while [[ $(jobs -r | wc -l) -ge $jobs ]]; do
      sleep 1
    done

    lm_caption_single "$img" "$model" "$prompt_file" &
  done

  wait

  echo "✅ All jobs completed"
}

txt_to_xmp() {
  local txt_dir="${1:-outputs}"
  local img_dir="${2:-.}"

  for txt in "$txt_dir"/*.txt; do
    [[ -e "$txt" ]] || { echo "⚠️ No txt files"; return 1; }

    local base="${txt:t:r}"
    local xmp_file="$img_dir/$base.xmp"

    if [[ -f "$xmp_file" ]]; then
      echo "Skipping $base (XMP exists)"
      continue
    fi

    local keywords
    keywords=$(head -n 1 "$txt")

    local caption
    caption=$(awk 'NR>1 && found {print} /^$/ {found=1}' "$txt")
    [[ -z "$caption" ]] && caption=$(tail -n +2 "$txt")

    local kw_array
    kw_array=("${(@s:,:)keywords}")

    local xmp_keywords=""
    for kw in "${kw_array[@]}"; do
      local kw_clean
      kw_clean=$(echo "$kw" | xargs)
      [[ -n "$kw_clean" ]] && xmp_keywords+="<rdf:li>$kw_clean</rdf:li>"
    done

    echo "Writing $xmp_file"

    cat > "$xmp_file" <<EOF
<x:xmpmeta xmlns:x="adobe:ns:meta/">
 <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:Description rdf:about=""
    xmlns:dc="http://purl.org/dc/elements/1.1/"
    xmlns:lr="http://ns.adobe.com/lightroom/1.0/">

    <lr:description>$caption</lr:description>

    <dc:description>
      <rdf:Alt>
        <rdf:li xml:lang="x-default">$caption</rdf:li>
      </rdf:Alt>
    </dc:description>

    <dc:subject>
      <rdf:Bag>
        $xmp_keywords
      </rdf:Bag>
    </dc:subject>

  </rdf:Description>
 </rdf:RDF>
</x:xmpmeta>
EOF
  done
}

cr3_to_xmp_pipeline() {
  local model="${1:-gemma-4-e4b}"
  local prompt_file="${2:-prompt.md}"
  shift 2

  local files=("$@")

  if [[ ${#files[@]} -eq 0 ]]; then
    files=( *.CR3 )
  fi

  echo "=== Step 1: CR3 → JPG ==="
  cr3_to_jpgs "${files[@]}" || return 1

  echo "=== Step 2: JPG → TXT ==="
  lm_batch_caption "$model" "$prompt_file" || return 1

  echo "=== Step 3: TXT → XMP ==="
  txt_to_xmp outputs . || return 1

  echo "=== Done ==="
}

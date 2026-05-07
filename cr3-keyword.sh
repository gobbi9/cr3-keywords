#!/usr/bin/env zsh

set -u

VERBOSE=0
DRY_RUN=0

log_verbose() {
  (( VERBOSE == 1 )) && echo "$*"
}

expand_path() {
  local p="$1"
  if [[ "$p" == '~/'* ]]; then
    p="$HOME/${p#\~/}"
  elif [[ "$p" == '~' ]]; then
    p="$HOME"
  fi
  print -r -- "$p"
}

xml_escape() {
  local s="$1"
  s=${s//&/&amp;}
  s=${s//</&lt;}
  s=${s//>/&gt;}
  s=${s//\"/&quot;}
  s=${s//\'/&apos;}
  print -r -- "$s"
}

render_progress() {
  local current="$1"
  local total="$2"
  local width=30

  (( total <= 0 )) && total=1

  local filled=$(( current * width / total ))
  local empty=$(( width - filled ))

  local filled_bar
  local empty_bar
  filled_bar=$(printf '%*s' "$filled" '' | tr ' ' '#')
  empty_bar=$(printf '%*s' "$empty" '' | tr ' ' '-')

  printf '\r[%s%s] %d/%d' "$filled_bar" "$empty_bar" "$current" "$total"
  (( current == total )) && printf '\n'
}

cr3_to_jpgs() {
  local jpg_dir="$1"
  shift

  local -a files
  files=("$@")

  if [[ ${#files[@]} -eq 0 ]]; then
    echo "❌ No CR3 files to convert"
    return 1
  fi

  (( DRY_RUN == 0 )) && mkdir -p "$jpg_dir"

  local total=${#files[@]}
  local idx=0

  for f in "${files[@]}"; do
    ((idx++))

    [[ -f "$f" ]] || {
      log_verbose "Skipping missing file: $f"
      render_progress "$idx" "$total"
      continue
    }

    local base="${f:t:r}"
    local out="$jpg_dir/${base}.jpg"

    if [[ -f "$out" ]]; then
      log_verbose "Skipping $f (JPG exists)"
      render_progress "$idx" "$total"
      continue
    fi

    if (( DRY_RUN == 1 )); then
      log_verbose "[dry-run] Would generate JPG: $out"
      render_progress "$idx" "$total"
      continue
    fi

    local tmp="$jpg_dir/${base}_tmp.jpg"

    if ! exiftool -b -PreviewImage "$f" > "$tmp"; then
      echo "❌ Failed to extract preview image from: $f"
      rm -f "$tmp"
      render_progress "$idx" "$total"
      continue
    fi

    local orientation
    orientation=$(exiftool -Orientation -n -s -s -s "$f" 2>/dev/null)

    case "$orientation" in
      3) magick "$tmp" -rotate 180 -resize 2000x2000 -quality 85 "$out" ;;
      6) magick "$tmp" -rotate 90  -resize 2000x2000 -quality 85 "$out" ;;
      8) magick "$tmp" -rotate 270 -resize 2000x2000 -quality 85 "$out" ;;
      *) magick "$tmp" -resize 2000x2000 -quality 85 "$out" ;;
    esac

    rm -f "$tmp"
    render_progress "$idx" "$total"
  done
}

lm_caption_single() {
  local img="$1"
  local model="$2"
  local prompt_file="$3"
  local output_dir="$4"
  local tmp_dir="$5"

  [[ -f "$img" ]] || {
    echo "❌ Image not found: $img"
    return 1
  }

  [[ -f "$prompt_file" ]] || {
    echo "❌ Prompt file not found: $prompt_file"
    return 1
  }

  local base="${img:t:r}"
  local out="$output_dir/${base}.txt"

  if [[ -f "$out" ]]; then
    log_verbose "Skipping $img (TXT exists)"
    return 0
  fi

  if (( DRY_RUN == 1 )); then
    log_verbose "[dry-run] Would send request for $img and write $out"
    return 0
  fi

  mkdir -p "$output_dir" "$tmp_dir"

  local b64_file="$tmp_dir/${base}.b64"
  local json_file="$tmp_dir/${base}.json"

  base64 -i "$img" | tr -d '\n' > "$b64_file"

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

  curl -s http://localhost:1234/v1/chat/completions \
    -H "Content-Type: application/json" \
    -d @"$json_file" | \
    jq -r '.choices[0].message.content' > "$out"

  rm -f "$b64_file" "$json_file"
}

lm_batch_caption() {
  local jpg_dir="$1"
  local output_dir="$2"
  local tmp_dir="$3"
  local model="${4:-gemma-4-e4b}"
  local prompt_file="${5:-prompt.md}"
  shift 5

  local -a images
  if [[ $# -gt 0 ]]; then
    images=("$@")
  else
    images=("$jpg_dir"/*.jpg(N))
  fi

  if [[ ${#images[@]} -eq 0 ]]; then
    echo "❌ No JPG files to caption"
    return 1
  fi

  local total=${#images[@]}
  local idx=0

  for img in "${images[@]}"; do
    ((idx++))

    [[ -f "$img" ]] || {
      log_verbose "Skipping missing JPG: $img"
      render_progress "$idx" "$total"
      continue
    }

    lm_caption_single "$img" "$model" "$prompt_file" "$output_dir" "$tmp_dir"
    render_progress "$idx" "$total"
  done
}

txt_to_xmp() {
  local txt_dir="$1"
  local xmp_dir="$2"
  shift 2

  local -a txt_files
  if [[ $# -gt 0 ]]; then
    txt_files=("$@")
  else
    txt_files=("$txt_dir"/*.txt(N))
  fi

  if [[ ${#txt_files[@]} -eq 0 ]]; then
    echo "❌ No TXT files to convert"
    return 1
  fi

  local total=${#txt_files[@]}
  local idx=0

  for txt in "${txt_files[@]}"; do
    ((idx++))

    [[ -f "$txt" ]] || {
      log_verbose "Skipping missing TXT: $txt"
      render_progress "$idx" "$total"
      continue
    }

    local base="${txt:t:r}"
    local xmp_file="$xmp_dir/$base.xmp"

    if [[ -f "$xmp_file" ]]; then
      log_verbose "Skipping $base (XMP exists)"
      render_progress "$idx" "$total"
      continue
    fi

    if (( DRY_RUN == 1 )); then
      log_verbose "[dry-run] Would write XMP: $xmp_file"
      render_progress "$idx" "$total"
      continue
    fi

    local keywords
    keywords=$(head -n 1 "$txt")

    local caption
    caption=$(awk 'NR>1 && found {print} /^$/ {found=1}' "$txt")
    [[ -z "$caption" ]] && caption=$(tail -n +2 "$txt")

    local caption_escaped
    caption_escaped=$(xml_escape "$caption")

    local keyword_items
    keyword_items=("${(@s:,:)keywords}")

    local xmp_keywords=""
    local keyword
    for keyword in "${keyword_items[@]}"; do
      keyword="${${keyword##[[:space:]]#}%%[[:space:]]#}"
      [[ -n "$keyword" ]] || continue
      xmp_keywords+="<rdf:li>$(xml_escape "$keyword")</rdf:li>"
    done

    cat > "$xmp_file" <<EOF
<x:xmpmeta xmlns:x="adobe:ns:meta/">
 <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:Description rdf:about=""
    xmlns:dc="http://purl.org/dc/elements/1.1/"
    xmlns:lr="http://ns.adobe.com/lightroom/1.0/">

    <lr:description>$caption_escaped</lr:description>

    <dc:description>
      <rdf:Alt>
        <rdf:li xml:lang="x-default">$caption_escaped</rdf:li>
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

    render_progress "$idx" "$total"
  done
}

cr3_to_xmp_pipeline() {
  local model="$1"
  local prompt_file="$2"
  local cr3_path="$3"
  shift 3

  [[ -d "$cr3_path" ]] || {
    echo "❌ CR3 path not found: $cr3_path"
    return 1
  }

  [[ -f "$prompt_file" ]] || {
    echo "❌ Prompt file not found: $prompt_file"
    return 1
  }

  local folder_name="${cr3_path:A:t}"
  local jpg_dir="jpgs/$folder_name"
  local output_dir="outputs/$folder_name"
  local tmp_dir="tmp/$folder_name"

  if (( DRY_RUN == 0 )); then
    mkdir -p "$jpg_dir" "$output_dir" "$tmp_dir"
  fi

  local -a files
  files=("$@")

  local -a cr3_files
  if [[ ${#files[@]} -eq 0 ]]; then
    cr3_files=("$cr3_path"/*.CR3(N) "$cr3_path"/*.cr3(N))
  else
    local item
    for item in "${files[@]}"; do
      local candidate
      if [[ "$item" == */* || "$item" == /* || "$item" == '~'* ]]; then
        candidate="$item"
      else
        candidate="$cr3_path/$item"
      fi

      candidate=$(expand_path "$candidate")

      if [[ -f "$candidate" ]]; then
        cr3_files+=("$candidate")
      else
        log_verbose "Skipping missing CR3: $candidate"
      fi
    done
  fi

  if [[ ${#cr3_files[@]} -eq 0 ]]; then
    echo "❌ No CR3 files found to process"
    return 1
  fi

  local -a selected_jpgs
  local -a selected_txts
  local cr3_file
  for cr3_file in "${cr3_files[@]}"; do
    local base="${cr3_file:t:r}"
    selected_jpgs+=("$jpg_dir/$base.jpg")
    selected_txts+=("$output_dir/$base.txt")
  done

  echo "=== Step 1: CR3 → JPG ==="
  cr3_to_jpgs "$jpg_dir" "${cr3_files[@]}" || return 1

  echo "=== Step 2: JPG → TXT ==="
  lm_batch_caption "$jpg_dir" "$output_dir" "$tmp_dir" "$model" "$prompt_file" "${selected_jpgs[@]}" || return 1

  echo "=== Step 3: TXT → XMP ==="
  txt_to_xmp "$output_dir" "$cr3_path" "${selected_txts[@]}" || return 1

  echo "=== Done ==="
}

usage() {
  cat <<'EOF'
Usage:
  ./cr3-keyword.sh [--verbose] [--dry-run] <cr3_path>
  ./cr3-keyword.sh [--verbose] [--dry-run] <cr3_path> IMG_0150.CR3
  ./cr3-keyword.sh [--verbose] [--dry-run] <cr3_path> IMG_0150.CR3 IMG_0151.CR3
  ./cr3-keyword.sh [--verbose] [--dry-run] <model> <prompt_file> <cr3_path> IMG_0150.CR3 IMG_0151.CR3

Flags:
  -v, --verbose   Print detailed per-file logs
  -n, --dry-run   Simulate actions without writing files or sending requests
  -h, --help      Show this help
EOF
}

main() {
  local model="gemma-4-e4b"
  local prompt_file="prompt.md"
  local cr3_path=""

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -v|--verbose)
        VERBOSE=1
        shift
        ;;
      -n|--dry-run)
        DRY_RUN=1
        shift
        ;;
      -h|--help)
        usage
        return 0
        ;;
      --)
        shift
        break
        ;;
      -*)
        echo "❌ Unknown flag: $1"
        usage
        return 1
        ;;
      *)
        break
        ;;
    esac
  done

  if [[ $# -lt 1 ]]; then
    usage
    return 1
  fi

  local first_arg
  first_arg=$(expand_path "$1")

  if [[ -d "$first_arg" ]]; then
    cr3_path="$first_arg"
    shift
  else
    if [[ $# -lt 3 ]]; then
      usage
      return 1
    fi

    model="$1"
    prompt_file=$(expand_path "$2")
    cr3_path=$(expand_path "$3")
    shift 3
  fi

  cr3_to_xmp_pipeline "$model" "$prompt_file" "$cr3_path" "$@"
}

if [[ "${(%):-%N}" == "$0" ]]; then
  main "$@"
fi

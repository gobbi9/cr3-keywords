# Repo-local overlay for cr3-keywords.
# Activate from repo root with: o

# `cr3` is available via mise PATH injection from `.mise.toml` (`{{config_root}}/bin`).

# Show repo-local man page staged by `make build`.
export def "man cr3" [] {
  let local_man = ($env.PWD | path join ".man")
  if not ($local_man | path exists) {
    error make --unspanned { msg: "Local man page not found. Run `make build` first." }
  }

  with-env { MANPATH: $local_man } {
    ^man cr3
  }
}

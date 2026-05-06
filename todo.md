Example argument for zsh shell script:
cr3_path=~/Pictures/raw/braunschweig-20260503

- jpgs should be saved in the `./jpgs/braunschweig-20260503` directory
- txt outputs should be saved in the `./outputs/braunschweig-20260503` directory
- tmp files should be saved in the `./tmp/braunschweig-20260503` directory
- xmp files should be saved in the `~/Pictures/raw/braunschweig-20260503` directory

Usages:

```bash
./cr3-keyword.sh <cr3_path>
./cr3-keyword.sh <cr3_path> IMG_0150.CR3 IMG_0151.CR3
./cr3-keyword.sh gemma-4-e4b prompt.md <cr3_path> IMG_0150.CR3 IMG_0151.CR3
./cr3-keyword.sh "~/Pictures/raw/braunschweig-20260503" IMG_0150.CR3
```

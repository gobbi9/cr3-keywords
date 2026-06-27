class Cr3Keywords < Formula
  desc "Local Lightroom AI keywording and captioning pipeline for Canon CR3 photos"
  homepage "https://github.com/gobbi9/cr3-keywords"
  version "{{VERSION}}"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/gobbi9/cr3-keywords/releases/download/v{{VERSION}}/cr3_{{VERSION}}_darwin-arm64"
      sha256 "{{SHA256_DARWIN_ARM64}}"
    else
      url "https://github.com/gobbi9/cr3-keywords/releases/download/v{{VERSION}}/cr3_{{VERSION}}_darwin-amd64"
      sha256 "{{SHA256_DARWIN_AMD64}}"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/gobbi9/cr3-keywords/releases/download/v{{VERSION}}/cr3_{{VERSION}}_linux-amd64"
      sha256 "{{SHA256_LINUX_AMD64}}"
    else
      odie "Unsupported architecture: #{Hardware::CPU.arch}"
    end
  end

  depends_on "exiftool"

  resource "manpage" do
    url "https://github.com/gobbi9/cr3-keywords/releases/download/v{{VERSION}}/cr3.1"
    sha256 "{{SHA256_MANPAGE}}"
  end

  def install
    bin.install Dir["cr3_{{VERSION}}_*"] .first => "cr3"

    resource("manpage").stage do
      man1.install "cr3.1"
    end
  end

  test do
    assert_match "Usage", shell_output("#{bin}/cr3 --help")
  end
end

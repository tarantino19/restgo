class Restapisummarizer < Formula
  desc "AI-powered REST API endpoint analyzer and summarizer"
  homepage "https://github.com/tarantino19/restgo"
  url "https://github.com/tarantino19/restgo/archive/v1.0.0.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"restapisummarizer"
  end

  test do
    assert_match "restapisummarizer", shell_output("#{bin}/restapisummarizer --version")
  end
end 
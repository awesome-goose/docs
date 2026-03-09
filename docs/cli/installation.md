# CLI Installation

Install the Goose CLI to scaffold applications and generate code.

## Installation

Download pre-built binaries from the [GitHub Releases page](https://github.com/awesome-goose/cli/releases).

### Linux

#### Install

```bash
# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    i686) ARCH="386" ;;
esac

# Download latest release
VERSION=$(curl -s https://api.github.com/repos/awesome-goose/cli/releases/latest | grep '"tag_name"' | cut -d'"' -f4)
curl -LO "https://github.com/awesome-goose/cli/releases/download/${VERSION}/goose-${VERSION}-linux-${ARCH}.tar.gz"

# Extract
tar -xzf "goose-${VERSION}-linux-${ARCH}.tar.gz"

# Install to /usr/local/bin
sudo mv "goose-${VERSION}-linux-${ARCH}" /usr/local/bin/goose
sudo chmod +x /usr/local/bin/goose

# Clean up
rm "goose-${VERSION}-linux-${ARCH}.tar.gz"

# Verify
goose --version
```

**Available architectures:**

- `amd64` - 64-bit (most common)
- `arm64` - ARM 64-bit (Raspberry Pi 4, AWS Graviton)
- `386` - 32-bit
- `arm` - ARM 32-bit (older Raspberry Pi)

#### Update

```bash
# Same as install - download and replace
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    i686) ARCH="386" ;;
esac

VERSION=$(curl -s https://api.github.com/repos/awesome-goose/cli/releases/latest | grep '"tag_name"' | cut -d'"' -f4)
curl -LO "https://github.com/awesome-goose/cli/releases/download/${VERSION}/goose-${VERSION}-linux-${ARCH}.tar.gz"
tar -xzf "goose-${VERSION}-linux-${ARCH}.tar.gz"
sudo mv "goose-${VERSION}-linux-${ARCH}" /usr/local/bin/goose
rm "goose-${VERSION}-linux-${ARCH}.tar.gz"
goose --version
```

#### Uninstall

```bash
sudo rm /usr/local/bin/goose
```

---

### macOS

#### Install

**Apple Silicon (M1/M2/M3/M4):**

```bash
# Download latest release
VERSION=$(curl -s https://api.github.com/repos/awesome-goose/cli/releases/latest | grep '"tag_name"' | cut -d'"' -f4)
curl -LO "https://github.com/awesome-goose/cli/releases/download/${VERSION}/goose-${VERSION}-darwin-arm64.tar.gz"

# Extract
tar -xzf "goose-${VERSION}-darwin-arm64.tar.gz"

# Install to /usr/local/bin
sudo mv "goose-${VERSION}-darwin-arm64" /usr/local/bin/goose
sudo chmod +x /usr/local/bin/goose

# Clean up
rm "goose-${VERSION}-darwin-arm64.tar.gz"

# Verify
goose --version
```

**Intel Mac:**

```bash
# Download latest release
VERSION=$(curl -s https://api.github.com/repos/awesome-goose/cli/releases/latest | grep '"tag_name"' | cut -d'"' -f4)
curl -LO "https://github.com/awesome-goose/cli/releases/download/${VERSION}/goose-${VERSION}-darwin-amd64.tar.gz"

# Extract
tar -xzf "goose-${VERSION}-darwin-amd64.tar.gz"

# Install to /usr/local/bin
sudo mv "goose-${VERSION}-darwin-amd64" /usr/local/bin/goose
sudo chmod +x /usr/local/bin/goose

# Clean up
rm "goose-${VERSION}-darwin-amd64.tar.gz"

# Verify
goose --version
```

> **Tip:** Run `uname -m` to check your architecture. `arm64` = Apple Silicon, `x86_64` = Intel.

**Auto-detect architecture:**

```bash
ARCH=$(uname -m)
[[ "$ARCH" == "x86_64" ]] && ARCH="amd64"

VERSION=$(curl -s https://api.github.com/repos/awesome-goose/cli/releases/latest | grep '"tag_name"' | cut -d'"' -f4)
curl -LO "https://github.com/awesome-goose/cli/releases/download/${VERSION}/goose-${VERSION}-darwin-${ARCH}.tar.gz"
tar -xzf "goose-${VERSION}-darwin-${ARCH}.tar.gz"
sudo mv "goose-${VERSION}-darwin-${ARCH}" /usr/local/bin/goose
rm "goose-${VERSION}-darwin-${ARCH}.tar.gz"
goose --version
```

#### Update

```bash
# Same as install - download and replace
ARCH=$(uname -m)
[[ "$ARCH" == "x86_64" ]] && ARCH="amd64"

VERSION=$(curl -s https://api.github.com/repos/awesome-goose/cli/releases/latest | grep '"tag_name"' | cut -d'"' -f4)
curl -LO "https://github.com/awesome-goose/cli/releases/download/${VERSION}/goose-${VERSION}-darwin-${ARCH}.tar.gz"
tar -xzf "goose-${VERSION}-darwin-${ARCH}.tar.gz"
sudo mv "goose-${VERSION}-darwin-${ARCH}" /usr/local/bin/goose
rm "goose-${VERSION}-darwin-${ARCH}.tar.gz"
goose --version
```

#### Uninstall

```bash
sudo rm /usr/local/bin/goose
```

---

### Windows

#### Install

**PowerShell (Run as Administrator):**

```powershell
# Set architecture (amd64, arm64, or 386)
$ARCH = "amd64"

# Get latest version
$VERSION = (Invoke-RestMethod -Uri "https://api.github.com/repos/awesome-goose/cli/releases/latest").tag_name

# Download
$URL = "https://github.com/awesome-goose/cli/releases/download/$VERSION/goose-$VERSION-windows-$ARCH.zip"
Invoke-WebRequest -Uri $URL -OutFile "goose.zip"

# Extract
Expand-Archive -Path "goose.zip" -DestinationPath "." -Force

# Create installation directory
New-Item -ItemType Directory -Force -Path "C:\Program Files\goose" | Out-Null

# Move binary
Move-Item "goose-$VERSION-windows-$ARCH.exe" "C:\Program Files\goose\goose.exe" -Force

# Clean up
Remove-Item "goose.zip"

# Add to PATH (permanent)
$currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($currentPath -notlike "*C:\Program Files\goose*") {
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;C:\Program Files\goose", "Machine")
}

Write-Host "Installation complete. Restart your terminal, then run: goose --version"
```

**Available architectures:**

- `amd64` - 64-bit (most common)
- `arm64` - ARM 64-bit (Windows on ARM)
- `386` - 32-bit

**Manual installation:**

1. Download the `.zip` file from [releases](https://github.com/awesome-goose/cli/releases)
2. Extract the archive
3. Move `goose-*-windows-*.exe` to `C:\Program Files\goose\goose.exe`
4. Add `C:\Program Files\goose` to your PATH:
   - Open **System Properties** → **Environment Variables**
   - Under **System variables**, find **Path** and click **Edit**
   - Click **New** and add `C:\Program Files\goose`
   - Click **OK** to save

#### Update

**PowerShell (Run as Administrator):**

```powershell
$ARCH = "amd64"
$VERSION = (Invoke-RestMethod -Uri "https://api.github.com/repos/awesome-goose/cli/releases/latest").tag_name
$URL = "https://github.com/awesome-goose/cli/releases/download/$VERSION/goose-$VERSION-windows-$ARCH.zip"

Invoke-WebRequest -Uri $URL -OutFile "goose.zip"
Expand-Archive -Path "goose.zip" -DestinationPath "." -Force
Move-Item "goose-$VERSION-windows-$ARCH.exe" "C:\Program Files\goose\goose.exe" -Force
Remove-Item "goose.zip"

goose --version
```

#### Uninstall

**PowerShell (Run as Administrator):**

```powershell
# Remove binary
Remove-Item "C:\Program Files\goose" -Recurse -Force

# Remove from PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
$newPath = ($currentPath -split ";" | Where-Object { $_ -ne "C:\Program Files\goose" }) -join ";"
[Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")

Write-Host "Goose CLI uninstalled"
```

---

## Installing a Specific Version

Replace `latest` with the desired version tag (e.g., `v1.0.0`):

**Linux/macOS:**

```bash
VERSION="v1.0.0"
# Then follow the install steps above
```

**Windows:**

```powershell
$VERSION = "v1.0.0"
# Then follow the install steps above
```

---

## Alternative: Using Go Install

If you have Go installed, you can also install via:

```bash
go install github.com/awesome-goose/cli@latest
```

This installs to `$GOPATH/bin`. Ensure this directory is in your PATH.

---

## Verifying Installation

```bash
goose --version
```

Expected output:

```
goose version v1.0.0
```

Run help to see available commands:

```bash
goose --help
```

---

## Verifying Download Integrity

Each release includes SHA256 checksums in `checksums.txt`.

**Linux/macOS:**

```bash
# Download checksums
curl -LO "https://github.com/awesome-goose/cli/releases/download/${VERSION}/checksums.txt"

# Verify
sha256sum -c checksums.txt --ignore-missing
```

**Windows PowerShell:**

```powershell
# Calculate hash of downloaded file
(Get-FileHash "goose.zip" -Algorithm SHA256).Hash

# Compare with checksums.txt from release page
```

---

## Troubleshooting

### Command Not Found

**Linux/macOS:**

```bash
# Check if binary exists
ls -la /usr/local/bin/goose

# Check PATH
echo $PATH | grep "/usr/local/bin"
```

**Windows:**

```powershell
# Check if binary exists
Test-Path "C:\Program Files\goose\goose.exe"

# Check PATH
$env:Path -split ";" | Select-String "goose"
```

### Permission Denied (Linux/macOS)

```bash
sudo chmod +x /usr/local/bin/goose
```

### macOS Security Warning

If macOS blocks the binary:

1. Go to **System Preferences** → **Security & Privacy**
2. Click **Allow Anyway** next to the blocked app message
3. Run `goose --version` again and click **Open**

Or remove the quarantine attribute:

```bash
xattr -d com.apple.quarantine /usr/local/bin/goose
```

---

## Next Steps

- [Creating Applications](creating-apps.md) - Create your first app
- [Generating Modules](generating-modules.md) - Add modules
- [Command Reference](reference.md) - Full documentation

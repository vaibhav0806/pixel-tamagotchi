const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const https = require("https");

const VERSION = require("./package.json").version;
const REPO = "vaibhav0806/pixel-tamagotchi";

function getPlatform() {
  const platform = process.platform;
  const arch = process.arch;

  const osMap = { linux: "linux", darwin: "darwin", win32: "windows" };
  const archMap = { x64: "amd64", arm64: "arm64" };

  const os = osMap[platform];
  const cpu = archMap[arch];

  if (!os || !cpu) {
    throw new Error(
      `Unsupported platform: ${platform} ${arch}. ` +
        `Download manually from https://github.com/${REPO}/releases`
    );
  }

  return { os, cpu, ext: platform === "win32" ? ".zip" : ".tar.gz" };
}

function download(url) {
  return new Promise((resolve, reject) => {
    https
      .get(url, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          return download(res.headers.location).then(resolve).catch(reject);
        }
        if (res.statusCode !== 200) {
          return reject(new Error(`Download failed: ${res.statusCode} ${url}`));
        }
        const chunks = [];
        res.on("data", (chunk) => chunks.push(chunk));
        res.on("end", () => resolve(Buffer.concat(chunks)));
        res.on("error", reject);
      })
      .on("error", reject);
  });
}

async function install() {
  const { os, cpu, ext } = getPlatform();
  const name = `pixel-tamagotchi_${os}_${cpu}`;
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${name}${ext}`;

  console.log(`Downloading pixel-tamagotchi v${VERSION} for ${os}/${cpu}...`);

  const data = await download(url);

  const binDir = path.join(__dirname, "bin");
  fs.mkdirSync(binDir, { recursive: true });

  const tmpFile = path.join(__dirname, `tmp${ext}`);
  fs.writeFileSync(tmpFile, data);

  const archiveBinName = os === "windows" ? "pixel-tamagotchi.exe" : "pixel-tamagotchi";
  const binName = os === "windows" ? "pixel-tamagotchi-bin.exe" : "pixel-tamagotchi-bin";

  try {
    if (ext === ".tar.gz") {
      execSync(`tar -xzf "${tmpFile}" -C "${binDir}" "${archiveBinName}"`, { stdio: "pipe" });
      // Rename to avoid overwriting the JS wrapper
      fs.renameSync(path.join(binDir, archiveBinName), path.join(binDir, binName));
    } else {
      // Windows zip — use PowerShell
      execSync(
        `powershell -command "Expand-Archive -Path '${tmpFile}' -DestinationPath '${binDir}' -Force"`,
        { stdio: "pipe" }
      );
      fs.renameSync(path.join(binDir, archiveBinName), path.join(binDir, binName));
    }

    const binPath = path.join(binDir, binName);
    if (os !== "windows") {
      fs.chmodSync(binPath, 0o755);
    }

    console.log(`Installed pixel-tamagotchi to ${binPath}`);
  } finally {
    fs.unlinkSync(tmpFile);
  }
}

install().catch((err) => {
  console.error(`Failed to install pixel-tamagotchi: ${err.message}`);
  process.exit(1);
});

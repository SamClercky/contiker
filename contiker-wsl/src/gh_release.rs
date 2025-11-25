use std::{
    fs::File,
    io::{BufWriter, Write},
    path::{Path, PathBuf},
};

use anyhow::{Context, bail};
use reqwest::{
    blocking::Client,
    header::{ACCEPT, HeaderMap, HeaderValue, USER_AGENT},
};
use serde::Deserialize;

#[derive(Debug, Clone, Deserialize)]
pub struct GHRelease {
    /// The assets related to a GH release
    pub assets: Vec<GHAsset>,
}

#[derive(Debug, Clone, Deserialize)]
pub struct GHAsset {
    /// The name of the asset
    pub name: String,
    /// The link to download the asset
    pub browser_download_url: String,
    /// The content type of the asset
    pub content_type: String,
}

impl GHAsset {
    pub fn download(&self, dir: &Path, client: &Client) -> anyhow::Result<PathBuf> {
        assert!(dir.is_dir(), "dir should be a directory");

        let mut headers = HeaderMap::new();
        headers.append(USER_AGENT, HeaderValue::from_static("Contiker-WSL"));
        headers.append(ACCEPT, HeaderValue::from_str(&self.content_type).unwrap());

        let mut res = client
            .get(&self.browser_download_url)
            .headers(headers)
            .send()
            .context("while requesting the asset")?;

        let out_file_path = dir.join(Path::new(&self.name));
        let out_file = File::create(out_file_path.clone())
            .context("while creating a file to write the downloaded distro to")?;
        let mut out_file = BufWriter::new(out_file);

        println!(
            "Dowloading {} into {}",
            &self.browser_download_url,
            out_file_path.display()
        );

        res.copy_to(&mut out_file)
            .context("while copying into file")?;

        out_file.flush()?;

        Ok(out_file_path)
    }
}

pub fn get_latest_release(client: &Client) -> anyhow::Result<GHRelease> {
    let mut headers = HeaderMap::new();
    headers.append(USER_AGENT, HeaderValue::from_static("Contiker-WSL"));
    headers.append(ACCEPT, HeaderValue::from_static("text/json"));

    let res = client
        .get("https://api.github.com/repos/SamClercky/contiker/releases/latest")
        .headers(headers)
        .send()
        .context("while retrieving the latest GH release")?;

    let release: String = res.text()?;
    let Ok(release) = serde_json::from_str(&release) else {
        bail!("could not parse github API response: {release}");
    };

    Ok(release)
}

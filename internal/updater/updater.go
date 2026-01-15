package updater

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
    "os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const RepoOwner = "Sw3bbl3"
const RepoName = "devdoctor"

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type githubRelease struct {
	TagName string         `json:"tag_name"`
	Assets  []releaseAsset `json:"assets"`
}

func LatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoOwner, RepoName)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		// Fallback to tags
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", RepoOwner, RepoName)
		resp2, err2 := client.Get(url)
		if err2 != nil {
			return "", err2
		}
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			return "", fmt.Errorf("unexpected status: %s", resp2.Status)
		}
		var tags []struct{ Name string `json:"name"` }
		if err := json.NewDecoder(resp2.Body).Decode(&tags); err != nil {
			return "", err
		}
		if len(tags) > 0 {
			return strings.TrimPrefix(tags[0].Name, "v"), nil
		}
		// Fallback to Go module resolution
		out, err := exec.Command("go", "list", "-m", "-json", fmt.Sprintf("github.com/%s/%s@latest", RepoOwner, RepoName)).Output()
		if err != nil {
			return "", errors.New("no tags found")
		}
		var mod struct{ Version string `json:"Version"` }
		if jerr := json.Unmarshal(out, &mod); jerr != nil {
			return "", jerr
		}
		return strings.TrimPrefix(mod.Version, "v"), nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var gr githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return "", err
	}
	return strings.TrimPrefix(gr.TagName, "v"), nil
}

func selectAsset(assets []releaseAsset) (releaseAsset, error) {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	var candidates []releaseAsset
	for _, a := range assets {
		name := strings.ToLower(a.Name)
		if strings.Contains(name, osName) && strings.Contains(name, arch) {
			candidates = append(candidates, a)
		}
	}
	if len(candidates) == 0 {
		return releaseAsset{}, errors.New("no matching asset found for platform")
	}
	return candidates[0], nil
}

func latestRelease() (githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoOwner, RepoName)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return githubRelease{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return githubRelease{}, errors.New("no releases")
	}
	if resp.StatusCode != http.StatusOK {
		return githubRelease{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var gr githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return githubRelease{}, err
	}
	return gr, nil
}

func destinationPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exe)
	name := "devdoctor"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	gobin := os.Getenv("GOBIN")
	if gobin == "" {
		gopath := os.Getenv("GOPATH")
		if gopath != "" {
			gobin = filepath.Join(gopath, "bin")
		} else {
			home, _ := os.UserHomeDir()
			gobin = filepath.Join(home, "go", "bin")
		}
	}
	dest := filepath.Join(gobin, name)
	if dest == exe {
		return dest, nil
	}
	if _, err := os.Stat(dest); err == nil {
		return dest, nil
	}
	return filepath.Join(dir, name), nil
}

func downloadWithProgress(url, outPath string) error {
	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	cl := resp.Header.Get("Content-Length")
	var total int64
	if cl != "" {
		if n, err := fmt.Sscanf(cl, "%d", &total); n == 1 && err == nil {
		}
	}
	var downloaded int64
	buf := make([]byte, 32*1024)
	lastPrint := time.Now()
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				return werr
			}
			downloaded += int64(n)
			if total > 0 && time.Since(lastPrint) > 500*time.Millisecond {
				pct := float64(downloaded) / float64(total) * 100
				fmt.Printf("[INFO] Downloading: %.1f%% (%.1f MB / %.1f MB)\r", pct, float64(downloaded)/1e6, float64(total)/1e6)
				lastPrint = time.Now()
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return readErr
		}
	}
	fmt.Print("\n")
	return nil
}

func UpdateToLatest(currentVersion string) (string, error) {
	gr, err := latestRelease()
	if err != nil {
		// Fallback: install from source
		cmd := exec.Command("go", "install", fmt.Sprintf("github.com/%s/%s/cmd/%s@latest", RepoOwner, RepoName, RepoName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}
		dest, derr := destinationPath()
		if derr != nil {
			return "", derr
		}
		return dest, nil
	}
	remote := strings.TrimPrefix(gr.TagName, "v")
	if remote == currentVersion {
		return "", fmt.Errorf("already up to date: %s", remote)
	}
	asset, err := selectAsset(gr.Assets)
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "devdoctor-update-*")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	if err := downloadWithProgress(asset.BrowserDownloadURL, tmpPath); err != nil {
		return "", err
	}
	dest, err := destinationPath()
	if err != nil {
		return "", err
	}
	if runtime.GOOS != "windows" {
		_ = os.Chmod(tmpPath, 0755)
	}
	if err := os.Rename(tmpPath, dest); err != nil {
		fallback := dest + ".new"
		if ferr := os.Rename(tmpPath, fallback); ferr != nil {
			return "", err
		}
		return fmt.Sprintf("downloaded to %s; replace existing binary after exit", fallback), nil
	}
	return dest, nil
}

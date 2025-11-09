# 🚀 GitHub Release Upload Guide

## 📦 Files Ready for Release v1.0.0

Tất cả files đã được chuẩn bị trong thư mục `releases/v1.0.0/`:

```
releases/v1.0.0/
├── ducla-agent_1.0.0_amd64.deb          (8.7M) - Ubuntu/Debian package
├── ducla-agent-1.0.0-1.x86_64.rpm       (8.7M) - RHEL/CentOS package  
├── ducla-agent-linux-amd64.tar.gz       (8.4M) - Binary distribution
├── checksums.txt                        (288B) - SHA256 checksums
└── RELEASE-NOTES.md                     (6.0K) - Release documentation
```

## 🎯 GitHub Release Steps

### 1. Tạo Release trên GitHub

1. Truy cập: https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases
2. Click **"Create a new release"**
3. Điền thông tin:

**Tag version:** `v1.0.0`  
**Release title:** `Ducla Cloud Agent v1.0.0 - First Stable Release`  
**Target:** `main` branch

### 2. Upload Files

Drag & drop hoặc click **"Attach binaries"** để upload các files:

- ✅ `ducla-agent_1.0.0_amd64.deb`
- ✅ `ducla-agent-1.0.0-1.x86_64.rpm`  
- ✅ `ducla-agent-linux-amd64.tar.gz`
- ✅ `checksums.txt`
- ✅ `RELEASE-NOTES.md`

### 3. Release Description

Copy nội dung từ `releases/v1.0.0/RELEASE-NOTES.md` vào phần description.

### 4. Publish Release

- ✅ Check **"Set as the latest release"**
- ✅ Click **"Publish release"**

## 🔗 Download Links Sau Khi Release

Sau khi publish, các links này sẽ hoạt động:

### Ubuntu/Debian Installation
```bash
wget https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0/ducla-agent_1.0.0_amd64.deb
sudo dpkg -i ducla-agent_1.0.0_amd64.deb
sudo systemctl enable ducla-agent
sudo systemctl start ducla-agent
```

### RHEL/CentOS Installation
```bash
wget https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0/ducla-agent-1.0.0-1.x86_64.rpm
sudo rpm -ivh ducla-agent-1.0.0-1.x86_64.rpm
sudo systemctl enable ducla-agent
sudo systemctl start ducla-agent
```

### Binary Installation
```bash
wget https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0/ducla-agent-linux-amd64.tar.gz
tar -xzf ducla-agent-linux-amd64.tar.gz
sudo cp ducla-agent /usr/local/bin/
sudo chmod +x /usr/local/bin/ducla-agent
```

### Verify Downloads
```bash
wget https://github.com/duclacloud/DUCLA-CLOUD-AGENT/releases/download/v1.0.0/checksums.txt
sha256sum -c checksums.txt
```

## ✅ Verification Checklist

Sau khi release, test các links:

- [ ] DEB package download và install
- [ ] RPM package download và install  
- [ ] Binary tar.gz download và extract
- [ ] Checksums verification
- [ ] Service start và API endpoints
- [ ] CLI commands hoạt động
- [ ] Man page available

## 🔐 Security Verification

### Package Checksums
```
5f51e17277262f203807dcae829aca87984728b38f1766bb952de8367238e644  ducla-agent_1.0.0_amd64.deb
b5c60969eea2d66a6cc6e088d69df5eb836293be56f44359a0e8bfe08df36768  ducla-agent-1.0.0-1.x86_64.rpm
214ab473ec177a4e567f2d3e3ee3194c9c08e6f38875c3cf23be3cd699eeb82a  ducla-agent-linux-amd64.tar.gz
```

### Verification Commands
```bash
# Verify DEB package
echo "5f51e17277262f203807dcae829aca87984728b38f1766bb952de8367238e644  ducla-agent_1.0.0_amd64.deb" | sha256sum -c

# Verify RPM package  
echo "b5c60969eea2d66a6cc6e088d69df5eb836293be56f44359a0e8bfe08df36768  ducla-agent-1.0.0-1.x86_64.rpm" | sha256sum -c

# Verify binary package
echo "214ab473ec177a4e567f2d3e3ee3194c9c08e6f38875c3cf23be3cd699eeb82a  ducla-agent-linux-amd64.tar.gz" | sha256sum -c
```

## 📢 Post-Release Actions

### 1. Update Documentation
- [ ] Update README.md với download links
- [ ] Update USER-GUIDE.md với correct URLs
- [ ] Update installation instructions

### 2. Announce Release
- [ ] GitHub Discussions post
- [ ] Social media announcement
- [ ] Email notifications to users
- [ ] Update project website

### 3. Monitor Release
- [ ] Watch for download statistics
- [ ] Monitor GitHub issues for problems
- [ ] Check user feedback
- [ ] Prepare hotfix if needed

## 🎉 Success Metrics

Sau 24 giờ, check:
- Download counts cho mỗi package
- GitHub stars và forks
- Issues reported
- User feedback

---

**🚀 Ready to make Ducla Cloud Agent v1.0.0 available to the world!**
# AWST Agent Telepítési Útmutató

Az AWST Agent egy reverse proxy, amely lehetővé teszi a biztonságos kapcsolódást a központi szerverhez.

## 📋 Előfeltételek

- Linux rendszer (Ubuntu 20.04 vagy újabb)
- Go 1.19+ (csak fordításhez)
- Docker és Docker Compose (a központi szerverhez)
- Root jogosultságok

## 🚀 Gyors telepítés

### 1. Agent fordítása és telepítése

```bash
cd Agent
make build
sudo make install

sudo awst-agent


# Indítás
sudo systemctl start awst-agent

# Leállítás
sudo systemctl stop awst-agent

# Újraindítás
sudo systemctl restart awst-agent

# Állapot ellenőrzés
sudo systemctl status awst-agent

# Automatikus indítás tiltása
sudo systemctl disable awst-agent

# Logok (utolsó 100 sor)
sudo journalctl -u awst-agent -n 100

# Logok folyamatos figyelése
sudo journalctl -u awst-agent -f



sudo make uninstall
# vagy manuálisan:
sudo systemctl stop awst-agent
sudo systemctl disable awst-agent
sudo rm -f /etc/systemd/system/awst-agent.service
sudo rm -f /usr/local/bin/awst-agent
sudo rm -rf /etc/awst
sudo systemctl daemon-reload
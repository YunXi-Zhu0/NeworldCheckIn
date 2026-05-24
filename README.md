# NeworldCheckin

每天自动签到的 Go 小工具。

## 配置

先准备 `config.yaml`：

```yaml
email: "your-email@example.com"
passwd: "your-password"
```

可直接复制 `config.yaml.example` 修改。

## Ubuntu 部署

直接执行安装脚本：

```bash
sudo bash ./scripts/install-ubuntu-checkin.sh
```

安装后会自动：

- 部署到 `/opt/neworld_check-in`
- 安装 `checkin.service` 和 `checkin.timer`
- 每天早上 `06:00-09:00` 随机签到一次

## Ubuntu 卸载

直接执行卸载脚本：

```bash
sudo bash ./scripts/uninstall-ubuntu-checkin.sh
```

卸载后会自动：

- 停止并禁用 `checkin.timer`
- 删除 `checkin.service` 和 `checkin.timer`
- 删除 `/opt/neworld_check-in`

## 常用命令

```bash
systemctl status checkin.timer
systemctl start checkin.service
journalctl -u checkin.service -n 50 --no-pager
```

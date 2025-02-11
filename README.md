# go-command-exec

Go学習用に書いたGoサンプルプログラムです。

SSHでリモートホストに接続して、コマンド結果を出力します。
SSHはパスワード認証と鍵認証両方に対応しています。

## 必要条件

少なくともGo 1.13から動くと思いますが、下記バージョンで動作を確認しています。

- Go 1.23.6 以上

## How to use

### Requirement

バイナリビルドに下記外部モジュールが必要です。

```bash
go get golang.org/x/crypto/ssh
go get golang.org/x/crypto/ssh/agent
go get gopkg.in/yaml.v2
```

### Build

バイナリにビルドする場合は下記のようにビルドします。

```bash
go build go-command-exec.go
```

ビルドせずに実行する場合は下記のように実行できます。

```bash
go run go-command-exec.go
```

### Setup

config.yamlファイルを編集して、接続先のサーバー情報と実行するコマンドを設定します。

```yaml
servers:
  - user: "vagrant"
    password: "vagrant"
    host: "192.168.56.201"
    port: "22"
    commands:
      - "uptime"
  - user: "vagrant"
    password: "vagrant"
    host: "192.168.56.202"
    port: "22"
    commands:
      - "df -k"
  - user: "vagrant"
    password: "vagrant"
    host: "192.168.56.203"
    port: "22"
    commands:
      - "free -m"
      - "sudo find /var/spool/postfix/ -type f | wc -l"
```

### Run

実行結果です。YAMLに定義された内容にしたがって順次SSH接続を行い下記のような出力が得られます。

```bash
❯ ./go-command-exec.go
Host: 192.168.56.201, ExecTime: 60.617088ms, ExecCommand: uptime
 15:27:22 up  6:11,  1 user,  load average: 0.00, 0.00, 0.00

Host: 192.168.56.202, ExecTime: 57.771437ms, ExecCommand: df -k
Filesystem                 1K-blocks    Used Available Use% Mounted on
devtmpfs                        4096       0      4096   0% /dev
tmpfs                         899424       0    899424   0% /dev/shm
tmpfs                         359772    6648    353124   2% /run
efivarfs                         256      33       219  13% /sys/firmware/efi/efivars
/dev/mapper/almalinux-root  13365248 3598912   9766336  27% /
/dev/sda2                     983040  337452    645588  35% /boot
/dev/sda1                     613160    7228    605932   2% /boot/efi
tmpfs                         179884       4    179880   1% /run/user/500

Host: 192.168.56.203, ExecTime: 67.66347ms, ExecCommand: free -m
               total        used        free      shared  buff/cache   available
Mem:            1756         339        1228           6         346        1416
Swap:           1639           0        1639

Host: 192.168.56.203, ExecTime: 82.15265ms, ExecCommand: sudo find /var/spool/postfix/ -type f | wc -l
6

```

# FlubioStudios SetupSuite
The SetupSuite was built to setup servers faster and more efficiantly

## Features

### -initial setup

## References
```shell
#!bin/bash
echo "Setting up server..."
adduser --disabled-password --gecos "" flubi0
groupAdd sshUserusermod -aG sshUser flubi0
update-alternatives --config editorecho 'flubi0 ALL=(ALL:ALL) ALL' | sudo EDITOR='tee -a' visudosu flubi0
mkdir /home/flubi0/.ssh/cat >> /home/flubi0/.ssh/authorized_keys << EOF
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC5G5f2XmkW/Zmv3G6YvrL4789hjYvidIWsVv7VjB2wAhrmH3gYzvI3SunBjVBzH/bgvEWffc21rKygAop4WSvZ6lF/J+Kipk7vu/yQN5RfcawJH/jv6J+J55yUViAYi+vqXTbkBzOlIY+WIBddKxiPrRVCQtQjcaN8HZwkpMXjUcJz9dIyUR+lQpFr9wCBr/sStKgTOxopMOIb9zsfNWB444ec6HgjgLLVd+SdD3g5fsvZR5dWWs9/dcCJSYRaoDD+ggzTXHSWO7sXKgSfq+hhQnizCZHPiDS69sMabZlIgwSOwXpVA/w5+nlYLsKcaPIoulAF2qU6v/3AQXQK7uRF
EOF
chown -R flubi0:flubi0 /home/flubi0/.ssh
chmod 700 /home/flubi0/.ssh
chmod 600 /home/flubi0/.ssh/authorized_keys
exit
cat >> /root/.bashrc << EOF
# ~/.bashrc: executed by bash(1) for non-login shells.

# Note: PS1 and umask are already set in /etc/profile. You should not
# need this unless you want different defaults for root.
# PS1='${debian_chroot:+($debian_chroot)}\h:\w\$ '
# umask 022

# You may uncomment the following lines if you want `ls' to be colorized:
export LS_OPTIONS='--color=auto'
# eval "`dircolors`"
alias ls='ls -la $LS_OPTIONS'
# alias ll='ls $LS_OPTIONS -l'
# alias l='ls $LS_OPTIONS -lA'
#
# Some more alias to avoid making mistakes:
# alias rm='rm -i'
# alias cp='cp -i'
# alias mv='mv -i'
PATH=$PATH:/usr/sbin
EOF
cp /root/.ssh/authorized_keys /root/.ssh/authorized_keys_copied_due_to_security
cat >> /etc/ssh/sshd_config << EOF
#       $OpenBSD: sshd_config,v 1.103 2018/04/09 20:41:22 tj Exp $

# This is the sshd server system-wide configuration file.  See
# sshd_config(5) for more information.

# This sshd was compiled with PATH=/usr/bin:/bin:/usr/sbin:/sbin

# The strategy used for options in the default sshd_config shipped with
# OpenSSH is to specify options with their default value where
# possible, but leave them commented.  Uncommented options override the
# default value.

Port 1101
#AddressFamily any
#ListenAddress 0.0.0.0
#ListenAddress ::

#HostKey /etc/ssh/ssh_host_rsa_key
#HostKey /etc/ssh/ssh_host_ecdsa_key
#HostKey /etc/ssh/ssh_host_ed25519_key

# Ciphers and keying
#RekeyLimit default none

# Logging
#SyslogFacility AUTH
#LogLevel INFO

# Authentication:

#LoginGraceTime 2m
PermitRootLogin no
#StrictModes yes
#MaxAuthTries 6
#MaxSessions 10
AllowGroups sshuser
PubkeyAuthentication yes

# Expect .ssh/authorized_keys2 to be disregarded by default in future.
#AuthorizedKeysFile     .ssh/authorized_keys .ssh/authorized_keys2

#AuthorizedPrincipalsFile none

#AuthorizedKeysCommand none
#AuthorizedKeysCommandUser nobody

# For this to work you will also need host keys in /etc/ssh/ssh_known_hosts
#HostbasedAuthentication no
# Change to yes if you don't trust ~/.ssh/known_hosts for
# HostbasedAuthentication
#IgnoreUserKnownHosts no
# Don't read the user's ~/.rhosts and ~/.shosts files
#IgnoreRhosts yes

# To disable tunneled clear text passwords, change to no here!
PasswordAuthentication no
PermitEmptyPasswords no

# Change to yes to enable challenge-response passwords (beware issues with
# some PAM modules and threads)
ChallengeResponseAuthentication no

# Kerberos options
#KerberosAuthentication no
#KerberosOrLocalPasswd yes
#KerberosTicketCleanup yes
#KerberosGetAFSToken no

# GSSAPI options
#GSSAPIAuthentication no
#GSSAPICleanupCredentials yes
#GSSAPIStrictAcceptorCheck yes
#GSSAPIKeyExchange no

# Set this to 'yes' to enable PAM authentication, account processing,
# and session processing. If this is enabled, PAM authentication will
# be allowed through the ChallengeResponseAuthentication and
# PasswordAuthentication.  Depending on your PAM configuration,
# PAM authentication via ChallengeResponseAuthentication may bypass
# the setting of "PermitRootLogin yes
# If you just want the PAM account and session checks to run without
# PAM authentication, then enable this but set PasswordAuthentication
# and ChallengeResponseAuthentication to 'no'.
UsePAM yes

#AllowAgentForwarding yes
#AllowTcpForwarding yes
#GatewayPorts no
X11Forwarding yes
#X11DisplayOffset 10
#X11UseLocalhost yes
#PermitTTY yes
PrintMotd no
#PrintLastLog yes
#TCPKeepAlive yes
#PermitUserEnvironment no
#Compression delayed
#ClientAliveInterval 0
#ClientAliveCountMax 3
#UseDNS no
#PidFile /var/run/sshd.pid
#MaxStartups 10:30:100
#PermitTunnel no
#ChrootDirectory none
#VersionAddendum none

# no default banner path
#Banner none

# Allow client to pass locale environment variables
AcceptEnv LANG LC_*

# override default of no subsystems
Subsystem sftp  /usr/lib/openssh/sftp-server

# Example of overriding settings on a per-user basis
#Match User anoncvs
#       X11Forwarding no
#       AllowTcpForwarding no
#       PermitTTY no
#       ForceCommand cvs server
PasswordAuthentication no
EOF
service sshd restart
apt-get upgrade && apt-get update -y
```

```
https://gist.github.com/kemelzaidan/9146624
```
```
https://gist.githubusercontent.com/kemelzaidan/4e16d95e81db1ed90a4a/raw/1fdea821501265a9275ef36cfdc2a76d09cda5e5/.bash_aliases
```
```
https://raw.githubusercontent.com/creationix/nvm/master/install.sh
```
```
https://github.com/ioniq-io/automated-server-setup
```
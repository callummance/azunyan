#!//bin/sh

sed -i -e 's/AllowTcpForwarding no/AllowTcpForwarding yes/g' /etc/ssh/sshd_config
sed -i -e 's/GatewayPorts no/GatewayPorts yes/g' /etc/ssh/sshd_config


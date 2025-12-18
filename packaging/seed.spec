Name:       seed
Version:    __VERSION__
Release:    1%{?dist}
Summary:    The Seed - Network Diagnostic Tool by Mustard Seed Networks
License:    BSL 1.1
URL:        https://github.com/krisarmstrong/seed
BuildArch:  __RPM_ARCH__

Requires:   libpcap, systemd, libcap

%description
The Seed is a professional-grade network diagnostic appliance designed
for network technicians and engineers. Plug it into any network jack and
instantly see link status, switch information, DHCP details, DNS health,
and gateway connectivity through a modern web interface.

Features:
- Real-time network diagnostics via web UI
- WiFi survey and heatmap generation
- Speed testing with iPerf3 integration
- SNMP device discovery
- DHCP rogue detection
- Cable diagnostics (TDR)

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}/etc/seed
mkdir -p %{buildroot}/lib/systemd/system
mkdir -p %{buildroot}/usr/share/seed
mkdir -p %{buildroot}/usr/local/seed/configs
mkdir -p %{buildroot}/usr/local/seed/logs

# Copy binaries
install -m 755 %{_repo_root}/dist/seed-linux-__ARCHITECTURE__ %{buildroot}/usr/local/bin/seed
install -m 755 %{_repo_root}/bin/iperf3-linux-__ARCHITECTURE__ %{buildroot}/usr/local/bin/iperf3

# Copy config files
cp -r %{_repo_root}/configs/* %{buildroot}/etc/seed/

# Copy web assets
cp -r %{_repo_root}/web/dist %{buildroot}/usr/share/seed/web

# Copy systemd service file
install -m 644 %{_repo_root}/packaging/seed.service %{buildroot}/lib/systemd/system/seed.service

%files
/usr/local/bin/seed
/usr/local/bin/iperf3
/etc/seed
/lib/systemd/system/seed.service
/usr/share/seed/web
%dir /usr/local/seed
%dir /usr/local/seed/configs
%dir /usr/local/seed/logs

%pre
# Create service user and group
getent group seed >/dev/null || groupadd -r seed
getent passwd seed >/dev/null || useradd -r -g seed -d /usr/local/seed -s /sbin/nologin -c "The Seed Service" seed

%post
# Set ownership of directories
chown -R seed:seed /usr/local/seed

# Copy default config if not exists
if [ ! -f /usr/local/seed/configs/seed.yaml ] && [ -f /etc/seed/seed.yaml ]; then
    cp /etc/seed/seed.yaml /usr/local/seed/configs/
    chown seed:seed /usr/local/seed/configs/seed.yaml
fi

# Set capabilities for raw socket access
setcap cap_net_raw=+ep /usr/local/bin/seed || true

%systemd_post seed.service

%preun
%systemd_preun seed.service

%postun
%systemd_postun_with_restart seed.service

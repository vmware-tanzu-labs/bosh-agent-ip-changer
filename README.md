# bosh-agent-ip-changer
Utility to update every BOSH agent's configured director IP

For Opsman deployed BOSH Directors, this utility quickly and easily updates all the BOSH deployed VMs to use the
new BOSH Director IP.

## Background
Typically the BOSH director IP should never change, however there may be certain scenarios where this is unavoidable such as:
1. After an AZ failure and you restore the BOSH director to an AZ with a different subnet.
1. After removing a secondary NIC from the BOSH director.
1. The BOSH director was originally deployed to an IP that is now needed for another use by your network team.

If you need to change the BOSH Director's IP for any reason, this tool can help. When you change a BOSH Director's IP
address all the BOSH deployed VMs will report back as unreachable because all the agents are still attempting to
connect to the director at it's old IP. The normal way to get the system back into a working state would be to run
`bosh cck` and recreate all VMs. There are a couple of downsides to recreating all VMs:

1. Drain scripts are not run which could cause downtime for some applications on TAS.
1. Execution time is _long_, since you're recreating all VMs.

This tool gives you another, better option.

## How it works
This tool grabs all deployments from your BOSH Director and then for each instance group in that deployment
it asks for the instance group VM password from Opsman. The tool then iterates over every VM in every deployment
and uses the instance group's password to SSH into each VM to reconfigure the BOSH agent to use the new BOSH
Director IP address and restart the BOSH agent. This is orders of magnitude faster and safer than recreating
all the VMs in every deployment. For a large foundation this will take under a minute, thereafter all BOSH
agents will report back as healthy. That's it, you're a hero!

## Usage

```shell
bosh-agent-ip-changer [OPTIONS] changeIP [changeIP-OPTIONS]

Help Options:
  -h, --help                     Show this help message

[changeIP command options]
          --ca-cert=             OpsManager CA certificate path or value [$OM_CA_CERT]
      -c, --client-id=           Client ID for the Ops Manager VM [$OM_CLIENT_ID]
      -s, --client-secret=       Client Secret for the Ops Manager VM [$OM_CLIENT_SECRET]
      -p, --password=            admin password for the Ops Manager VM [$OM_PASSWORD]
      -k, --skip-ssl-validation  skip ssl certificate validation during http requests [$OM_SKIP_SSL_VALIDATION]
      -u, --username=            admin username for the Ops Manager VM [$OM_USERNAME]
      -t, --target=              location of the Ops Manager VM [$OM_TARGET]
          --debug                sets log level to debug
          --old-director-ip=     previous ip of bosh director
          --new-director-ip=     new ip of bosh director
```

The `bosh-agent-ip-changer` uses the same environment variables as the [OM CLI](https://github.com/pivotal-cf/om).
The below example usage assumes you've already set up your OM CLI connection information, like OM_USERNAME,
OM_TARGET etc.

```shell
bosh-agent-ip-changer changeIP --old-director-ip 192.168.1.2 --new-director-ip 192.168.1.11
```

## Building
Go version 1.22+ is required.
```shell
go build -o ./bin/bosh-agent-ip-changer ./cmd
```

## Changing the BOSH Director's IP
This is out of scope of this tool, as this tool only updates the BOSH agent config, but here are some 
general instructions you can use as a starting point to change your BOSH Director's IP.

1. First enable Opsman Advanced Mode via the following OM CLI command:
```shell
om curl -x PUT -p /api/v0/staged/infrastructure/locked -d '{"locked" : "false"}'
```

2. Login to the Opsman UI and open the BOSH Director tile. On the `Assign AZs and Network` tab change the assigned
network as required.

3. The next set of steps require you to SSH into your Opsman VM and execute a few different commands. First we decrypt
the actual-installation.yml and installation.yml files:

```shell
sudo -u tempest-web SECRET_KEY_BASE="s" RAILS_ENV=production \
  /home/tempest-web/tempest/web/scripts/decrypt \
  /var/tempest/workspaces/default/actual-installation.yml \
  /tmp/actual-installation.yml

sudo -u tempest-web SECRET_KEY_BASE="s" RAILS_ENV=production \
  /home/tempest-web/tempest/web/scripts/decrypt \
  /var/tempest/workspaces/default/installation.yml \
  /tmp/installation.yml
```

4. Create a file on the Opsman VM, let's name it `/tmp/installation-ops.yml` with the following content:
```yaml
- type: remove
  path: /products/installation_name=p-bosh/director_ssl
- type: remove
  path: /products/installation_name=p-bosh/uaa_ssl
- type: remove
  path: /products/installation_name=p-bosh/blobstore_certificate
- type: remove
  path: /products/installation_name=p-bosh/director_agent_ssl
- type: remove
  path: /products/installation_name=p-bosh/credhub_ssl
- type: remove
  path: /products/installation_name=p-bosh/director_metrics_server_certificate
- type: remove
  path: /products/installation_name=p-bosh/director_metrics_server_client_certificate
- type: remove
  path: /products/installation_name=p-bosh/nats_server_certificate?
- type: remove
  path: /products/installation_name=p-bosh/bosh_metrics_forwarder_ssl
- type: remove
  path: /products/installation_name=p-bosh/director_configuration/allocated_director_ips/0
```

5. Use this Opsfile to reset or remove the various TLS certs that have the BOSH Director's old IP in their SAN by
running the following commands:
```shell
bosh int -o /tmp/installation-ops.yml /tmp/installation.yml > /tmp/new-installation.yml
bosh int -o /tmp/installation-ops.yml /tmp/actual-installation.yml > /tmp/new-actual-installation.yml
```

6. With the cert entries reset in the Opsman DB we re-encrypt files back to original location and restart the
Opsman web server:
```shell
sudo -u tempest-web SECRET_KEY_BASE="s" RAILS_ENV=production \
  /home/tempest-web/tempest/web/scripts/encrypt \
  /tmp/new-installation.yml \
  /var/tempest/workspaces/default/installation.yml

sudo -u tempest-web SECRET_KEY_BASE="s" RAILS_ENV=production \
  /home/tempest-web/tempest/web/scripts/encrypt \
  /tmp/new-actual-installation.yml \
  /var/tempest/workspaces/default/actual-installation.yml

sudo service tempest-web stop && sudo service tempest-web start
```

7. Login to the Opsman UI and `Apply Changes` to the Director tile only. Do *not* apply changes to any
of the other product tiles!

8. Use this tool to update all the BOSH deployed VMs agent configuration to the
new BOSH Director IP.

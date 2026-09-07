# Traefik edge routing

This directory contains the Portainer/Docker Swarm edge configuration used on
the Alibaba VPS for `chat-prospect.iainfinito.com.br` and the existing finance,
debounce, and Pedistore domains.

## Runtime layout

- `traefik.yaml` is the Swarm stack definition for the primary Traefik edge.
- `dynamic.yml` contains file-provider routers for services that still run on
  the host (`finance` systemd services and the legacy Traefik endpoint).
- The live host mounts the dynamic file at
  `/root/traefik-dynamic/dynamic.yml` and deploys the stack as `traefik`.
- Portainer/Swarm owns ports `80` and `443`; the host Nginx service is not part
  of the runtime.

## Important operational detail

The finance Next.js service must listen on `0.0.0.0:3100`, not only on
`127.0.0.1`, so the Traefik container can reach it through
`host.docker.internal`.

Do not commit ACME state, API keys, or provider credentials. The live ACME
storage remains in the external Docker volume `volume_swarm_certificates`.


## Deploy via Portainer + Traefik

The frontend is packaged for the existing Docker Swarm control plane. In Portainer, create/update a stack from `docker-stack.yml` after publishing the image `prospect-frontend:latest` to the registry used by the Swarm nodes.

The stack attaches to the external `network_public` network and exposes only the internal container port 80. Traefik owns the public route and TLS for `chat-prospect.iainfinito.com.br` through the `letsencryptresolver` resolver defined by the main Traefik stack.

Do not add an Nginx host reverse-proxy block for this application. Before rollout, the main Traefik service must be healthy on ports 80 and 443; the current legacy `traefik_estudioai` listener on 8080/8443 is not the production entrypoint.

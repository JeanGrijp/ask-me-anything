FROM n8nio/n8n:latest

USER root

# Instala o only-allow, cria pasta de extensões e instala o pacote localmente
RUN npm install -g only-allow && \
  mkdir -p /data/custom && \
  cd /data/custom && \
  npm init -y && \
  npm install n8n-nodes-evolution-api

ENV N8N_CUSTOM_EXTENSIONS=/data/custom

USER node
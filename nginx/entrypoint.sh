
# Garante que o script pare se houver um erro
set -e

# Substitui as variáveis de ambiente no template e cria o arquivo de configuração final
envsubst '${NGINX_FRONTEND_NAMES} ${API_DOMAIN}' < /etc/nginx/templates/nginx.conf.template > /etc/nginx/conf.d/app.conf

# Executa o comando original do contêiner Nginx (CMD)
exec "$@"
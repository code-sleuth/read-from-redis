FROM redis:7

MAINTAINER Ibrahim Mbaziira <code.ibra@gmail.com>

# startup script
COPY ./docker-data/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod 755 /docker-entrypoint.sh

EXPOSE 7005

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["redis-cluster"]

FROM ubuntu:latest
RUN apt-get update 
RUN apt-get install -y curl
COPY script.sh /script.sh 
COPY payload.json /payload.json
RUN chmod +x script.sh
RUN chmod +x payload.json
ENTRYPOINT ["/script.sh"]

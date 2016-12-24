FROM ubuntu:15.10
RUN apt-get update && \
    apt-get install -y vim wget
RUN apt-get install -y bash-completion
RUN wget -O /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.4.4/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl
ENV SSL_SCRIPT_COMMIT 98660ada3d800f653fc1f105771b5173f9d1a019
RUN wget -O /usr/bin/update-rancher-ssl https://raw.githubusercontent.com/rancher/rancher/${SSL_SCRIPT_COMMIT}/server/bin/update-rancher-ssl && \
    chmod +x /usr/bin/update-rancher-ssl
ENV KUBE_SERVER http://localhost:8080
COPY kubectld /usr/bin/
COPY kubectld.sh /usr/bin
COPY completion/kubectl /etc/bash_completion.d/
ENV SERVER http://localhost:8080
ENV LISTEN :8091
ENV HOME /tmp
RUN cp /etc/skel/.bashrc /tmp/ && \
    echo 'PS1="> "' >> /tmp/.bashrc && \
    echo . /etc/bash_completion >> /tmp/.bashrc
USER nobody
WORKDIR $HOME 
CMD ["kubectld.sh"]

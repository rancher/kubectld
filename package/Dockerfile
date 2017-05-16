FROM rancher/agent-base:v0.3.0
RUN apt-get update && \
    apt-get install -y vim wget
RUN apt-get install -y bash-completion
RUN wget -O /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.6.0/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl
RUN wget -O helm.tar.gz http://storage.googleapis.com/kubernetes-helm/helm-v2.3.0-linux-amd64.tar.gz && \
    tar xvzf helm.tar.gz && \
    mv ./linux-amd64/helm /usr/bin/ && \
    chmod +x /usr/bin/helm
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
    echo . /etc/bash_completion >> /tmp/.bashrc && \
    echo 'alias k="kubectl"' >> /tmp/.bashrc && \
    echo 'alias ks="kubectl -n kube-system"' >> /tmp/.bashrc
WORKDIR $HOME 
CMD ["kubectld.sh"]

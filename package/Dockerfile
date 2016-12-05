FROM ubuntu:15.10
RUN apt-get update && \
    apt-get install -y vim wget
RUN apt-get install -y bash-completion
RUN wget -O /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.4.4/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl
RUN wget http://storage.googleapis.com/kubernetes-helm/helm-v2.0.0-linux-386.tar.gz && \
    tar xvzf helm-v2.0.0-linux-386.tar.gz && \
    mv ./linux-386/helm /usr/bin/ && \
    chmod +x /usr/bin/helm
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

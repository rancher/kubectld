FROM ubuntu:15.10
RUN apt-get update && \
    apt-get install -y vim wget
RUN apt-get install -y bash-completion
RUN wget -O /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.3.0/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl
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

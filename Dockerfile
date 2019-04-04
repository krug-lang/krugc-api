FROM golang

RUN apt-get update && apt-get install -y ca-certificates git-core ssh

ADD keys/my_key_rsa /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/

ADD . /go/src/github.com/hugobrains/caasper

RUN go get github.com/hugobrains/caasper
RUN go install github.com/hugobrains/caasper

EXPOSE 8001
FROM scratch
ADD cerner_ca ./cerner_ca
ADD ca-certificates.crt /etc/ssl/certs/
COPY sample.pdf ./sample.pdf
EXPOSE 9000
ENTRYPOINT ["/cerner_ca"]
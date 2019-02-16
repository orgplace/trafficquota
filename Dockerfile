FROM scratch

COPY ./trafficquotad /trafficquotad

ENTRYPOINT [ "/trafficquotad" ]

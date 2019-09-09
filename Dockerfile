FROM golang:1.12.5-stretch as builder

RUN mkdir /backy-scraper
WORKDIR /backy-scraper

ADD go.mod .

RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/backy-scraper .

FROM flaviostutz/ceph-client:13.2.5

RUN apt-get update && apt-get -y install python3-alembic python3-dateutil python3-prettytable python3-psutil python3-setproctitle python3-shortuuid python3-sqlalchemy python3-psycopg2 netcat 
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y python3-boto python3-azure-storage
RUN apt-get install -y librados-dev librbd-dev rbd-nbd nbd-client 
RUN apt install curl -y

RUN wget https://github.com/jairsjunior/backy2/raw/master/dist/backy2_2.9.18_all.deb -O /backy2.deb
RUN dpkg -i backy2.deb
RUN rm /backy2.deb
RUN apt-get -y -f install

COPY --from=builder /go/bin/backy-scraper /bin/backy-scraper

COPY checkIsUp.sh /
COPY diff-bkp.sh /
ADD backy.cfg.template /
RUN chmod +x diff-bkp.sh
RUN touch /var/log/backy.log

ENV SOURCE_DATA_PATH ''

#file (will store in /var/lib/backy2/data)
#s3 (must be configured with ENVs below)
ENV TARGET_DATA_BACKEND 'file'

ENV S3_AWS_ACCESS_KEY_ID ''
ENV S3_AWS_SECRET_ACCESS_KEY ''
ENV S3_AWS_HOST ''
ENV S3_AWS_PORT '443'
ENV S3_AWS_HTTPS 'true'
ENV S3_AWS_BUCKET_NAME ''

ENV AZURE_ACCESS_KEY_ID ''
ENV AZURE_SECRET_ACCESS_KEY ''
ENV AZURE_BUCKET_NAME ''

ENV SIMULTANEOUS_WRITES '3'
ENV MAX_BANDWIDTH_WRITE '0'
ENV SIMULTANEOUS_READS '10'
ENV MAX_BANDWIDTH_READ '0'
ENV PROTECT_YOUNG_BACKUP_DAYS '6'

ENV PRE_POST_TIMEOUT '7200'
ENV PRE_BACKUP_COMMAND ''
ENV POST_BACKUP_COMMAND ''
ENV INDEX_DB_ADDRESS 'sqlite:////var/lib/backy2/backy.sqlite'
ENV CONDUCTOR_WORKNAME 'backup'

CMD ["bash","checkIsUp.sh"]

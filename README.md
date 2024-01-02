# Dedup
dedup removes the duplicate files for a given directory

# Start dedup server

```
docker run -dit -p 9090:8080 --rm \
  -v dir1:/tmp/dedup/dir1 -v dir2:/tmp/dedup/dir2 \
  email4tong/dedup:latest
```
The mounted directory dir1 and dir2 will be processed altogether, you may
mount many directories to the container under /tmp/dedup directory to process
more files

To start drurun:
```
curl http://localhost:9090/start?dryrun=true 
```
To start real run which also removes the duplicated files:
```
curl http://localhost:9090/start?dryrun=false
```
To stop the process:
```
curl http://localhost:9090/stop
```

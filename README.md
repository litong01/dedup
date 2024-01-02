# dedup
dedup removes the duplicate files for a given directory

# To use dedup with docker

```
docker run -dit -p 9090:8080 --rm \
  -v dir1:/tmp/dedup/dir1 -v dir2:/tmp/dedup/dir2 \
  email4tong/dedup:latest
```

Then use browser to hit the server to start or stop the dedup
process.

To start drurun,
curl http://localhost:9090/start?drurun=true 

To start real run which also removes the duplicated files
curl http://localhost:9090/start?drurun=false

To stop the process
curl http://localhost:9090/stop

# dedup
dedup removes the duplicate files for a given directory

# To use dedup with docker

```
docker run -dit -p 9090:8080 --rm \
  -v dir1:/tmp/dedup/dir1 -v dir2:/tmp/dedup/dir2 \
  email4tong/dedup:latest
```
docker run -d -p 9189:9189 \
--name=gluster_exporter \
--restart=always -u root \
-v /usr/bin/docker:/usr/bin/docker \
-v /var/run/docker.sock:/var/run/docker.sock \
-v /usr/lib64/libltdl.so.7:/usr/lib/x86_64-linux-gnu/libltdl.so.7 \
gluster_exporter:latest

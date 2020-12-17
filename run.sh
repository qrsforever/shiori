## !/bin/bash

# SHIORI_DIR=/blog/source/_shiori SOCKS5_PROXY=127.0.0.1:1080 /data/shiori/shiori serve -p 8698

docker run -dit --name shiori --network host \
    --restart unless-stopped --env SOCKS5_PROXY=127.0.0.1:1080 \
    --volume /blog/source/_shiori:/srv/shiori \
    --entrypoint /bin/sh shiori \
    -c "shiori serve -p 8698"
